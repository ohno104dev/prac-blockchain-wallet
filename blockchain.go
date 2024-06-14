package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "BLOCKCHAIN REWARD SYSTEM (e.g. minting & fees)"
	MINING_REWARD     = 1.0
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp\t%d\n", b.timestamp)
	fmt.Printf("nonce\t\t%d\n", b.nonce)
	fmt.Printf("previous_hash\t%x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	// fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		previousHash [32]byte       `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		previousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	BlockchainAddress string
}

func NewBlockchain(addr string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.BlockchainAddress = addr
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
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

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.BlockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	fmt.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) bool {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
	return false
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{nonce, previousHash, 0, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	// fmt.Println(guessHashStr)
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0

	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("  sender_blockchain_address\t%s\n", t.senderBlockchainAddress)
	fmt.Printf("  recipient_blockchain_address\t%s\n", t.recipientBlockchainAddress)
	fmt.Printf("  value\t\t\t\t%1f\n", t.value)
}

func main() {
	minerBlockchainAddr := "miners_blockchain_addr"
	bc := NewBlockchain(minerBlockchainAddr)
	bc.Print()

	bc.AddTransaction("Rock", "Morty", 137)
	bc.Mining()
	bc.Print()

	bc.AddTransaction("Stan", "Francine", 23)
	bc.AddTransaction("Fry", "Ben", 42)
	bc.Mining()
	bc.Print()
}
