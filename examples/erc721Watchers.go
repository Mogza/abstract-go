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

	nftAddr := common.HexToAddress("NFT_CONTRACT_ADDRESS")
	nft, _ := clients.NewERC721(wsClient, nftAddr, "")

	transferCh := make(chan clients.ERC721TransferEvent)
	approvalCh := make(chan clients.ERC721ApprovalEvent)

	nft.WatchTransfers(ctx, transferCh)
	nft.WatchApprovals(ctx, approvalCh)

	for {
		select {
		case t := <-transferCh:
			fmt.Printf("🔁 NFT Transfer: %s -> %s tokenID:%s\n", t.From.Hex(), t.To.Hex(), t.TokenID.String())
		case a := <-approvalCh:
			fmt.Printf("✅ NFT Approval: %s approved %s tokenID:%s\n", a.Owner.Hex(), a.Approved.Hex(), a.TokenID.String())
		}
	}
}
