package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	// Connect to node via WebSocket
	wsClient, err := clients.DialWS("wss://api.testnet.abs.xyz/ws")
	if err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()

	// Token address to watch
	token := common.HexToAddress("0xYourERC20Token")

	// Channels to receive events
	transferCh := make(chan clients.ERC20TransferEvent)
	approvalCh := make(chan clients.ERC20ApprovalEvent)

	// Start watching Transfer events
	err = clients.WatchERC20Transfers(wsClient, token, nil, nil, transferCh)
	if err != nil {
		log.Fatal(err)
	}

	// Start watching Approval events
	err = clients.WatchERC20Approvals(wsClient, token, nil, nil, approvalCh)
	if err != nil {
		log.Fatal(err)
	}

	// Listen and print events
	for {
		select {
		case t := <-transferCh:
			fmt.Printf("🔁 Transfer: %s -> %s : %s tokens\n", t.From.Hex(), t.To.Hex(), t.Value.String())
		case a := <-approvalCh:
			fmt.Printf("✅ Approval: %s approved %s : %s tokens\n", a.Owner.Hex(), a.Spender.Hex(), a.Value.String())
		}
	}
}
