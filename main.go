package main

import (
	"fmt"

	"github.com/ohno104dev/prac-blockchain-wallet-go/wallet"
)

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "Morty", 137.0)
	fmt.Printf("Signature %s", t.GenerateSignature())
}
