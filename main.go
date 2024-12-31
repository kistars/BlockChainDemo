package main

import (
	"blockchain/block"
	"fmt"
	"strconv"
)

func main() {
	bc := block.NewBlockchain()

	bc.AddBlock("Send 1 BTC to Alice")
	bc.AddBlock("Send 2 BTCs to Bob")

	for _, b := range bc.Blocks {
		fmt.Printf("Prev.hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)

		// validate proof of work
		pof := block.NewProofOfWork(b)
		fmt.Printf("Validate: %s\n", strconv.FormatBool(pof.Validate()))

		fmt.Println()
	}
}
