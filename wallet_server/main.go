package main

import (
	"flag"
	"log"
)

func Init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP Number for Online Wallet")
	gateway := flag.String("gateway", "http://127.0.0.1:3333", "Blockchain Gateway")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	log.Println("Starting wallet server on port:", *port, "using blockchain node", *gateway, "as gateway")
	app.Run()
}

type A struct {
}

func (A) Hello() {}

type IA interface {
	Hello()
}

var _ IA = (*A)(nil)
