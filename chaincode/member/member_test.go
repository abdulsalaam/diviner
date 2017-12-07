package main

import (
	"os"
	"testing"

	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbc "diviner/protos/common"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	priv1 bccsp.Key
	priv2 bccsp.Key
	m1    *pbm.Member
	m2    *pbm.Member
	stub  *shim.MockStub

	fee      = 0.001
	balance  = 10000.0
	register = []byte("register")
	query    = []byte("query")
	transfer = []byte("transfer")
)

func TestMain(m *testing.M) {

	priv1, _ = csp.GeneratePrivateTempKey()
	priv2, _ = csp.GeneratePrivateTempKey()

	m1, _ = pbm.NewMember(priv1, 0.0)
	m2, _ = pbm.NewMember(priv2, 0.0)

	stub = ccc.NewMockStub("member", NewMemberChaincode())
	ccc.MockInit(stub)

	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	b1, _ := pbm.Marshal(m1)
	v1, _ := pbc.NewVerification(priv1, b1)
	bv1, _ := pbc.Marshal(v1)

	resp := ccc.MockInvoke(stub, register, b1, bv1)
	if !ccc.OK(&resp) {
		t.Fatalf("register failed: %s", resp.Message)
	}

	tmp, _ := pbm.Unmarshal(resp.Payload)

	if tmp.Address != m1.Address {
		t.Fatalf("address not match: %s, %s", tmp.Address, m1.Address)
	}

	if tmp.Balance != balance {
		t.Fatalf("balance failed, must be %v, but %v", balance, tmp.Balance)
	}

	if tmp.Blocked {
		t.Fatalf("block status failed, must be %v, but %v", false, tmp.Blocked)
	}

	if len(tmp.Assets) != 0 {
		t.Fatalf("size of asset must be %d, but %d", 0, len(tmp.Assets))
	}
}
