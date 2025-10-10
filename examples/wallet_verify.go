package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	w, _ := clients.NewWallet()
	msg := []byte("Verify this message!")

	sig, _ := w.SignMessageEIP191(msg)
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))
	hash := crypto.Keccak256([]byte(prefix), msg)

	valid, err := clients.VerifySignature(hash, sig, w.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("🧾 Signature valid?", valid)
}
