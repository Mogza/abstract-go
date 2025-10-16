package clients

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

// NewSubscriptionManager creates a new SubscriptionManager for managing event subscriptions.
// It holds references to the client, subscriptions, and manages goroutines.
func NewSubscriptionManager(c *Client) *SubscriptionManager {
	return &SubscriptionManager{client: c}
}

// Close cancels all active subscriptions and waits for goroutines to finish.
// Ensures a clean shutdown of the SubscriptionManager.
func (m *SubscriptionManager) Close() {
	if m.cancel != nil {
		m.cancel()
	}
	for _, sub := range m.subs {
		sub.Unsubscribe()
	}
	m.wg.Wait()
}

// SubscribeNewHeads subscribes to new block headers as soon as blocks are mined.
// Requires a WebSocket connection; sends headers to the provided channel.
func (c *Client) SubscribeNewHeads(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	if !c.isWS {
		return nil, fmt.Errorf("SubscribeNewHeads requires a WebSocket connection")
	}

	return c.Eth.SubscribeNewHead(ctx, ch)
}

// SubscribeLogs subscribes to smart contract event logs matching the given filter query.
// Requires a WebSocket connection; sends logs to the provided channel.
func (c *Client) SubscribeLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	if !c.isWS {
		return nil, fmt.Errorf("SubscribeLogs requires a WebSocket connection")
	}
	return c.Eth.SubscribeFilterLogs(ctx, query, ch)
}

// SubscribePendingTxs subscribes to new transactions entering the mempool.
// Requires a WebSocket connection; sends transaction hashes to the provided channel.
func (c *Client) SubscribePendingTxs(ctx context.Context, ch chan<- common.Hash) (ethereum.Subscription, error) {
	if !c.isWS {
		return nil, fmt.Errorf("SubscribePendingTxs requires a WebSocket connection")
	}
	return c.RpcClient.EthSubscribe(ctx, ch, "newPendingTransactions")
}

// --- Unified Event Watching ---

type EventHandler func(vLog types.Log) error

// WatchContractEvent watches for a specific contract event with optional indexed filters.
// Invokes the handler for each matching log; requires a WebSocket client.
func (c *Client) WatchContractEvent(ctx context.Context, contractAddr common.Address, abiJSON string, eventName string, filter map[string][]common.Address, handler EventHandler) error {
	if !c.isWS {
		return fmt.Errorf("WatchContractEvent requires a WS client")
	}

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return err
	}

	retEvent, exists := parsedABI.Events[eventName]
	if !exists {
		return fmt.Errorf("event %s not found in ABI", eventName)
	}

	topics := [][]common.Hash{{retEvent.ID}}

	// Apply indexed filters if provided
	for i, input := range retEvent.Inputs {
		if !input.Indexed {
			continue
		}
		key := input.Name
		if addrs, ok := filter[key]; ok && len(addrs) > 0 {
			topicHashes := make([]common.Hash, len(addrs))
			for j, addr := range addrs {
				topicHashes[j] = common.HexToHash(addr.Hex())
			}
			// fill empty topics until this index
			for len(topics) <= i {
				topics = append(topics, nil)
			}
			topics[i] = topicHashes
		}
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics:    topics,
	}

	logsCh := make(chan types.Log)
	sub, err := c.SubscribeLogs(ctx, query, logsCh)
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
				if err := handler(vLog); err != nil {
					fmt.Println("handler error:", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// SubscribeNewHeads subscribes to new block headers and invokes the handler for each.
// Manages subscription lifecycle and goroutine cleanup.
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

// SubscribeLogs subscribes to contract logs matching the filter and invokes the handler.
// Manages subscription lifecycle and goroutine cleanup.
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

//  --- SubscriptionManager : Helper ---

// SubscribePendingTxs subscribes to new pending transactions and invokes the handler for each.
// Manages subscription lifecycle and goroutine cleanup.
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
