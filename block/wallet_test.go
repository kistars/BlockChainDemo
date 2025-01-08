package block

import (
	"fmt"
	"testing"
)

func TestGenerateAddr(t *testing.T) {
	wallet := NewWallet()
	addrBytes := wallet.GetAddress()
	fmt.Println(string(addrBytes))
}
