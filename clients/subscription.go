package clients

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type SubscriptionManager struct {
	client *Client
	subs   []event.Subscription
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewSubscriptionManager(c *Client) *SubscriptionManager {
	return &SubscriptionManager{client: c}
}

func (m *SubscriptionManager) Close() {
	if m.cancel != nil {
		m.cancel()
	}
	for _, sub := range m.subs {
		sub.Unsubscribe()
	}
	m.wg.Wait()
}

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

// SubscribeNewHeads helper
func (m *SubscriptionManager) SubscribeNewHeads(handler func(*types.Header)) error {
	if !m.client.isWS {
		return fmt.Errorf("SubscribeNewHeads requires a WebSocket connection")
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	headers := make(chan *types.Header)
	sub, err := m.client.SubscribeNewHeads(ctx, headers)
	if err != nil {
		return err
	}
	m.subs = append(m.subs, sub)

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			select {
			case err := <-sub.Err():
				log.Println("NewHeads subscription error:", err)
				return
			case header := <-headers:
				handler(header)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

// SubscribeLogs helper
func (m *SubscriptionManager) SubscribeLogs(query ethereum.FilterQuery, handler func(types.Log)) error {
	if !m.client.isWS {
		return fmt.Errorf("SubscribeLogs requires a WebSocket connection")
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	logsCh := make(chan types.Log)
	sub, err := m.client.SubscribeLogs(ctx, query, logsCh)
	if err != nil {
		return err
	}
	m.subs = append(m.subs, sub)

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			select {
			case err := <-sub.Err():
				log.Println("Logs subscription error:", err)
				return
			case vLog := <-logsCh:
				handler(vLog)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

// SubscribePendingTxs Helper
func (m *SubscriptionManager) SubscribePendingTxs(handler func(common.Hash)) error {
	if !m.client.isWS {
		return fmt.Errorf("SubscribePendingTxs requires a WebSocket connection")
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	txCh := make(chan common.Hash)
	sub, err := m.client.SubscribePendingTxs(ctx, txCh)
	if err != nil {
		return err
	}
	m.subs = append(m.subs, sub)

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			select {
			case err := <-sub.Err():
				log.Println("PendingTxs subscription error:", err)
				return
			case tx := <-txCh:
				handler(tx)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
