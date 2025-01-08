package block

import (
	"encoding/hex"
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

const utxoBucket = "utxoBucket"

// cache of blocks
type UTXOSet struct {
	Blockchain *Blockchain
}

func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err1 := tx.DeleteBucket(bucketName)
		if err1 != nil && !errors.Is(err1, bolt.ErrBucketNotFound) {
			log.Fatal(err1)
		}
		_, err1 = tx.CreateBucket(bucketName)
		if err1 != nil {
			log.Fatal(err1)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for id, outs := range UTXO {
			key, _ := hex.DecodeString(id)
			err2 := b.Put(key, outs.Serialize())
			if err2 != nil {
				log.Fatal(err2)
			}
		}

		return nil
	})
}

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	upspendableOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		// iterate the utxo bucket
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outputs := DeserializeOutputs(v)

			for idx, out := range outputs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					upspendableOutputs[txID] = append(upspendableOutputs[txID], idx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return accumulated, upspendableOutputs
}

func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	UTXOs := make([]TXOutput, 0)
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outputs := DeserializeOutputs(v)

			for _, output := range outputs.Outputs {
				if output.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, output)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)

	}

	return UTXOs
}

func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, transaction := range block.Transactions {
			if !transaction.IsCoinbase() {
				for _, input := range transaction.Vin {
					data := b.Get(input.TxID)
					outputs := DeserializeOutputs(data)
					updateOutputs := TXOutputs{}

					// find unspent outputs
					for outIdx, output := range outputs.Outputs {
						if input.Vout != outIdx {
							updateOutputs.Outputs = append(updateOutputs.Outputs, output)
						}
					}

					if len(updateOutputs.Outputs) == 0 {
						err1 := b.Delete(input.TxID)
						if err1 != nil {
							log.Fatal(err1)
						}
					} else {
						// update UTXO set
						err2 := b.Put(input.TxID, updateOutputs.Serialize())
						if err2 != nil {
							log.Fatal(err2)
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			for _, output := range transaction.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, output)
			}
			err3 := b.Put(transaction.ID, newOutputs.Serialize())
			if err3 != nil {
				log.Fatal(err3)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}
