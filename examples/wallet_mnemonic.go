package main

import (
	"fmt"
	"log"

	"github.com/mogza/abstract-go/clients"
)

func main() {
	mnemonic, err := clients.GenerateMnemonic(128) // 12 words
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("🧩 Generated mnemonic:", mnemonic)

	wallet, err := clients.NewWalletFromMnemonic(mnemonic, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Derived wallet address:", wallet.Address.Hex())
	fmt.Println("🔑 Private key:", wallet.ExportPrivateKeyHex())
}
