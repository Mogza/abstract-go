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
	wsClient, err := clients.DialWS("wss://api.testnet.abs.xyz/ws")
	if err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()

	nftAddr := common.HexToAddress("NFT_CONTRACT_ADDRESS")
	nft, err := clients.NewERC721(wsClient, nftAddr, "")
	if err != nil {
		log.Fatal(err)
	}

	transferCh := make(chan clients.ERC721TransferEvent)
	approvalCh := make(chan clients.ERC721ApprovalEvent)

	if err := nft.WatchTransfers(ctx, transferCh); err != nil {
		log.Fatal(err)
	}
	if err := nft.WatchApprovals(ctx, approvalCh); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case t := <-transferCh:
			fmt.Printf("🔁 NFT Transfer: %s -> %s tokenID:%s\n", t.From.Hex(), t.To.Hex(), t.TokenID.String())
		case a := <-approvalCh:
			fmt.Printf("✅ NFT Approval: %s approved %s tokenID:%s\n", a.Owner.Hex(), a.Approved.Hex(), a.TokenID.String())
		}
	}
}
