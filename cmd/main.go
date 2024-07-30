package main

import (
	"fmt"

	"github.com/ohno104dev/prac-blockchain-wallet-go/utils"
)

func main() {
	myAddress := utils.GetHost()
	fmt.Println(utils.FindNeighbors(myAddress, 3333, 0, 3, 3333, 3336))
}
