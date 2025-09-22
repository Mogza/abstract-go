package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	client, err := clients.DialWS("wss://api.testnet.abs.xyz/ws")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	pendingCh := make(chan common.Hash)
	sub, err := client.SubscribePendingTxs(context.Background(), pendingCh)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📡 Subscribed to pending transactions")

	for {
		select {
		case err := <-sub.Err():
			log.Println("Subscription error:", err)
			return
		case txHash := <-pendingCh:
			fmt.Println("📝 Pending tx:", txHash.Hex())
		}
	}
}
