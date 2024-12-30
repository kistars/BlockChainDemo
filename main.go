package main

import (
	"blockchain/block"
	"fmt"
)

func main() {
	bc := block.NewBlockchain()

	bc.AddBlock("Send 1 BTC to Alice")
	bc.AddBlock("Send 2 BTCs to Bob")

	for _, b := range bc.Blocks {
		fmt.Printf("Prev.hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Println()
	}
}
