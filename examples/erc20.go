package main

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	ctx := context.Background()

	client, _ := clients.DialHTTP("https://api.testnet.abs.xyz")
	defer client.Eth.Close()

	wallet, _ := clients.FromPrivateKey("YOUR_PRIVATE_KEY")
	token := common.HexToAddress("TOKEN_ADDRESS")
	recipient := common.HexToAddress("RECIPIENT_ADDRESS")

	erc20, _ := clients.NewERC20(client, token, "")

	balance, _ := erc20.BalanceOf(ctx, wallet.Address)
	name, _ := erc20.Name(ctx)
	symbol, _ := erc20.Symbol(ctx)
	decimals, _ := erc20.Decimals(ctx)

	fmt.Printf("%s (%s) Decimals:%d Balance:%s\n", name, symbol, decimals, balance.String())

	amtFloat := new(big.Float).SetFloat64(0.001)
	decimalsBig := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	amtFloat.Mul(amtFloat, decimalsBig)
	amount := new(big.Int)
	amtFloat.Int(amount)

	tx, _ := erc20.Transfer(ctx, wallet, recipient, amount)
	fmt.Println("ERC20 Transfer Tx Hash:", tx.Hash().Hex())
}
