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
	fmt.Println("👛 Wallet address:", wallet.Address.Hex())

	jsonBytes, err := wallet.ExportKeystoreJSON("YOUR_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("📦 Exported keystore JSON (truncated):", string(jsonBytes)[:80], "...")

	imported, err := clients.ImportKeystoreJSON(jsonBytes, "YOUR_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Imported wallet address:", imported.Address.Hex())
}
