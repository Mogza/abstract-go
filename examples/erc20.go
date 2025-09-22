package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	ctx := context.Background()

	client, err := clients.DialHTTP("https://api.testnet.abs.xyz")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Eth.Close()

	// Use your private key
	wallet, err := clients.FromPrivateKey("YOUR_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}

	token := common.HexToAddress("TOKEN_ADDRESS")
	recipient := common.HexToAddress("RECIPIENT_ADDRESS")

	// Read ERC20 info
	balance, _ := clients.ERC20BalanceOf(ctx, client, token, wallet.Address)
	name, _ := clients.ERC20Name(ctx, client, token)
	symbol, _ := clients.ERC20Symbol(ctx, client, token)
	decimals, _ := clients.ERC20Decimals(ctx, client, token)

	fmt.Printf("Token: %s (%s) Decimals: %d\nBalance: %s\n", name, symbol, decimals, balance.String())

	amtFloat := new(big.Float).SetFloat64(0.001)
	decimalsBig := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	amtFloat.Mul(amtFloat, decimalsBig)
	amount := new(big.Int)
	amtFloat.Int(amount)

	tx, err := clients.ERC20Transfer(ctx, wallet, client, token, recipient, amount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ERC20 Transfer Tx Hash:", tx.Hash().Hex())
}
