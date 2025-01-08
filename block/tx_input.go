package block

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	TxID      []byte // id of transaction that input belongs to
	Vout      int    // previous TX output's index
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
