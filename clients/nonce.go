package clients

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type NonceManager struct {
	client *Client
	addr   common.Address
	mu     sync.Mutex
	nonce  uint64
	init   bool
}

// NewNonceManager creates a nonce manager for a wallet
func NewNonceManager(client *Client, addr common.Address) *NonceManager {
	return &NonceManager{
		client: client,
		addr:   addr,
	}
}

// Next returns the next nonce safely
func (nm *NonceManager) Next(ctx context.Context) (uint64, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if !nm.init {
		n, err := nm.client.NonceAt(ctx, nm.addr)
		if err != nil {
			return 0, err
		}
		nm.nonce = n
		nm.init = true
	}

	nonce := nm.nonce
	nm.nonce++
	return nonce, nil
}

// Reset allows forcing sync with blockchain (if needed)
func (nm *NonceManager) Reset() {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.init = false
}
