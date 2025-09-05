package clients

import (
	"context"
	"crypto/ecdsa"
	"math/big"

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

// SendETH signs and sends ETH to a recipient using EIP-1559
func (w *Wallet) SendETH(ctx context.Context, client *Client, to common.Address, amount *big.Int) error {
	nonce, err := client.NonceAt(ctx, w.Address)
	if err != nil {
		return err
	}

	gasTipCap, err := client.Eth.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	gasFeeCap, err := client.Eth.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   nil, // set later with signer
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       21000,
		To:        &to,
		Value:     amount,
		Data:      nil,
	})

	chainID, err := client.Eth.NetworkID(ctx)
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), w.PrivateKey)
	if err != nil {
		return err
	}

	return client.SendTransaction(ctx, signedTx)
}
