package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"github.com/ohno104dev/prac-blockchain-wallet-go/utils"
)

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		privateKey,
		publicKey,
		sender,
		recipient,
		value,
	}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {

	if tr.SenderPrivateKey == nil ||
		tr.SenderPublicKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Value == nil {
		return false
	}

	return true
}
