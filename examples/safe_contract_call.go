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

	nm := clients.NewNonceManager(client, wallet.Address)

	contract := common.HexToAddress("0xContractAddress")
	abiJSON := `[{"inputs":[{"internalType":"uint256","name":"value","type":"uint256"}],"name":"setValue","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	tx, err := wallet.SafeContractCall(ctx, client, contract, abiJSON, "setValue", nm, big.NewInt(42))
	if err != nil {
		fmt.Println("SafeContractCall failed:", err)
	} else {
		fmt.Println("SafeContractCall tx:", tx.Hash().Hex())
	}
}
