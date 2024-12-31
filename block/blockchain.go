package block

import (
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const latestKey = "latest" // specify the latest block's hash

// A blockchain can have multiple branches, and it’s the longest of them that’s considered main
type Blockchain struct {
	tip []byte   // latest block's hash
	db  *bolt.DB // store the blocks
}

func (bc *Blockchain) CloseDB() {
	_ = bc.db.Close()
}

func (bc *Blockchain) AddBlock(data string) {
	// get previous block's hash
	var prevHash []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		prevHash = b.Get([]byte("latest"))
		return nil
	})

	// create the next block
	newBlock := NewBlock(data, prevHash)
	bc.tip = newBlock.Hash

	// insert into db
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		_ = b.Put(newBlock.Hash, newBlock.Serialization())
		_ = b.Put([]byte(latestKey), newBlock.Hash)

		return nil
	})
}

// genesis block
func GenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	//return &Blockchain{Blocks: []*Block{GenesisBlock()}}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := GenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			_ = b.Put(genesis.Hash, genesis.Serialization())
			_ = b.Put([]byte(latestKey), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte(latestKey))
		}
		return nil
	})

	bc := Blockchain{tip: tip, db: db}

	return &bc
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bi := &BlockchainIterator{currentHash: bc.tip, db: bc.db}
	return bi
}
