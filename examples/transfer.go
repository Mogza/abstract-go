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

	// Connect to Abstract RPC
	client, err := clients.Dial("https://api.testnet.abs.xyz")
	if err != nil {
		log.Fatal("Failed to connect to Abstract RPC:", err)
	}
	defer client.Eth.Close()

	// Create a new wallet (or use FromPrivateKey)
	wallet, err := clients.NewWallet()
	if err != nil {
		log.Fatal("Failed to create wallet:", err)
	}

	fmt.Println("Wallet Address:", wallet.Address.Hex())

	// Example recipient (replace with a real one)
	recipient := common.HexToAddress("0x0000000000000000000000000000000000000000")

	// Amount to send (0.01 ETH)
	amount := big.NewInt(0)
	amount.SetString("10000000000000000", 10) // 0.01 ETH in wei

	fmt.Println("Sending 0.01 ETH to:", recipient.Hex())

	// Send ETH
	err = wallet.SendETH(ctx, client, recipient, amount)
	if err != nil {
		log.Fatal("Failed to send ETH:", err)
	}

	fmt.Println("✅ Transaction sent successfully!")
}
