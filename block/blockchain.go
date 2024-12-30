package block

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]   // get the latest block
	newBlock := NewBlock(data, prevBlock.Hash) // create a new block
	bc.Blocks = append(bc.Blocks, newBlock)    // add it to the chain
}

// genesis block
func GenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{GenesisBlock()}}
}
