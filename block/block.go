package block

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int // 'number once', arbitrary number that's only used once
}

// calculate the hash of a block
func (b *Block) SetHash() {
	ts := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, ts}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// Create a new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}

	pof := NewProofOfWork(block)
	nonce, hash := pof.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}
