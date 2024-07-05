package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ohno104dev/prac-blockchain-wallet-go/blockchain"
	"github.com/ohno104dev/prac-blockchain-wallet-go/wallet"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainNode struct {
	port uint16
}

func NewBlockchainNode(port uint16) *BlockchainNode {
	return &BlockchainNode{port}
}

func (bcn *BlockchainNode) Port() uint16 {
	return bcn.port
}

func (bcn *BlockchainNode) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minerWallet.BlockchainAddress(), bcn.Port())
		cache["blockchain"] = bc
		// XXX: DON'T DO THIS, only shown for demonstration purposes
		log.Printf("Public key %v", minerWallet.PublicKey())
		log.Printf("Private key %v", minerWallet.PrivateKey())
		log.Printf("Blockchain Address %v", minerWallet.BlockchainAddress())
	}

	return bc
}

func (bcn *BlockchainNode) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("Error: Invalid HTTP method")
	}
}

func (bcn *BlockchainNode) Run() {
	http.HandleFunc("/", bcn.GetChain)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.port)), nil))
}
