package csp

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger/fabric/bccsp"
)

var (
	priv1 bccsp.Key
	priv2 bccsp.Key
	pub1  bccsp.Key
	pub2  bccsp.Key
	hello = []byte("hello world")
	sign1 []byte
	sign2 []byte
)

func TestMain(m *testing.M) {
	var err error
	prov, err = NewCSP(os.TempDir())
	if err != nil {
		fmt.Printf("init error: %v", err)
		os.Exit(-1)
	}

	priv1, err = GeneratePrivateKey()
	if err != nil {
		fmt.Printf("generate private key error: %v", err)
		os.Exit(-1)
	}

	if priv1.Symmetric() {
		fmt.Printf("must be asymmetic key")
		os.Exit(-1)
	}

	priv2, err = GeneratePrivateKey()
	if err != nil {
		fmt.Printf("generate private key error: %v", err)
		os.Exit(-1)
	}

	pub1, err = priv1.PublicKey()
	if err != nil {
		fmt.Printf("get public key error: %v", err)
		os.Exit(-1)
	}

	pub2, err = priv2.PublicKey()
	if err != nil {
		fmt.Printf("get public key error: %v", err)
		os.Exit(-1)
	}

	sign1, err = Sign(priv1, hello)
	if err != nil {
		fmt.Printf("sign error: %v", err)
		os.Exit(-1)
	}

	sign2, err = Sign(priv2, hello)
	if err != nil {
		fmt.Printf("sign error: %v", err)
		os.Exit(-1)
	}

	if ret := m.Run(); ret != 0 {
		fmt.Printf("Failed testing")
		os.Exit(-1)
	}

	os.Exit(0)
}

func TestGetPubKeyBytes(t *testing.T) {
	bytes1, err := GetPublicKeyBytes(priv1)
	if err != nil {
		t.Fatal(err)
	}

	bytes2, err := pub1.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	if hex.EncodeToString(bytes1) != hex.EncodeToString(bytes2) {
		t.Fatal("pub keys not match")
	}
}

func TestHash(t *testing.T) {
	testdata := [][]string{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"1234567890", "c775e7b757ede630cd0aa1113bd102661ab38829ca52a6422ab782862f268646"},
		{"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWSYZ", "24593a088c32f9744051bc3503f4092eaee97f5cd5809f6e9c70ad4a14ff11f2"},
	}

	for _, x := range testdata {
		hash, err := Hash([]byte(x[0]))
		if err != nil {
			t.Fatal(err)
		}

		if x[1] != hex.EncodeToString(hash) {
			t.Fatalf("hash fail: %q, %q, %q", x[0], x[1], hex.EncodeToString(hash))
		}
	}
}

func TestSign(t *testing.T) {

	if hex.EncodeToString(sign1) == hex.EncodeToString(sign2) {
		t.Fatalf("two signature must not be equal")
	}

	v1, err := Verify(pub1, sign1, hello)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}

	if !v1 {
		t.Fatalf("verify not work")
	}

	v2, err := Verify(pub2, sign2, hello)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}

	if !v2 {
		t.Fatalf("verify not work")
	}

	v3, err := Verify(pub2, sign1, hello)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}

	if v3 {
		t.Fatalf("verify must be wrong with wrong pub key")
	}

	v4, err := Verify(pub1, sign2, hello)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}

	if v4 {
		t.Fatalf("verify must be wrong with wrong pub key")
	}
}

func TestImportPriv(t *testing.T) {
	ski := hex.EncodeToString(priv1.SKI())

	priv3, err := ImportPrivFromSKI(ski)

	if err != nil {
		t.Fatal(err)
	}

	pub3, err := priv3.PublicKey()

	if err != nil {
		t.Fatal(err)
	}

	v1, err := Verify(pub3, sign1, hello)
	if err != nil {
		t.Fatal(err)
	}

	if !v1 {
		t.Fatal("verify error")
	}

	sign3, err := Sign(priv3, hello)
	if err != nil {
		t.Fatal(err)
	}

	v2, err := Verify(pub1, sign3, hello)
	if err != nil {
		t.Fatal(err)
	}

	if !v2 {
		t.Fatal("verify error")
	}
}

func TestImportPub(t *testing.T) {
	bytes, err := pub1.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	pub3, err := ImportPubFromRaw(bytes)
	if err != nil {
		t.Fatal(err)
	}

	v, err := Verify(pub3, sign1, hello)
	if err != nil {
		t.Fatal(err)
	}

	if !v {
		t.Fatal("verify error")
	}
}

/*func TestAddress(t *testing.T) {
	addr1, err := Base58Address(pub1)
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := Base58Address(pub2)
	if err != nil {
		t.Fatal(err)
	}

	if addr1 == addr2 {
		t.Fatal("address must not be equal")
	}
}*/
