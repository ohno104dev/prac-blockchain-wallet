package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ohno104dev/prac-blockchain-wallet-go/blockchain"
	"github.com/ohno104dev/prac-blockchain-wallet-go/utils"
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
		log.Println("========bc:", bc)
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("Error: Invalid HTTP method")
	}
}

func (bcn *BlockchainNode) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*blockchain.Transaction `json:"transactions"`
			Length       int                       `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})

		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("ERROR missing field(S)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.String2PublicKey(*t.SenderPublicKey)
		signature := utils.String2Signature(*t.Signature)
		bc := bcn.GetBlockchain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))

	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("ERROR missing field(S)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.String2PublicKey(*t.SenderPublicKey)
		signature := utils.String2Signature(*t.Signature)
		bc := bcn.GetBlockchain()
		isUpdated := bc.AddTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))

	case http.MethodDelete:
		bc := bcn.GetBlockchain()
		bc.ClearTransactionPool()
		io.WriteString(w, string(utils.JsonStatus("success")))

	default:
		log.Println("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)

	}
}

func (bcn *BlockchainNode) Mine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bcn.GetBlockchain()
		isMine := bc.Mining()

		var m []byte
		if !isMine {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) StartMining(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bcn.GetBlockchain()
		bc.StartMining()

		m := utils.JsonStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) Amount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		blockchainAddress := r.URL.Query().Get("blockchain_address")
		amount := bcn.GetBlockchain().CalculateTotalAmount(blockchainAddress)

		ar := &blockchain.AmountResponse{Amount: amount}
		m, _ := ar.MarshalJSON()

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) Consensus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		bc := bcn.GetBlockchain()
		replaced := bc.ResolveConflicts()

		w.Header().Add("Content-Type", "application/json")
		if replaced {
			io.WriteString(w, string(utils.JsonStatus("success")))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:
		log.Printf("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) Run() {
	bcn.GetBlockchain().Run()

	http.HandleFunc("/", bcn.GetChain)
	http.HandleFunc("/transactions", bcn.Transactions)
	http.HandleFunc("/mine", bcn.Mine)
	http.HandleFunc("/mine/start", bcn.StartMining)
	http.HandleFunc("/amount", bcn.Amount)
	http.HandleFunc("/consensus", bcn.Consensus)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.port)), nil))
}
