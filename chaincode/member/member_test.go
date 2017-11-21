package main

import (
	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbm "diviner/protos/member"
	"testing"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	query  = []byte("query")
	create = []byte("create")
	update = []byte("update")
	priv1  bccsp.Key
	priv2  bccsp.Key
	priv3  bccsp.Key
	m1     *pbm.Member
	m2     *pbm.Member
	m3     *pbm.Member
	stub   *shim.MockStub
)

func TestMain(m *testing.M) {
	priv1, _ = csp.GeneratePrivateTempKey()
	priv2, _ = csp.GeneratePrivateTempKey()
	priv3, _ = csp.GeneratePrivateTempKey()
	m1, _ = pbm.NewMember(priv1, 0.0)
	m2, _ = pbm.NewMember(priv2, 10.0)
	m3, _ = pbm.NewMember(priv3, 1000.0)
	stub = ccc.NewMockStub("member", NewMemberChaincode())

	ccc.MockInit(stub)

	m.Run()
}

func TestQueryMember(t *testing.T) {
	resp := ccc.MockInvoke(stub, query, []byte("abc"))
	if ccc.OK(&resp) {
		t.Fatal("can not query non-existed id")
	}
}
func TestCreateMember(t *testing.T) {
	b1, _ := pbm.Marshal(m1)
	b2, _ := pbm.Marshal(m2)

	resp := ccc.MockInvoke(stub, create, b1)
	if !ccc.OK(&resp) {
		t.Fatalf("create membe failed: %v\n", resp.Message)
	}

	resp = ccc.MockInvoke(stub, create, b2)
	if !ccc.OK(&resp) {
		t.Fatalf("create membe failed: %v\n", resp.Message)
	}

	resp = ccc.MockInvoke(stub, query, []byte(m1.Id))
	if resp.GetStatus() != shim.OK {
		t.Fatalf("query member failed: %v", resp.Message)
	}

	mm1, err := pbm.Unmarshal(resp.Payload)
	if err != nil {
		t.Fatal("data structure not match")
	}

	if mm1.Address != m1.Address || mm1.Balance != m1.Balance || mm1.Id != m1.Id {
		t.Fatal("data not match")
	}

	resp = ccc.MockInvoke(stub, create, b1)
	if ccc.OK(&resp) {
		t.Fatal("can not create an existed member")
	}
}

func TestUpdateMember(t *testing.T) {
	b3, _ := pbm.Marshal(m3)
	resp := ccc.MockInvoke(stub, update, b3)
	if ccc.OK(&resp) {
		t.Fatal("can not update an non-existed member")
	}

	m1.Balance = 1000.0
	b1, _ := pbm.Marshal(m1)
	resp = ccc.MockInvoke(stub, update, b1)
	if !ccc.OK(&resp) {
		t.Fatalf("update error: %v", resp.Message)
	}

	resp = ccc.MockInvoke(stub, query, []byte(m1.Id))
	if !ccc.OK(&resp) {
		t.Fatalf("query error: %v", resp.Message)
	}

	mm1, _ := pbm.Unmarshal(resp.Payload)
	if mm1.Id != m1.Id || mm1.Address != m1.Address || mm1.Balance != m1.Balance {
		t.Fatal("data not match")
	}
}
