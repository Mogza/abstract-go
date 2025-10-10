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

	recipients := []common.Address{
		common.HexToAddress("0xAddress1"),
		common.HexToAddress("0xAddress2"),
	}

	testAmount := big.NewInt(0)
	testAmount.SetString("10000000000000000", 10) // 0.01 ETH
	testAmount2 := big.NewInt(0)
	testAmount2.SetString("20000000000000000", 10) // 0.02 ETH

	amounts := []*big.Int{
		testAmount,  // 0.01 ETH
		testAmount2, // 0.02 ETH
	}

	txs, err := wallet.BatchSendETH(ctx, client, recipients, amounts, nm)
	if err != nil {
		fmt.Println("BatchSendETH failed:", err)
	} else {
		for i, tx := range txs {
			fmt.Printf("ETH tx %d: %s\n", i, tx.Hash().Hex())
		}
	}
}
