package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	client, err := clients.Dial("wss://api.testnet.abs.xyz/ws")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHeads(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📡 Subscribed to new heads")

	for {
		select {
		case err := <-sub.Err():
			log.Println("Subscription error:", err)
			return
		case header := <-headers:
			fmt.Println("⛓ New block:", header.Number, "at", time.Unix(int64(header.Time), 0))
		}
	}
}
