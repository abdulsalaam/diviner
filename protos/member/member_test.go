package member

import (
	"diviner/common/csp"
	"testing"
)

func TestMember(t *testing.T) {
	priv, _ := csp.GeneratePrivateTempKey()

	m, err := NewMember(priv, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	b, err := Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}

	if m.Address != m2.Address || m.Balance != m2.Balance {
		t.Fatal("decode failed")
	}

}
