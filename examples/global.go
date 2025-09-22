package main

import (
	"context"
	"fmt"
	"log"
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

	// Create wallet
	wallet, err := clients.NewWallet()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wallet Address:", wallet.Address.Hex())

	// Send ETH
	recipient := common.HexToAddress("RECIPIENT_ADDRESS")
	ethAmount := big.NewInt(0)
	ethAmount.SetString("5000000000000000", 10) // 0.005 ETH
	if err := wallet.SendETH(ctx, client, recipient, ethAmount); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ ETH sent!")

	// ERC20 transfer
	token := common.HexToAddress("TOKEN_ADDRESS")
	erc20Amount := big.NewInt(1000) // Example
	tx, err := clients.ERC20Transfer(ctx, wallet, client, token, recipient, erc20Amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ ERC20 sent, Tx Hash:", tx.Hash().Hex())

	// Read ERC20 balance
	balance, err := clients.ERC20BalanceOf(ctx, client, token, wallet.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ERC20 Balance:", balance.String())
}
