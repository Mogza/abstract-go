package clients

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const MinimalERC721ABI = `[
	{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"type":"function"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"approved","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"}
]`

type ERC721 struct {
	client *Client
	addr   common.Address
	abi    abi.ABI
}

type ERC721TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenID *big.Int
}

type ERC721ApprovalEvent struct {
	Owner    common.Address
	Approved common.Address
	TokenID  *big.Int
}

// --- Constructor ---
func NewERC721(client *Client, contractAddr common.Address, abiJSON string) (*ERC721, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if abiJSON == "" {
		abiJSON = MinimalERC721ABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	return &ERC721{
		client: client,
		addr:   contractAddr,
		abi:    parsedABI,
	}, nil
}

// --- Read-only calls ---
func (e *ERC721) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	data, _ := e.abi.Pack("balanceOf", owner)
	res, err := e.client.CallContract(ctx, ethereum.CallMsg{To: &e.addr, Data: data})
	if err != nil {
		return nil, err
	}
	balance := new(big.Int)
	return balance, e.abi.UnpackIntoInterface(&balance, "balanceOf", res)
}

func (e *ERC721) OwnerOf(ctx context.Context, tokenID *big.Int) (common.Address, error) {
	data, _ := e.abi.Pack("ownerOf", tokenID)
	res, err := e.client.CallContract(ctx, ethereum.CallMsg{To: &e.addr, Data: data})
	if err != nil {
		return common.Address{}, err
	}
	var owner common.Address
	return owner, e.abi.UnpackIntoInterface(&owner, "ownerOf", res)
}

func (e *ERC721) TokenURI(ctx context.Context, tokenID *big.Int) (string, error) {
	data, _ := e.abi.Pack("tokenURI", tokenID)
	res, err := e.client.CallContract(ctx, ethereum.CallMsg{To: &e.addr, Data: data})
	if err != nil {
		return "", err
	}
	var uri string
	return uri, e.abi.UnpackIntoInterface(&uri, "tokenURI", res)
}

// --- Write transactions ---
func (e *ERC721) TransferFrom(ctx context.Context, wallet *Wallet, from, to common.Address, tokenID *big.Int) (*types.Transaction, error) {
	data, _ := e.abi.Pack("transferFrom", from, to, tokenID)
	nm := NewNonceManager(e.client, wallet.Address)
	return wallet.BuildAndSendTx(ctx, e.client, &e.addr, big.NewInt(0), data, nm)
}

// --- Watchers ---
func (e *ERC721) WatchTransfers(ctx context.Context, ch chan<- ERC721TransferEvent) error {
	if !e.client.isWS {
		return fmt.Errorf("WatchTransfers requires a WS client")
	}
	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.addr},
		Topics:    [][]common.Hash{{transferSig}},
	}
	logsCh := make(chan types.Log)
	sub, err := e.client.SubscribeLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				fmt.Println("subscription error:", err)
				return
			case vLog := <-logsCh:
				event, err := parseERC721TransferLog(vLog)
				if err != nil {
					continue
				}
				ch <- *event
			}
		}
	}()
	return nil
}

func (e *ERC721) WatchApprovals(ctx context.Context, ch chan<- ERC721ApprovalEvent) error {
	if !e.client.isWS {
		return fmt.Errorf("WatchApprovals requires a WS client")
	}
	approvalSig := crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.addr},
		Topics:    [][]common.Hash{{approvalSig}},
	}
	logsCh := make(chan types.Log)
	sub, err := e.client.SubscribeLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				fmt.Println("subscription error:", err)
				return
			case vLog := <-logsCh:
				event, err := parseERC721ApprovalLog(vLog)
				if err != nil {
					continue
				}
				ch <- *event
			}
		}
	}()
	return nil
}

// --- Log parsers ---
func parseERC721TransferLog(vLog types.Log) (*ERC721TransferEvent, error) {
	event := new(ERC721TransferEvent)
	if len(vLog.Topics) < 4 {
		return nil, errors.New("invalid transfer log")
	}
	event.From = common.HexToAddress(vLog.Topics[1].Hex())
	event.To = common.HexToAddress(vLog.Topics[2].Hex())
	event.TokenID = new(big.Int).SetBytes(vLog.Topics[3].Bytes())
	return event, nil
}

func parseERC721ApprovalLog(vLog types.Log) (*ERC721ApprovalEvent, error) {
	event := new(ERC721ApprovalEvent)
	if len(vLog.Topics) < 4 {
		return nil, errors.New("invalid approval log")
	}
	event.Owner = common.HexToAddress(vLog.Topics[1].Hex())
	event.Approved = common.HexToAddress(vLog.Topics[2].Hex())
	event.TokenID = new(big.Int).SetBytes(vLog.Topics[3].Bytes())
	return event, nil
}
