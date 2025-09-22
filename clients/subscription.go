package clients

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SubscribeNewHeads subscribes to new block as soon as blocked are mined
func (c *Client) SubscribeNewHeads(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	if !c.isWS {
		return nil, fmt.Errorf("SubscribeNewHeads requires a WebSocket connection")
	}

	return c.Eth.SubscribeNewHead(ctx, ch)
}

// SubscribeLogs subscribes to smart contract event logs
func (c *Client) SubscribeLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	if !c.isWS {
		return nil, fmt.Errorf("SubscribeLogs requires a WebSocket connection")
	}
	return c.Eth.SubscribeFilterLogs(ctx, query, ch)
}

// SubscribePendingTxs subscribes to new transactions entering the mempool
func (c *Client) SubscribePendingTxs(ctx context.Context, ch chan<- common.Hash) (ethereum.Subscription, error) {
	if !c.isWS {
		return nil, fmt.Errorf("SubscribePendingTxs requires a WebSocket connection")
	}
	return c.RpcClient.EthSubscribe(ctx, ch, "newPendingTransactions")
}
