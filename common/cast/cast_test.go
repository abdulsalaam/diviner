package cast

import (
	"diviner/common/util"
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

	pi2, err := BytesToFloat64(bytes)
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

	x1, err := ToBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	x2, err := BytesToBool(x1)
	if err != nil {
		t.Fatal(err)
	}

	if x2 != false {
		t.Fatalf("data not match: %v, %v", x2, false)
	}

	x1, err = ToBytes(true)
	if err != nil {
		t.Fatal(err)
	}

	x2, err = BytesToBool(x1)
	if err != nil {
		t.Fatal(err)
	}

	if x2 != true {
		t.Fatalf("data not match: %v, %v", x2, true)
	}

	/*_, err = ToBytes(123)
	if err != nil {
		t.Fatal(err)
	}*/

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

func TestStringsToByteArray(t *testing.T) {
	a := "A"
	b := "B"
	c := "C"

	ab := []byte(a)
	bb := []byte(b)
	cb := []byte(c)

	result := StringsToByteArray(a, b, c)
	if len(result) != 3 {
		t.Fatal("length failed")
	}

	if util.CmpByteArray(result[0], ab) != 0 ||
		util.CmpByteArray(result[1], bb) != 0 ||
		util.CmpByteArray(result[2], cb) != 0 {
		t.Fatal("content failed")
	}
}
