package clients

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	Eth  *ethclient.Client
	isWS bool
}

// Dial connects to an Abstract RPC node
func Dial(url string) (*Client, error) {
	eth, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		Eth:  eth,
		isWS: strings.HasPrefix(url, "ws"),
	}, nil
}

// Close closes websocket connection
func (c *Client) Close() {
	c.Eth.Close()
}

// BalanceAt queries the balance of an address
func (c *Client) BalanceAt(ctx context.Context, addr common.Address) (*big.Int, error) {
	return c.Eth.BalanceAt(ctx, addr, nil)
}

// NonceAt queries the account nonce
func (c *Client) NonceAt(ctx context.Context, addr common.Address) (uint64, error) {
	return c.Eth.NonceAt(ctx, addr, nil)
}

// GasPrice returns current gas price
func (c *Client) GasPrice(ctx context.Context) (*big.Int, error) {
	return c.Eth.SuggestGasPrice(ctx)
}

// CallContract performs a read-only contract call
func (c *Client) CallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	return c.Eth.CallContract(ctx, msg, nil)
}

// SendTransaction sends a signed transaction
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.Eth.SendTransaction(ctx, tx)
}
