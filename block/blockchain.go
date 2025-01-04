package block

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const latestKey = "latest" // specify the latest block's hash
const genesisCoinbaseData = "genesis_coinbase"

// A blockchain can have multiple branches, and it’s the longest of them that’s considered main
type Blockchain struct {
	tip []byte   // latest block's hash
	db  *bolt.DB // store the blocks
}

func (bc *Blockchain) CloseDB() {
	_ = bc.db.Close()
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	// get previous block's hash
	var prevHash []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		prevHash = b.Get([]byte("latest"))
		return nil
	})

	// create the next block
	newBlock := NewBlock(transactions, prevHash)
	bc.tip = newBlock.Hash

	// insert into db
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		_ = b.Put(newBlock.Hash, newBlock.Serialization())
		_ = b.Put([]byte(latestKey), newBlock.Hash)

		return nil
	})
}

func NewBlockchain(address string) *Blockchain {
	//return &Blockchain{Blocks: []*Block{GenesisBlock()}}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			coinbaseTX := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := GenesisBlock([]*Transaction{coinbaseTX})

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

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOutputs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the out spent?
				if spentTXOutputs[txID] != nil {
					for _, spendOut := range spentTXOutputs[txID] {
						if spendOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTXOutputs[inTxID] = append(spentTXOutputs[inTxID], in.Vout)
					}
				}
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	transactions := bc.FindUnspentTransactions(address)

	for _, tx := range transactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

/*
address: receiver's address
amount: amount of sending coins
The method iterates over all unspent transactions and accumulates their values,
return balance and valid outputs
*/
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, output := range tx.Vout {
			if output.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += output.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {

}
