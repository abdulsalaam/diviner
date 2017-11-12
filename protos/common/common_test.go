package common

import (
	"diviner/common/csp"
	"testing"
)

func TestVerification(t *testing.T) {
	var expired int64 = 5 * 60

	priv, _ := csp.GeneratePrivateTempKey()

	in := []byte("test")

	v, err := NewVerification(priv, in)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := Verify(v, in, expired)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("verify failed")
	}
}

func TestWrongData(t *testing.T) {
	var expired int64 = 5 * 60
	priv, _ := csp.GeneratePrivateTempKey()
	pub, _ := priv.PublicKey()

	in := []byte("test")
	hash, _ := csp.Hash([]byte("test1"))

	_, err := NewVerification(pub, in)
	if err == nil {
		t.Fatal("can not use public key")
	}

	v, _ := NewVerification(priv, in)
	old := v.Hash

	v.Hash = hash
	ok, err := Verify(v, in, expired)
	if ok && err == nil {
		t.Fatal("verify failed")
	}

	v.Hash = old
	v.Time.Seconds -= 6 * 60

	ok, err = Verify(v, in, expired)
	if ok && err == nil {
		t.Fatal("verify failed")
	}

}
