package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	ctx := context.Background()
	client, _ := clients.DialHTTP("https://api.testnet.abs.xyz")
	defer client.Eth.Close()

	wallet, _ := clients.FromPrivateKey("YOUR_PRIVATE_KEY")
	nftAddr := common.HexToAddress("NFT_CONTRACT_ADDRESS")
	nft, _ := clients.NewERC721(client, nftAddr, "")

	tokenID := big.NewInt(1)

	balance, _ := nft.BalanceOf(ctx, wallet.Address)
	owner, _ := nft.OwnerOf(ctx, tokenID)
	uri, _ := nft.TokenURI(ctx, tokenID)

	fmt.Println("NFT Balance:", balance)
	fmt.Println("Token Owner:", owner.Hex())
	fmt.Println("Token URI:", uri)

	recipient := common.HexToAddress("RECIPIENT_ADDRESS")
	tx, _ := nft.TransferFrom(ctx, wallet, wallet.Address, recipient, tokenID)
	fmt.Println("Transfer Tx Hash:", tx.Hash().Hex())
}
