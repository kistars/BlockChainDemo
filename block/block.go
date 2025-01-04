package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int // 'number once', arbitrary number that's only used once
}

// calculate the hash of a block
//func (b *Block) SetHash() {
//	ts := []byte(strconv.FormatInt(b.Timestamp, 10))
//	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, ts}, []byte{})
//	hash := sha256.Sum256(headers)
//
//	b.Hash = hash[:]
//}

// genesis block
func GenesisBlock(coinbase []*Transaction) *Block {
	return NewBlock(coinbase, []byte{})
}

// Create a new block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}

	pof := NewProofOfWork(block)
	nonce, hash := pof.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// serialization
func (b *Block) Serialization() []byte {
	var buf bytes.Buffer // store the serialized data
	encoder := gob.NewEncoder(&buf)
	_ = encoder.Encode(b)

	return buf.Bytes()
}

func Deserializaion(d []byte) *Block {
	var block Block
	buf := bytes.NewReader(d)
	decoder := gob.NewDecoder(buf)
	_ = decoder.Decode(&block)
	return &block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
