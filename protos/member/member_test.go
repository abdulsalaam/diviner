package member

import (
	"diviner/common/base58"
	"diviner/common/csp"
	"testing"
)

func TestMember(t *testing.T) {
	priv, _ := csp.GeneratePrivateTempKey()

	m, err := NewMemberWithPrivateKey(priv, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	pub, _ := priv.PublicKey()
	addr, _ := base58.Address(pub)

	if m.Address != addr {
		t.Fatalf("address failed, must be %s, but %s", addr, m.Address)
	}

	if m.Blocked {
		t.Fatalf("block status failed, must be false, but true")
	}

	if m.Balance != 100.0 {
		t.Fatalf("balance failed, must be %v, but %v", 100.0, m.Balance)
	}

	if m.Subsidy != 0.0 {
		t.Fatalf("subsidy failed, must be 0.0, but %v", m.Subsidy)
	}

	if len(m.Assets) != 0 {
		t.Fatalf("asset size failed, must be 0, but %v", len(m.Assets))
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
