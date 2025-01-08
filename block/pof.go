package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24 //

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// SHA-256, left shift 'targetBits'
	target = target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{Block: b, Target: target}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.HashTransactions(),
		IntToHex(pow.Block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})

	return data
}

// do mining
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int // hashInt is the integer representation of hash
	var hash [32]byte
	nonce := 0
	maxNonce := math.MaxInt64

	fmt.Printf("Mining the block targeting %x\n", pow.Target)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 { // less than target
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data) // unique
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.Target) == -1
}
