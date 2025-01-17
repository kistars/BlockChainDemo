package block

import (
	"bytes"
	"fmt"
	"testing"
)

func TestIntToHex(t *testing.T) {
	a := IntToHex(16)
	fmt.Printf("%x\n", a)
}

func TestBytesJoin(t *testing.T) {
	var a []byte
	var b []byte
	a = []byte("a")
	b = []byte("b")
	bytesArr := [][]byte{a, b}

	res := bytes.Join(bytesArr, []byte(""))

	fmt.Println(string(res))
}
