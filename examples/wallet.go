package main

import (
	"fmt"
	"log"

	"github.com/mogza/abstract-go/clients"
)

func main() {
	wallet, err := clients.NewWallet()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New Wallet Created:")
	fmt.Println("Address:", wallet.Address.Hex())
	fmt.Println("Private Key (KEEP SECRET):", wallet.PrivateKey.D.Text(16))
}
