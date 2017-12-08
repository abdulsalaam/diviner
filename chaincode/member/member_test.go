package member

import (
	"os"
	"testing"

	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"diviner/common/cast"
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

func TestRegister(t *testing.T) {
	b1, _ := pbm.Marshal(m1)
	v1, _ := pbc.NewVerification(priv1, b1)
	bv1, _ := pbc.Marshal(v1)

	resp := ccc.MockInvoke(stub, register, b1, bv1)
	if ccc.NotOK(&resp) {
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

	if tmp.Subsidy != 0.0 {
		t.Fatalf("subsidy must be 0.0, but %v", tmp.Subsidy)
	}

	resp = ccc.MockInvoke(stub, register, b1, bv1)
	if ccc.OK(&resp) {
		t.Fatalf("can not register an existed member")
	}

	m1 = tmp
}

func TestQueryAndFind(t *testing.T) {

	addr := []byte(m1.Address)
	v, _ := pbc.NewVerification(priv1, addr)
	bv, _ := pbc.Marshal(v)

	resp := ccc.MockInvoke(stub, query, addr, bv)
	if ccc.NotOK(&resp) {
		t.Fatalf("query failed: %d, %s", resp.Status, resp.Message)
	}

	tmp, _ := pbm.Unmarshal(resp.Payload)
	if tmp.String() != m1.String() {
		t.Fatalf("data not match: %s, %s", tmp.String(), m1.String())
	}

	resp = ccc.MockInvoke(stub, []byte("find"), addr)
	if ccc.NotOK(&resp) {
		t.Fatalf("find failed: %d, %s", resp.Status, resp.Message)
	}

	tmp, _ = pbm.Unmarshal(resp.Payload)
	if tmp.String() != m1.String() {
		t.Fatalf("data not match: %s, %s", tmp.String(), m1.String())
	}

	addr = []byte("abc")
	v, _ = pbc.NewVerification(priv1, addr)
	bv, _ = pbc.Marshal(v)
	resp = ccc.MockInvoke(stub, query, addr, bv)
	if resp.Status != ccc.FORBIDDEN {
		t.Fatalf("can not query others")
	}

	resp = ccc.MockInvoke(stub, []byte("find"), addr)
	if resp.Status != ccc.NOTFOUND {
		t.Fatalf("find non-existed member code failed: %d, %s", resp.Status, resp.Message)
	}

	addr = []byte(m2.Address)
	v, _ = pbc.NewVerification(priv2, addr)
	bv, _ = pbc.Marshal(v)

	resp = ccc.MockInvoke(stub, query, addr, bv)
	if resp.Status != ccc.NOTFOUND {
		t.Fatalf("can not query non-existed member")
	}
}

func TestTransfer(t *testing.T) {
	b2, _ := pbm.Marshal(m2)
	v2, _ := pbc.NewVerification(priv1, b2)
	bv2, _ := pbc.Marshal(v2)

	resp := ccc.MockInvoke(stub, register, b2, bv2)
	m2, _ = pbm.Unmarshal(resp.Payload)

	if m1.Balance != m2.Balance || m1.Balance != balance {
		t.Fatalf("initial balance failed, all must be %v, but %v, %v", balance, m1.Balance, m2.Balance)
	}

	amount := 100.0
	tx, _ := pbm.NewTransfer(m2.Address, amount)

	txbytes, _ := pbm.MarshalTransfer(tx)
	vtx, _ := pbc.NewVerification(priv1, txbytes)
	vtxbytes, _ := pbc.Marshal(vtx)

	resp = ccc.MockInvoke(stub, transfer, txbytes, vtxbytes)
	if ccc.NotOK(&resp) {
		t.Fatalf("transfer failed: %s", resp.Message)
	}

	m1, _ = pbm.Unmarshal(resp.Payload)
	minus := amount * (1.0 + fee)
	if m1.Balance != (balance - minus) {
		t.Fatalf("source balance failed, must be %v, but %v", balance-minus, m1.Balance)
	}

	addr := []byte(m2.Address)
	v2, _ = pbc.NewVerification(priv2, addr)
	bv2, _ = pbc.Marshal(v2)

	resp = ccc.MockInvoke(stub, query, addr, bv2)
	if ccc.NotOK(&resp) {
		t.Fatalf("can not find target member")
	}

	m2, _ = pbm.Unmarshal(resp.Payload)
	if m2.Balance != balance+amount {
		t.Fatalf("target balance failed, must be %v, but %v", balance+amount, m2.Balance)
	}
}

func TestBalanceAndSubsidy(t *testing.T) {
	m1, _, _ = ccu.GetMemberAndCheck(stub, m1.Address)

	old := m1.Balance
	amount := 1000.0
	amountBytes, _ := cast.ToBytes(amount)

	resp := ccc.MockInvoke(stub, []byte("balance"), []byte(m1.Address), amountBytes)
	if ccc.NotOK(&resp) {
		t.Fatalf("update balance failed: %d, %s", resp.Status, resp.Message)
	}

	m1, _, _ = ccu.GetMemberAndCheck(stub, m1.Address)
	if m1.Balance != (old + amount) {
		t.Fatalf("member balace failed, must be %v, but %v", old+amount, m1.Balance)
	}

	old = m1.Balance
	oldSubsidy := m1.Subsidy
	resp = ccc.MockInvoke(stub, []byte("subsidy"), []byte(m1.Address), amountBytes)
	if ccc.NotOK(&resp) {
		t.Fatalf("subsidy failed: %d, %s", resp.Status, resp.Message)
	}

	m1, _, _ = ccu.GetMemberAndCheck(stub, m1.Address)
	if m1.Balance != (old + amount) {
		t.Fatalf("member balace failed, must be %v, but %v", old+amount, m1.Balance)
	}
	if m1.Subsidy != (oldSubsidy + amount) {
		t.Fatalf("member subsidy failed, must be %v, but %v", oldSubsidy+amount, m1.Subsidy)
	}
}

func TestMain(m *testing.M) {

	priv1, _ = csp.GeneratePrivateTempKey()
	priv2, _ = csp.GeneratePrivateTempKey()

	m1, _ = pbm.NewMemberWithPrivateKey(priv1, 0.0)
	m2, _ = pbm.NewMemberWithPrivateKey(priv2, 0.0)

	stub = ccc.NewMockStub("member", NewMemberChaincode())
	ccc.MockInit(stub)

	os.Exit(m.Run())
}
