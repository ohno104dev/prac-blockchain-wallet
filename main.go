package main

import (
	"fmt"

	"github.com/ohno104dev/prac-blockchain-wallet-go/wallet"
)

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
}
