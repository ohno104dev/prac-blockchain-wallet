package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type BlockchainNode struct {
	port uint16
}

func NewBlockchainNode(port uint16) *BlockchainNode {
	return &BlockchainNode{port}
}

func (bcn *BlockchainNode) Port() uint16 {
	return bcn.port
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func (bcn *BlockchainNode) Run() {
	http.HandleFunc("/", HelloWorld)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.port)), nil))
}
