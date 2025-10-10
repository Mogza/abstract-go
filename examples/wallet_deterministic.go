package main

import (
	"fmt"
	"log"

	"github.com/mogza/abstract-go/clients"
)

func main() {
	mnemonic, _ := clients.GenerateMnemonic(128) // 12 words
	fmt.Println("🌱 Using fixed mnemonic:", mnemonic)

	for i := uint32(0); i < 3; i++ {
		w, err := clients.NewDeterministicWallet(mnemonic, i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Wallet #%d → %s\n", i, w.Address.Hex())
	}
}
