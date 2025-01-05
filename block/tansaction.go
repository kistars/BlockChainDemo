package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 10 // rewards

type Transaction struct {
	ID       []byte
	VInputs  []TXInput
	VOutputs []TXOutput
}

type TXInput struct {
	TxID      []byte // id of transaction that input belongs to
	OutputIdx int    // previous TX output's index
	Signature []byte
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingKey := HashPubKey(in.PubKey)
	return bytes.Compare(lockingKey, pubKeyHash) == 0
}

// Output stores the coins
type TXOutput struct {
	Value      int    // coin count
	PubKeyHash []byte // public key hash
}

func NewTXOutput(value int, address string) *TXOutput {
	output := &TXOutput{value, nil}
	output.Lock([]byte(address))
	return output
}

// address: receiver's address
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen] // 1byte version + pub key hash + 4bytes checksum
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.VInputs) == 1 && len(tx.VInputs[0].TxID) == 0 && tx.VInputs[0].OutputIdx == -1
}

func (tx *Transaction) Serialize() []byte {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encode.Bytes()
}

// return the hash of the transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()

	return &tx
}

//func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
//	return in.ScriptSig == unlockingData
//}
//
//func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
//	return out.ScriptPubKey == unlockingData
//}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := TXInput{
				TxID:      txID,
				OutputIdx: out,
				//ScriptSig: from,
			}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	output := NewTXOutput(amount, to) // locked with receiver's address
	outputs = append(outputs, *output)

	if acc > amount {
		// locked with sender's address
		out := NewTXOutput(acc-amount, from)
		outputs = append(outputs, *out) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	//bc.SignTransaction(tx,

	return &tx
}

func (tx *Transaction) Sign(priKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	// we need to access the outputs referenced in the inputs of the transaction
	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.VInputs {
		prevTX := prevTXs[hex.EncodeToString(vin.TxID)]
		txCopy.VInputs[inID].Signature = nil
		txCopy.VInputs[inID].PubKey = prevTX.VOutputs[vin.OutputIdx].PubKeyHash
		txCopy.Hash()
		txCopy.VInputs[inID].PubKey = nil

		r, s, _ := ecdsa.Sign(rand.Reader, &priKey, txCopy.ID)
		signature := append(r.Bytes(), s.Bytes()...)
		txCopy.VInputs[inID].Signature = signature
	}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.VInputs {
		prevTX := prevTXs[hex.EncodeToString(vin.TxID)]
		txCopy.VInputs[inID].Signature = nil
		txCopy.VInputs[inID].PubKey = prevTX.VOutputs[vin.OutputIdx].PubKeyHash
		txCopy.Hash()
		txCopy.VInputs[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)]) // first half
		s.SetBytes(vin.Signature[(sigLen / 2):]) // second half

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.Signature[:(keyLen / 2)]) // first half
		y.SetBytes(vin.Signature[(keyLen / 2):]) // second half

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
}

// return parts of a transaction that need to be signed
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.VInputs {
		input := TXInput{vin.TxID, vin.OutputIdx, nil, nil}
		inputs = append(inputs, input)
	}

	for _, vout := range tx.VOutputs {
		output := TXOutput{vout.Value, vout.PubKeyHash}
		outputs = append(outputs, output)
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}
