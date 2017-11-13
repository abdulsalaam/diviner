package base58

import (
	"diviner/common/csp"
	"testing"
)

var testdata = [][]string{
	{"hello", "Cn8eVZg"},
	{"1234567890", "3mJr7AoUCHxNqd"},
	{"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "QNoRV1sxwosKt47hNWXhBXyUVZxxn47YeXd9H4Pp8Td6CkwUX7BtmKBHwFXgVGtgGRUS1so"},
}

func TestEncode(t *testing.T) {
	for _, x := range testdata {
		bytes := []byte(x[0])

		if Encode(bytes) != x[1] {
			t.Fatalf("encoding error: %q, %q, %q", x[0], x[1], Encode(bytes))
		}
	}
}

func TestDecode(t *testing.T) {
	for _, x := range testdata {
		bytes := Decode(x[1])
		if x[0] != string(bytes) {
			t.Fatalf("decoding error: %q", x[0])
		}
	}
}

func TestPubKeyToAddress(t *testing.T) {
	priv, _ := csp.GeneratePrivateTempKey()
	pub, _ := priv.PublicKey()

	_, err := Address(pub)
	if err != nil {
		t.Fatal(err)
	}
}
