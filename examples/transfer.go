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

	wallet, err := clients.FromPrivateKey("YOUR_PRIVATE_KEY_HERE")
	if err != nil {
		log.Fatal(err)
	}

	recipient := common.HexToAddress("RECIPIENT_ADDRESS")
	amount := big.NewInt(0)
	amount.SetString("10000000000000000", 10) // 0.01 ETH

	fmt.Println("Sending 0.01 ETH to:", recipient.Hex())

	if err := wallet.SendETH(ctx, client, recipient, amount); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ ETH Transfer sent!")
}
