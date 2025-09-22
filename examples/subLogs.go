package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	client, err := clients.Dial("wss://api.testnet.abs.xyz")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Filter logs for a specific contract address
	contract := common.HexToAddress("CONTRACT_ADDRESS")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📡 Subscribed to logs for contract:", contract.Hex())

	for {
		select {
		case err := <-sub.Err():
			log.Println("Subscription error:", err)
			return
		case vLog := <-logsCh:
			fmt.Printf("📜 Log from %s | Block %d | Tx %s\n",
				vLog.Address.Hex(), vLog.BlockNumber, vLog.TxHash.Hex())
		}
	}
}
