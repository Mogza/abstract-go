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

const MinimalERC20ABI = `[
	{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
	{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
	{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"anonymous":false,"inputs":[
		{"indexed":true,"name":"from","type":"address"},
		{"indexed":true,"name":"to","type":"address"},
		{"indexed":false,"name":"value","type":"uint256"}],
		"name":"Transfer","type":"event"},
	{"anonymous":false,"inputs":[
		{"indexed":true,"name":"owner","type":"address"},
		{"indexed":true,"name":"spender","type":"address"},
		{"indexed":false,"name":"value","type":"uint256"}],
		"name":"Approval","type":"event"}
]`

type ERC20 struct {
	client *Client
	addr   common.Address
	abi    abi.ABI
}

type ERC20TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type ERC20ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

// NewERC20 creates an ERC20 contract binding with the given client and address.
// Parses the provided ABI JSON or uses the minimal default.
func NewERC20(client *Client, token common.Address, abiJSON string) (*ERC20, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if abiJSON == "" {
		abiJSON = MinimalERC20ABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	return &ERC20{
		client: client,
		addr:   token,
		abi:    parsedABI,
	}, nil
}

// BalanceOf returns the token balance of the given address.
// Calls the ERC20 `balanceOf` method as a read-only contract call.
func (t *ERC20) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	data, _ := t.abi.Pack("balanceOf", owner)
	msg := ethereum.CallMsg{To: &t.addr, Data: data}
	res, err := t.client.CallContract(ctx, msg)
	if err != nil {
		return nil, err
	}
	balance := new(big.Int)
	err = t.abi.UnpackIntoInterface(&balance, "balanceOf", res)
	return balance, err
}

// Allowance returns the remaining tokens a spender can spend from an owner's account.
// Calls the ERC20 `allowance` method as a read-only contract call.
func (t *ERC20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	data, _ := t.abi.Pack("allowance", owner, spender)
	msg := ethereum.CallMsg{To: &t.addr, Data: data}
	res, err := t.client.CallContract(ctx, msg)
	if err != nil {
		return nil, err
	}
	allowance := new(big.Int)
	err = t.abi.UnpackIntoInterface(&allowance, "allowance", res)
	return allowance, err
}

// Name returns the name of the ERC20 token.
// Calls the ERC20 `name` method as a read-only contract call.
func (t *ERC20) Name(ctx context.Context) (string, error) {
	data, _ := t.abi.Pack("name")
	msg := ethereum.CallMsg{To: &t.addr, Data: data}
	res, err := t.client.CallContract(ctx, msg)
	if err != nil {
		return "", err
	}
	var name string
	err = t.abi.UnpackIntoInterface(&name, "name", res)
	return name, err
}

// Symbol returns the symbol of the ERC20 token.
// Calls the ERC20 `symbol` method as a read-only contract call.
func (t *ERC20) Symbol(ctx context.Context) (string, error) {
	data, _ := t.abi.Pack("symbol")
	msg := ethereum.CallMsg{To: &t.addr, Data: data}
	res, err := t.client.CallContract(ctx, msg)
	if err != nil {
		return "", err
	}
	var symbol string
	err = t.abi.UnpackIntoInterface(&symbol, "symbol", res)
	return symbol, err
}

// Decimals returns the number of decimals used by the ERC20 token.
// Calls the ERC20 `decimals` method as a read-only contract call.
func (t *ERC20) Decimals(ctx context.Context) (uint8, error) {
	data, _ := t.abi.Pack("decimals")
	msg := ethereum.CallMsg{To: &t.addr, Data: data}
	res, err := t.client.CallContract(ctx, msg)
	if err != nil {
		return 0, err
	}
	var decimals uint8
	err = t.abi.UnpackIntoInterface(&decimals, "decimals", res)
	return decimals, err
}

// Transfer sends a transaction to transfer tokens to another address.
// Calls the ERC20 `transfer` method as a write transaction.
func (t *ERC20) Transfer(ctx context.Context, wallet *Wallet, to common.Address, amount *big.Int) (*types.Transaction, error) {
	data, _ := t.abi.Pack("transfer", to, amount)
	nm := NewNonceManager(t.client, wallet.Address)
	return wallet.BuildAndSendTx(ctx, t.client, &t.addr, big.NewInt(0), data, nm)
}

// TransferFrom sends a transaction to transfer tokens from one address to another.
// Calls the ERC20 `transferFrom` method as a write transaction.
func (t *ERC20) TransferFrom(ctx context.Context, wallet *Wallet, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	data, _ := t.abi.Pack("transferFrom", from, to, amount)
	nm := NewNonceManager(t.client, wallet.Address)
	return wallet.BuildAndSendTx(ctx, t.client, &t.addr, big.NewInt(0), data, nm)
}

// Approve sends a transaction to approve a spender for a specific amount.
// Calls the ERC20 `approve` method as a write transaction.
func (t *ERC20) Approve(ctx context.Context, wallet *Wallet, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	data, _ := t.abi.Pack("approve", spender, amount)
	nm := NewNonceManager(t.client, wallet.Address)
	return wallet.BuildAndSendTx(ctx, t.client, &t.addr, big.NewInt(0), data, nm)
}

// WatchTransfers subscribes to Transfer events and sends them to the provided channel.
// Requires a WebSocket client; parses logs into ERC20TransferEvent structs.
func (t *ERC20) WatchTransfers(ctx context.Context, from, to *common.Address, ch chan<- ERC20TransferEvent) error {
	if !t.client.isWS {
		return fmt.Errorf("WatchTransfers requires a WS client")
	}

	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	query := ethereum.FilterQuery{
		Addresses: []common.Address{t.addr},
		Topics:    [][]common.Hash{{transferSig}},
	}
	if from != nil {
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(from.Hex())})
	}
	if to != nil {
		if len(query.Topics) < 2 {
			query.Topics = append(query.Topics, nil)
		}
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(to.Hex())})
	}

	logsCh := make(chan types.Log)
	sub, err := t.client.SubscribeLogs(ctx, query, logsCh)
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
				event, err := t.parseTransferLog(vLog)
				if err != nil {
					continue
				}
				ch <- *event
			}
		}
	}()
	return nil
}

// WatchApprovals subscribes to Approval events and sends them to the provided channel.
// Requires a WebSocket client; parses logs into ERC20ApprovalEvent structs.
func (t *ERC20) WatchApprovals(ctx context.Context, owner, spender *common.Address, ch chan<- ERC20ApprovalEvent) error {
	if !t.client.isWS {
		return fmt.Errorf("WatchApprovals requires a WS client")
	}

	approvalSig := crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

	query := ethereum.FilterQuery{
		Addresses: []common.Address{t.addr},
		Topics:    [][]common.Hash{{approvalSig}},
	}
	if owner != nil {
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(owner.Hex())})
	}
	if spender != nil {
		if len(query.Topics) < 2 {
			query.Topics = append(query.Topics, nil)
		}
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(spender.Hex())})
	}

	logsCh := make(chan types.Log)
	sub, err := t.client.SubscribeLogs(ctx, query, logsCh)
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
				event, err := t.parseApprovalLog(vLog)
				if err != nil {
					continue
				}
				ch <- *event
			}
		}
	}()
	return nil
}

// parseTransferLog parses a Transfer event log into an ERC20TransferEvent struct.
// Returns an error if the log format is invalid.
func (t *ERC20) parseTransferLog(vLog types.Log) (*ERC20TransferEvent, error) {
	if len(vLog.Topics) < 3 {
		return nil, fmt.Errorf("invalid ERC20 Transfer log")
	}
	event := &ERC20TransferEvent{
		From:  common.HexToAddress(vLog.Topics[1].Hex()),
		To:    common.HexToAddress(vLog.Topics[2].Hex()),
		Value: new(big.Int),
	}

	// Use the instance ABI
	err := t.abi.UnpackIntoInterface(event, "Transfer", vLog.Data)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// parseApprovalLog parses an Approval event log into an ERC20ApprovalEvent struct.
// Returns an error if the log format is invalid.
func (t *ERC20) parseApprovalLog(vLog types.Log) (*ERC20ApprovalEvent, error) {
	if len(vLog.Topics) < 3 {
		return nil, fmt.Errorf("invalid ERC20 Approval log")
	}
	event := &ERC20ApprovalEvent{
		Owner:   common.HexToAddress(vLog.Topics[1].Hex()),
		Spender: common.HexToAddress(vLog.Topics[2].Hex()),
		Value:   new(big.Int),
	}

	// Use the instance ABI
	err := t.abi.UnpackIntoInterface(event, "Approval", vLog.Data)
	if err != nil {
		return nil, err
	}

	return event, nil
}
