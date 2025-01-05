package block

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const addressChecksumLen = 4
const version = 0x00

// A wallet is nothing but a key pair
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() *Wallet {
	priKey, pubKey := newKeyPair()
	wallet := &Wallet{priKey, pubKey}
	return wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	priKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(priKey.X.Bytes(), priKey.Y.Bytes()...)

	return *priKey, pubKey
}

func (w Wallet) GetAddress() []byte {
	pushKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{version}, pushKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte {
	pubKeySha256 := sha256.Sum256(pubKey)

	RIPEMD160Hash := ripemd160.New()
	_, _ = RIPEMD160Hash.Write(pubKeySha256[:])
	pubRIPEMD160 := RIPEMD160Hash.Sum(nil)

	return pubRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
