package block

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int // 'number once', arbitrary number that's only used once
	Height        int
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Height:        height,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// serialization
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer // store the serialized data
	encoder := gob.NewEncoder(&buf)
	_ = encoder.Encode(b)

	return buf.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	buf := bytes.NewReader(d)
	decoder := gob.NewDecoder(buf)
	_ = decoder.Decode(&block)
	return &block
}

func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}

	// The root of the tree will serve as the unique identifier of block’s transactions
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}
