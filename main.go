package main

import (
	"blockchain/block"
	"fmt"
	"strconv"
)

func main() {
	bc := block.NewBlockchain()
	defer bc.CloseDB()

	//bc.AddBlock("Send 3 BTC to Alice")
	//bc.AddBlock("Send 5 BTCs to Bob")

	iter := bc.Iterator()

	for {
		b := iter.Next()
		fmt.Printf("Prev.hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)

		// validate
		pof := block.NewProofOfWork(b)
		fmt.Printf("validate: %s\n", strconv.FormatBool(pof.Validate()))
		fmt.Println()

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

}
