package block

import "testing"

func TestIntToHex(t *testing.T) {
	a := IntToHex(1111)
	t.Log(string(a))
}
