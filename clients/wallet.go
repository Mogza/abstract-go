package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

// NewWallet creates a new random wallet
func NewWallet() (*Wallet, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: key,
		Address:    crypto.PubkeyToAddress(key.PublicKey),
	}, nil
}

// FromPrivateKey creates wallet from an existing private key
func FromPrivateKey(hexKey string) (*Wallet, error) {
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: key,
		Address:    crypto.PubkeyToAddress(key.PublicKey),
	}, nil
}

// SendETH signs and sends ETH to a recipient using EIP-1559 with gas estimation
func (w *Wallet) SendETH(ctx context.Context, client *Client, to common.Address, amount *big.Int) error {
	if client.isWS {
		return fmt.Errorf("SendETH requires an HTTP connection, not WebSocket")
	}

	// Get account nonce
	nonce, err := client.NonceAt(ctx, w.Address)
	if err != nil {
		return err
	}

	gasTipCap, err := client.Eth.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	baseFee, err := client.Eth.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	maxFee := new(big.Int).Add(baseFee, gasTipCap)

	// Estimate gas limit dynamically
	msg := ethereum.CallMsg{
		From:  w.Address,
		To:    &to,
		Value: amount,
		Data:  nil,
	}
	gasLimit, err := client.Eth.EstimateGas(ctx, msg)
	if err != nil {
		return err
	}

	// Get chain ID
	chainID, err := client.Eth.NetworkID(ctx)
	if err != nil {
		return err
	}

	// Create EIP-1559 tx
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: maxFee, // Using same as tip if base fee not available; can also fetch base fee if RPC supports
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      nil,
	})

	// Sign tx
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), w.PrivateKey)
	if err != nil {
		return err
	}

	// Send tx
	return client.SendTransaction(ctx, signedTx)
}
