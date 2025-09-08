package clients

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const erc20ABI = `[{
	"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"
},{
	"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"
},{
	"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"
},{
	"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"
},{
	"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"
},{
	"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"
},{
	"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"
}]`

var parsedERC20ABI, _ = abi.JSON(strings.NewReader(erc20ABI))

// --- Read-only calls ---

func ERC20BalanceOf(ctx context.Context, client *Client, token common.Address, owner common.Address) (*big.Int, error) {
	data, err := parsedERC20ABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, err
	}
	result, err := client.Eth.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	values, err := parsedERC20ABI.Unpack("balanceOf", result)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func ERC20Allowance(ctx context.Context, client *Client, token common.Address, owner, spender common.Address) (*big.Int, error) {
	data, err := parsedERC20ABI.Pack("allowance", owner, spender)
	if err != nil {
		return nil, err
	}
	result, err := client.Eth.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	values, err := parsedERC20ABI.Unpack("allowance", result)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func ERC20Name(ctx context.Context, client *Client, token common.Address) (string, error) {
	data, err := parsedERC20ABI.Pack("name")
	if err != nil {
		return "", err
	}
	result, err := client.Eth.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return "", err
	}
	values, err := parsedERC20ABI.Unpack("name", result)
	if err != nil {
		return "", err
	}
	return values[0].(string), nil
}

func ERC20Symbol(ctx context.Context, client *Client, token common.Address) (string, error) {
	data, err := parsedERC20ABI.Pack("symbol")
	if err != nil {
		return "", err
	}
	result, err := client.Eth.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return "", err
	}
	values, err := parsedERC20ABI.Unpack("symbol", result)
	if err != nil {
		return "", err
	}
	return values[0].(string), nil
}

func ERC20Decimals(ctx context.Context, client *Client, token common.Address) (uint8, error) {
	data, err := parsedERC20ABI.Pack("decimals")
	if err != nil {
		return 0, err
	}
	result, err := client.Eth.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return 0, err
	}
	values, err := parsedERC20ABI.Unpack("decimals", result)
	if err != nil {
		return 0, err
	}
	return values[0].(uint8), nil
}

// --- Write transactions (signed) ---

func ERC20Transfer(ctx context.Context, wallet *Wallet, client *Client, token, to common.Address, amount *big.Int) (*types.Transaction, error) {
	data, err := parsedERC20ABI.Pack("transfer", to, amount)
	if err != nil {
		return nil, err
	}
	return sendERC20Tx(ctx, wallet, client, token, data)
}

func ERC20Approve(ctx context.Context, wallet *Wallet, client *Client, token, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	data, err := parsedERC20ABI.Pack("approve", spender, amount)
	if err != nil {
		return nil, err
	}
	return sendERC20Tx(ctx, wallet, client, token, data)
}

// --- Internal helper for sending signed ERC20 tx ---
func sendERC20Tx(ctx context.Context, wallet *Wallet, client *Client, token common.Address, data []byte) (*types.Transaction, error) {
	nonce, err := client.NonceAt(ctx, wallet.Address)
	if err != nil {
		return nil, err
	}

	gasTipCap, err := client.Eth.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	gasFeeCap, err := client.Eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: wallet.Address,
		To:   &token,
		Data: data,
	}
	gasLimit, err := client.Eth.EstimateGas(ctx, msg)
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &token,
		Gas:       gasLimit,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Value:     big.NewInt(0), // ERC20 transfers use data, not ETH
		Data:      data,
	})

	chainID, err := client.Eth.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, err
	}

	return signedTx, nil
}
