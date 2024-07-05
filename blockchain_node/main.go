package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain Node: ")
}

func main() {
	// wMiner := wallet.NewWallet()
	// wAlice := wallet.NewWallet()
	// wBob := wallet.NewWallet()

	// // wallet tansaction request
	// t := wallet.NewTransaction(wAlice.PrivateKey(), wAlice.PublicKey(), wAlice.BlockchainAddress(), wBob.BlockchainAddress(), 23.0)

	// // blockchain node transaction request handling
	// blockchain := blockchain.NewBlockchain(wMiner.BlockchainAddress())
	// isAdded := blockchain.AddTransaction(wAlice.BlockchainAddress(), wBob.BlockchainAddress(), 23.0, wAlice.PublicKey(), t.GenerateSignature())

	// fmt.Println("Transaction add to transaction pool?", isAdded)

	// blockchain.Mining()
	// blockchain.Print()

	// fmt.Printf("Miner has %.1f\n", blockchain.CalculateTotalAmount(wMiner.BlockchainAddress()))
	// fmt.Printf("Alice has %.1f\n", blockchain.CalculateTotalAmount(wAlice.BlockchainAddress())) // XXX: should check total amount
	// fmt.Printf("Bob has %.1f\n", blockchain.CalculateTotalAmount(wBob.BlockchainAddress()))

	port := flag.Uint("port", 3333, "TCP Port Number for Blockchain Node")
	flag.Parse()

	app := NewBlockchainNode(uint16(*port))
	log.Default().Println("Starting blockchain node on port:", *port)
	app.Run()
}
