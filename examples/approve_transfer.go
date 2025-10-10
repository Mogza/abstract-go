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

	token := common.HexToAddress("0xYourTokenAddress")
	recipient := common.HexToAddress("0xRecipientAddress")
	spender := common.HexToAddress("0xSpenderAddress")
	amount := big.NewInt(0)
	amount.SetString("10000000000000000", 10) // 0.01 ETH

	approveTx, transferTx, err := wallet.ApproveAndTransferERC20(ctx, client, token, recipient, spender, amount, nm)
	if err != nil {
		fmt.Println("ApproveAndTransferERC20 failed:", err)
	} else {
		fmt.Println("Approve tx:", approveTx.Hash().Hex())
		fmt.Println("Transfer tx:", transferTx.Hash().Hex())
	}
}
