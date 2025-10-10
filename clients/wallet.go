package clients

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

// BuildAndSendTx signs and sends ETH to a recipient using EIP-1559 with gas estimation
func (w *Wallet) BuildAndSendTx(ctx context.Context, client *Client, to *common.Address, value *big.Int, data []byte, nm *NonceManager) (*types.Transaction, error) {
	if client.isWS {
		return nil, fmt.Errorf("BuildAndSendTx requires an HTTP connection, not WebSocket")
	}

	// Get next nonce safely
	nonce, err := nm.Next(ctx)
	if err != nil {
		return nil, err
	}

	// Gas suggestion
	gasTipCap, err := client.Eth.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	baseFee, err := client.Eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	maxFee := new(big.Int).Add(baseFee, gasTipCap)

	// Estimate gas with optional buffer
	msg := ethereum.CallMsg{
		From:  w.Address,
		To:    to,
		Value: value,
		Data:  data,
	}
	gasLimit, err := client.EstimateGasWithBuffer(ctx, msg, 10) // +10% buffer
	if err != nil {
		return nil, err
	}

	chainID, err := client.Eth.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: maxFee,
		Gas:       gasLimit,
		To:        to,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), w.PrivateKey)
	if err != nil {
		return nil, err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, err
	}

	return signedTx, nil
}

// ExportKeystoreJSON exports the wallet as an encrypted keystore JSON (go-ethereum format).
// Use a strong password. Returns the JSON bytes.
func (w *Wallet) ExportKeystoreJSON(password string) ([]byte, error) {
	if w == nil || w.PrivateKey == nil {
		return nil, errors.New("wallet or private key nil")
	}
	id := uuid.New()
	keyStruct := &keystore.Key{
		Id:         id,
		Address:    w.Address,
		PrivateKey: w.PrivateKey,
	}
	// standard scrypt params
	keyjson, err := keystore.EncryptKey(keyStruct, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	return keyjson, nil
}

// SignMessageEIP191 signs a human-readable message with the `\x19Ethereum Signed Message:\n<len>` prefix.
// Returns the 65-byte [R|S|V] signature (where V is 27/28 converted to 0/1 for go-ethereum compat).
func (w *Wallet) SignMessageEIP191(message []byte) ([]byte, error) {
	if w == nil || w.PrivateKey == nil {
		return nil, errors.New("wallet or private key nil")
	}
	// prefix as in personal_sign / EIP-191
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	hash := crypto.Keccak256([]byte(prefix), message)
	sig, err := crypto.Sign(hash, w.PrivateKey)
	if err != nil {
		return nil, err
	}
	// Normalize V to 27/28 (Ethereum standard)
	if sig[64] < 27 {
		sig[64] += 27
	}
	return sig, nil
}

// SignHash signs a 32-byte digest using the account private key (generic).
// Use this for signing EIP-712 hashes (you provide the typed data hash).
func (w *Wallet) SignHash(digest []byte) ([]byte, error) {
	if w == nil || w.PrivateKey == nil {
		return nil, errors.New("wallet or private key nil")
	}
	if len(digest) != 32 {
		// allow if user passes non-32 but still hashable: we hash it to 32 bytes
		tmp := sha256.Sum256(digest)
		digest = tmp[:]
	}
	sig, err := crypto.Sign(digest, w.PrivateKey)
	if err != nil {
		return nil, err
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	return sig, nil
}

// ExportPrivateKeyHex is a convenience wrapper returning the 0x-prefixed hex private key
func (w *Wallet) ExportPrivateKeyHex() string {
	return "0x" + w.PrivateKeyHex()
}
