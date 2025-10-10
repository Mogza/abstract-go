package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	ctx := context.Background()

	wsClient, _ := clients.DialWS("wss://api.testnet.abs.xyz/ws")
	defer wsClient.Close()

	token := common.HexToAddress("ERC20_TOKEN_ADDRESS")
	erc20, _ := clients.NewERC20(wsClient, token, "")

	transferCh := make(chan clients.ERC20TransferEvent)
	approvalCh := make(chan clients.ERC20ApprovalEvent)

	erc20.WatchTransfers(ctx, nil, nil, transferCh)
	erc20.WatchApprovals(ctx, nil, nil, approvalCh)

	for {
		select {
		case t := <-transferCh:
			fmt.Printf("🔁 Transfer: %s -> %s : %s tokens\n", t.From.Hex(), t.To.Hex(), t.Value.String())
		case a := <-approvalCh:
			fmt.Printf("✅ Approval: %s approved %s : %s tokens\n", a.Owner.Hex(), a.Spender.Hex(), a.Value.String())
		}
	}
}
