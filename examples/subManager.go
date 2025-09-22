package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	client, err := clients.DialWS("wss://api.testnet.abs.xyz/ws")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	manager := clients.NewSubscriptionManager(client)

	// Subscribe to new blocks
	err = manager.SubscribeNewHeads(func(h *types.Header) {
		fmt.Println("⛓ Block:", h.Number.Uint64())
	})
	if err != nil {
		return
	}

	// Subscribe to logs from a contract
	err = manager.SubscribeLogs(
		ethereum.FilterQuery{Addresses: []common.Address{common.HexToAddress("CONTRACT_ADDRESS")}},
		func(l types.Log) {
			fmt.Println("📜 Log event at block", l.BlockNumber)
		},
	)
	if err != nil {
		return
	}

	// Subscribe to pending txs
	err = manager.SubscribePendingTxs(func(tx common.Hash) {
		fmt.Println("📝 Pending tx:", tx.Hex())
	})
	if err != nil {
		return
	}

	// Keep running
	select {}
}
