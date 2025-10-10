package clients

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const MinimalERC721ABI = `[
	{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"type":"function"}
]`

type ERC721 struct {
	client *Client
	addr   common.Address
	abi    abi.ABI
}

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

// BalanceOf returns how many NFTs the owner owns
func (e *ERC721) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	data, err := e.abi.Pack("balanceOf", owner)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &e.addr,
		Data: data,
	}

	res, err := e.client.CallContract(ctx, msg)
	if err != nil {
		return nil, err
	}

	balance := new(big.Int)
	err = e.abi.UnpackIntoInterface(&balance, "balanceOf", res)
	return balance, err
}

// OwnerOf returns the owner of a specific token
func (e *ERC721) OwnerOf(ctx context.Context, tokenId *big.Int) (common.Address, error) {
	data, err := e.abi.Pack("ownerOf", tokenId)
	if err != nil {
		return common.Address{}, err
	}

	msg := ethereum.CallMsg{
		To:   &e.addr,
		Data: data,
	}

	res, err := e.client.CallContract(ctx, msg)
	if err != nil {
		return common.Address{}, err
	}

	var owner common.Address
	err = e.abi.UnpackIntoInterface(&owner, "ownerOf", res)
	return owner, err
}

// TokenURI returns the metadata URI of a token
func (e *ERC721) TokenURI(ctx context.Context, tokenId *big.Int) (string, error) {
	data, err := e.abi.Pack("tokenURI", tokenId)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   &e.addr,
		Data: data,
	}

	res, err := e.client.CallContract(ctx, msg)
	if err != nil {
		return "", err
	}

	var uri string
	err = e.abi.UnpackIntoInterface(&uri, "tokenURI", res)
	return uri, err
}

// TransferFrom transfers an NFT from `from` to `to`
func (e *ERC721) TransferFrom(ctx context.Context, wallet *Wallet, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	data, err := e.abi.Pack("transferFrom", from, to, tokenId)
	if err != nil {
		return nil, err
	}

	nm := NewNonceManager(e.client, wallet.Address)
	return wallet.BuildAndSendTx(ctx, e.client, &e.addr, big.NewInt(0), data, nm)
}
