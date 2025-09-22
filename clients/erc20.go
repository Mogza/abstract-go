package clients

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type ERC20ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

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

func parseTransferLog(vLog types.Log) (*ERC20TransferEvent, error) {
	event := new(ERC20TransferEvent)
	err := parsedERC20ABI.UnpackIntoInterface(event, "Transfer", vLog.Data)
	if err != nil {
		return nil, err
	}
	// Parse indexed fields
	event.From = common.HexToAddress(vLog.Topics[1].Hex())
	event.To = common.HexToAddress(vLog.Topics[2].Hex())
	return event, nil
}

func parseApprovalLog(vLog types.Log) (*ERC20ApprovalEvent, error) {
	event := new(ERC20ApprovalEvent)
	err := parsedERC20ABI.UnpackIntoInterface(event, "Approval", vLog.Data)
	if err != nil {
		return nil, err
	}
	event.Owner = common.HexToAddress(vLog.Topics[1].Hex())
	event.Spender = common.HexToAddress(vLog.Topics[2].Hex())
	return event, nil
}

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

//  --- Watch Functions ---

func WatchERC20Transfers(client *Client, token common.Address, from, to *common.Address, ch chan<- ERC20TransferEvent) error {
	if !client.isWS {
		return fmt.Errorf("WatchERC20Transfers requires a WS client")
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{token},
	}
	if from != nil {
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(from.Hex())})
	}
	if to != nil {
		if len(query.Topics) == 0 {
			query.Topics = append(query.Topics, nil)
		}
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(to.Hex())})
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeLogs(context.Background(), query, logsCh)
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
				event, err := parseTransferLog(vLog)
				if err != nil {
					continue
				}
				ch <- *event
			}
		}
	}()

	return nil
}

func WatchERC20Approvals(client *Client, token common.Address, owner, spender *common.Address, ch chan<- ERC20ApprovalEvent) error {
	if !client.isWS {
		return fmt.Errorf("WatchERC20Approvals requires a WS client")
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{token},
	}
	if owner != nil {
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(owner.Hex())})
	}
	if spender != nil {
		if len(query.Topics) == 0 {
			query.Topics = append(query.Topics, nil)
		}
		query.Topics = append(query.Topics, []common.Hash{common.HexToHash(spender.Hex())})
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeLogs(context.Background(), query, logsCh)
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
				event, err := parseApprovalLog(vLog)
				if err != nil {
					continue
				}
				ch <- *event
			}
		}
	}()

	return nil
}
