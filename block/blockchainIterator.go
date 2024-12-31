package block

import "github.com/boltdb/bolt"

type BlockchainIterator struct {
	currentHash []byte   // current block's hash
	db          *bolt.DB // the whole blockchain
}

func (bi *BlockchainIterator) Next() *Block {
	var block *Block
	_ = bi.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		data := b.Get(bi.currentHash)
		block = Deserializaion(data)
		return nil
	})

	bi.currentHash = block.PrevBlockHash // reverse order (from new to old blocks)

	return block
}
