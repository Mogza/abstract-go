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
	defer client.Close()

	wallet, err := clients.FromPrivateKey("YOUR_WALLET_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}
	nm := clients.NewNonceManager(client, wallet.Address)

	recipient := common.HexToAddress("RECIPIENT_ADDRESS")
	amount := big.NewInt(0)
	amount.SetString("10000000000000000", 10) // 0.01 ETH

	fmt.Println("Sending 0.01 ETH to:", recipient.Hex())

	tx, err := wallet.BuildAndSendTx(ctx, client, &recipient, amount, nil, nm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ ETH Transfer sent! Tx hash:", tx.Hash().Hex())
}
