package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	w, err := clients.NewWallet()
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte("Hello Abstract!")

	sig, err := w.SignMessageEIP191(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("🖊 Signed message:", string(msg))
	fmt.Printf("📜 Signature (hex): 0x%x\n", sig)

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))
	hash := crypto.Keccak256([]byte(prefix), msg)

	addr, err := clients.RecoverAddressFromSignature(hash, sig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("🔍 Recovered address:", addr.Hex())
	fmt.Println("✅ Original address:", w.Address.Hex())
}
