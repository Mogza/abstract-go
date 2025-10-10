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
	client, _ := clients.DialHTTP("https://api.testnet.abs.xyz")
	defer client.Eth.Close()

	wallet, err := clients.FromPrivateKey("YOUR_WALLET_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}

	erc721Addr := common.HexToAddress("YOUR_ERC721_CONTRACT_ADDRESS")
	nft, _ := clients.NewERC721(client, erc721Addr, "")

	owner := wallet.Address
	tokenId := big.NewInt(1)

	balance, _ := nft.BalanceOf(ctx, owner)
	fmt.Println("NFT Balance:", balance)

	ownerAddr, _ := nft.OwnerOf(ctx, tokenId)
	fmt.Println("Token owner:", ownerAddr.Hex())

	uri, _ := nft.TokenURI(ctx, tokenId)
	fmt.Println("Token URI:", uri)

	tx, _ := nft.TransferFrom(ctx, wallet, owner, common.HexToAddress("RECIPIENT"), tokenId)
	fmt.Println("Transfer tx hash:", tx.Hash().Hex())
}
