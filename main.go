package main

import (
	"blockchain/block"
)

func main() {
	bc := block.NewBlockchain("Tom")
	defer bc.CloseDB()

	cli := block.NewCLI(bc)
	cli.Run()
}
