package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ohno104dev/prac-blockchain-wallet-go/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "BLOCKCHAIN REWARD SYSTEM (e.g. minting & fees)"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	BlockchainAddress string
	port              uint16
	mux               sync.Mutex
}

func NewBlockchain(addr string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.BlockchainAddress = addr
	bc.CreateBlock(0, b.Hash())
	bc.port = port

	return bc
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
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
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.TransactionPool()) == 0 {
		return false
	}

	bc.AddTransaction(MINING_SENDER, bc.BlockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	fmt.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)
	// TODO
	// sync
	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// XXX: For easy test
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: not enough balance in wallet")
		// 	return false
		// }

		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Could not verify transaction")
	}

	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))

	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
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

func (bc *Blockchain) CalculateTotalAmount(blockchainAddr string) float32 {
	var totalAmount float32 = 0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddr == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddr == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}

	return totalAmount
}
