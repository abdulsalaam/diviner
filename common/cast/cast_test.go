package cast

import (
	"encoding/hex"
	"math"
	"testing"
)

func TestCast(t *testing.T) {

	pi := math.Pi

	bytes, err := ToBytes(pi)

	if err != nil {
		t.Fatal(err)
	}

	pi2, err := ToFloat64(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if pi2 != pi {
		t.Fatalf("must to equals: %v, %v\n", pi2, pi)
	}

	hello := []byte("hello world")

	hello2, err := ToBytes(hello)
	if err != nil {
		t.Fatal(err)
	}

	if hex.EncodeToString(hello) != hex.EncodeToString(hello2) {
		t.Fatal("byte array not match")
	}

	b3, err := ToBytes(hello, pi)
	if err != nil {
		t.Fatal(err)
	}

	if hex.EncodeToString(b3) != (hex.EncodeToString(hello) + hex.EncodeToString(bytes)) {
		t.Fatal("concate not match")
	}

	str := "hello world"
	strBytes := []byte(str)
	strBytes2, err := ToBytes(str)
	if err != nil {
		t.Fatal(err)
	}

	if hex.EncodeToString(strBytes) != hex.EncodeToString(strBytes2) {
		t.Fatal("string not match")
	}
}
