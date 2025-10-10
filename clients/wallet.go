package clients

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func (w *Wallet) ApproveAndTransferERC20(ctx context.Context, client *Client, token, recipient, spender common.Address, amount *big.Int, nm *NonceManager) (*types.Transaction, *types.Transaction, error) {

	erc20Token, err := NewERC20(client, token, "")
	if err != nil {
		return nil, nil, err
	}

	// Approve spender
	approveTx, err := erc20Token.Approve(ctx, w, spender, amount)
	if err != nil {
		return nil, nil, err
	}

	// Transfer tokens from spender to recipient
	transferTx, err := erc20Token.TransferFrom(ctx, w, spender, recipient, amount)
	if err != nil {
		return approveTx, nil, err
	}

	return approveTx, transferTx, nil
}

func (w *Wallet) BatchSendETH(ctx context.Context, client *Client, recipients []common.Address, amounts []*big.Int, nm *NonceManager) ([]*types.Transaction, error) {

	if len(recipients) != len(amounts) {
		return nil, fmt.Errorf("recipients and amounts length mismatch")
	}

	var txs []*types.Transaction

	for i, to := range recipients {
		tx, err := w.BuildAndSendTx(ctx, client, &to, amounts[i], nil, nm)
		if err != nil {
			return txs, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

func (w *Wallet) SafeContractCall(ctx context.Context, client *Client, contract common.Address, abiJSON string, method string, nm *NonceManager, params ...interface{}) (*types.Transaction, error) {

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	// Optional: simulate call first
	msg := ethereum.CallMsg{
		From: w.Address,
		To:   &contract,
		Data: data,
	}
	if _, err := client.CallContract(ctx, msg); err != nil {
		return nil, fmt.Errorf("simulation failed: %w", err)
	}

	tx, err := w.BuildAndSendTx(ctx, client, &contract, big.NewInt(0), data, nm)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
