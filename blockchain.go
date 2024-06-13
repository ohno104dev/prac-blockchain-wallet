package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	b := &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
	}

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp\t%d\n", b.timestamp)
	fmt.Printf("nonce\t\t%d\n", b.nonce)
	fmt.Printf("previous_hash\t%x\n", b.previousHash)
	fmt.Printf("transactions\t%s\n", b.transactions)
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int      `json:"nonce"`
		previousHash [32]byte `json:"previous_hash"`
		Timestamp    int64    `json:"timestamp"`
		Transactions []string `json:"transactions"`
	}{
		Nonce:        b.nonce,
		previousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 15), i, strings.Repeat("=", 15))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("#", 39))
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func main() {
	bc := NewBlockchain()
	bc.Print()

	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(23, previousHash)
	bc.Print()

	previousHash = bc.LastBlock().Hash()
	bc.CreateBlock(42, previousHash)
	bc.Print()
}
