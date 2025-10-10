package clients

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	Eth       *ethclient.Client
	RpcClient *rpc.Client
	isWS      bool
}

// DialHTTP creates a client for HTTP connections (query & tx)
func DialHTTP(url string) (*Client, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, fmt.Errorf("DialHTTP requires an http:// or https:// URL")
	}

	eth, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		Eth:  eth,
		isWS: false,
	}, nil
}

// DialWS creates a client for WebSocket connections (subscriptions)
func DialWS(url string) (*Client, error) {
	if !strings.HasPrefix(url, "ws") {
		return nil, fmt.Errorf("DialWS requires a ws:// or wss:// URL")
	}

	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		Eth:       ethclient.NewClient(rpcClient),
		RpcClient: rpcClient,
		isWS:      true,
	}, nil
}

// Close closes websocket connection
func (c *Client) Close() {
	c.Eth.Close()
}

// BalanceAt queries the balance of an address
func (c *Client) BalanceAt(ctx context.Context, addr common.Address) (*big.Int, error) {
	if c.isWS {
		return nil, fmt.Errorf("BalanceAt requires an HTTP connection, not WebSocket")
	}

	return c.Eth.BalanceAt(ctx, addr, nil)
}

// NonceAt queries the account nonce
func (c *Client) NonceAt(ctx context.Context, addr common.Address) (uint64, error) {
	if c.isWS {
		return 0, fmt.Errorf("NonceAt requires an HTTP connection, not WebSocket")
	}

	return c.Eth.NonceAt(ctx, addr, nil)
}

// GasPrice returns current gas price
func (c *Client) GasPrice(ctx context.Context) (*big.Int, error) {
	if c.isWS {
		return nil, fmt.Errorf("GasPrice requires an HTTP connection, not WebSocket")
	}

	return c.Eth.SuggestGasPrice(ctx)
}

// CallContract performs a read-only contract call
func (c *Client) CallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	if c.isWS {
		return nil, fmt.Errorf("CallContract requires an HTTP connection, not WebSocket")
	}

	return c.Eth.CallContract(ctx, msg, nil)
}

// SendTransaction sends a signed transaction
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if c.isWS {
		return fmt.Errorf("SendTransaction requires an HTTP connection, not WebSocket")
	}

	return c.Eth.SendTransaction(ctx, tx)
}

// EstimateGasWithBuffer estimates gas for a CallMsg and applies a buffer percentage.
// bufferPercent is e.g., 10 for +10%.
func (c *Client) EstimateGasWithBuffer(ctx context.Context, msg ethereum.CallMsg, bufferPercent uint64) (uint64, error) {
	gas, err := c.Eth.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}

	// Apply buffer
	buffer := gas * bufferPercent / 100
	gasWithBuffer := gas + buffer

	return gasWithBuffer, nil
}
