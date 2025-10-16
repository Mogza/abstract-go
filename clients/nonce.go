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

// NewNonceManager creates a NonceManager for safely tracking and incrementing nonces.
// Associates the manager with a specific client and address.
func NewNonceManager(client *Client, addr common.Address) *NonceManager {
	return &NonceManager{
		client: client,
		addr:   addr,
	}
}

// Next returns the next available nonce for the address, safely incrementing it.
// Initializes from the blockchain if not already synced.
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

// Reset forces the NonceManager to resync the nonce from the blockchain on next use.
// Useful if a transaction is dropped or nonce is out of sync.
func (nm *NonceManager) Reset() {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.init = false
}
