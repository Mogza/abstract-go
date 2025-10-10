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
	defer client.Eth.Close()

	wallet, err := clients.FromPrivateKey("YOUR_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}

	nftAddr := common.HexToAddress("NFT_CONTRACT_ADDRESS")
	nft, err := clients.NewERC721(client, nftAddr, "")
	if err != nil {
		log.Fatal(err)
	}

	tokenID := big.NewInt(1)

	balance, err := nft.BalanceOf(ctx, wallet.Address)
	if err != nil {
		log.Fatal(err)
	}
	owner, err := nft.OwnerOf(ctx, tokenID)
	if err != nil {
		log.Fatal(err)
	}
	uri, err := nft.TokenURI(ctx, tokenID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NFT Balance:", balance)
	fmt.Println("Token Owner:", owner.Hex())
	fmt.Println("Token URI:", uri)

	recipient := common.HexToAddress("RECIPIENT_ADDRESS")
	tx, err := nft.TransferFrom(ctx, wallet, wallet.Address, recipient, tokenID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transfer Tx Hash:", tx.Hash().Hex())
}
