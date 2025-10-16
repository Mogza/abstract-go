package clients

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// NewWallet creates a new random Ethereum wallet.
// Returns a Wallet struct with a generated private key and address.
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

// FromPrivateKey creates a wallet from an existing hex-encoded private key.
// Returns a Wallet struct with the provided private key and derived address.
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

// GenerateMnemonic generates a new BIP39 mnemonic phrase with the given strength.
// Strength determines the number of words (e.g., 128 bits = 12 words).
func GenerateMnemonic(strength int) (string, error) {
	entropy, err := bip39.NewEntropy(strength)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// NewWalletFromMnemonic derives the first Ethereum account from a BIP39 mnemonic.
// Uses the standard path m/44'/60'/0'/0/0 and optional passphrase.
func NewWalletFromMnemonic(mnemonic, passphrase string) (*Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, passphrase)
	// master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	// Derivation path: m/44'/60'/0'/0/0
	// standard hardened derivation: index + 0x80000000
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coin, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	account, _ := coin.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)
	addrKey, err := change.NewChildKey(0)
	if err != nil {
		return nil, err
	}
	privKeyBytes := addrKey.Key
	priv, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: priv,
		Address:    crypto.PubkeyToAddress(priv.PublicKey),
	}, nil
}

// WalletFromHex creates a Wallet from a hex-encoded private key string.
// Strips 0x prefix if present and returns the Wallet struct.
func WalletFromHex(hexKey string) (*Wallet, error) {
	// strip 0x
	hexKey = strings.TrimPrefix(hexKey, "0x")
	priv, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: priv,
		Address:    crypto.PubkeyToAddress(priv.PublicKey),
	}, nil
}

// PrivateKeyHex returns the wallet's private key as a hex string (no 0x prefix).
// Returns an empty string if the wallet or private key is nil.
func (w *Wallet) PrivateKeyHex() string {
	if w == nil || w.PrivateKey == nil {
		return ""
	}
	return hex.EncodeToString(crypto.FromECDSA(w.PrivateKey))
}

// ImportKeystoreJSON imports a go-ethereum keystore JSON and decrypts it.
// Returns a Wallet struct with the decrypted private key and address.
func ImportKeystoreJSON(keyjson []byte, password string) (*Wallet, error) {
	keyStruct, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: keyStruct.PrivateKey,
		Address:    keyStruct.Address,
	}, nil
}

// RecoverAddressFromSignature recovers the Ethereum address from a signature and digest.
// The digest must be the original 32-byte hash that was signed.
func RecoverAddressFromSignature(digest []byte, sig []byte) (common.Address, error) {
	var zero common.Address
	if len(sig) != 65 {
		return zero, errors.New("signature must be 65 bytes")
	}

	sigCopy := make([]byte, 65)
	copy(sigCopy, sig)

	// Ensure v is in {0,1} for go-ethereum
	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}

	pubkey, err := crypto.SigToPub(digest, sigCopy)
	if err != nil {
		return zero, err
	}

	return crypto.PubkeyToAddress(*pubkey), nil
}

// VerifySignature checks that signature was produced by expected address for the given digest.
// Returns true if the recovered address matches the expected address.
func VerifySignature(digest []byte, sig []byte, expected common.Address) (bool, error) {
	addr, err := RecoverAddressFromSignature(digest, sig)
	if err != nil {
		return false, err
	}
	return addr == expected, nil
}

// NewDeterministicWallet returns a wallet derived from a seed phrase + index using the ETH derivation path.
// It is deterministic and useful for tests/dev.
func NewDeterministicWallet(mnemonic string, index uint32) (*Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	// m/44'/60'/0'/0/index
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coin, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	account, _ := coin.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)
	addrKey, err := change.NewChildKey(index)
	if err != nil {
		return nil, err
	}
	priv, err := crypto.ToECDSA(addrKey.Key)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: priv,
		Address:    crypto.PubkeyToAddress(priv.PublicKey),
	}, nil
}
