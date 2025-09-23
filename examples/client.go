package main

import (
	"context"
	"fmt"
	"log"

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

	address := common.HexToAddress("0x0000000000000000000000000000000000000000")
	balance, err := client.BalanceAt(ctx, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ETH Balance:", balance.String())
}
