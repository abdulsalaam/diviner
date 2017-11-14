package market

import (
	"diviner/common/csp"
	"testing"

	ccc "diviner/chaincode/common"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/google/uuid"
)

var (
	query    = "query"
	create   = "create"
	approve  = "approve"
	stub     = ccc.NewMockStub("market", NewMarketChaincode())
	member   *pbm.Member
	title    = "gogogo"
	outcomes = []string{"yes", "no"}
	event    *pbl.Event
)

func TestMain(m *testing.M) {

	ccc.MockInit(stub, nil)

	priv, _ := csp.GeneratePrivateTempKey()
	member, _ = pbm.NewMember(priv, 10000.0)
	b1, _ := pbm.Marshal(member)

	event, _ = pbl.NewEvent(member.Id, title, outcomes[0], outcomes[1])
	e1, _ := pbl.MarshalEvent(event)

	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	stub.PutState(member.Id, b1)
	stub.PutState(event.Id, e1)
	stub.MockTransactionEnd(txid)

	m.Run()
}

func TestCreateWithFund(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, "create", "fund", member.Id, event.Id, "100.0")
	if !ccc.OK(&resp) {
		t.Fatalf("create market failed: %s", resp.Message)
	}

	mkt, _ := pbl.UnmarshalMarket(resp.Payload)

	resp = ccc.MockInvokeWithString(stub, "query", mkt.Id)
	if !ccc.OK(&resp) {
		t.Fatal("query market failed")
	}

	mkt2, _ := pbl.UnmarshalMarket(resp.Payload)

	if !pbl.CmpMarket(mkt, mkt2) {
		t.Fatal("data error")
	}

	mb, existed, err := ccc.GetStateAndCheck(stub, member.Id)
	if err != nil {
		t.Fatal("query user failed")
	}

	if !existed {
		t.Fatal("user not found")
	}

	member2, _ := pbm.Unmarshal(mb)

	if member2.Balance != (member.Balance - mkt.Fund) {
		t.Fatal("member balance failed")
	}

	member.Balance = member.Balance - mkt.Fund
}

func TestCreateWithLiquidity(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, "create", "liquidity", member.Id, event.Id, "100.0")
	if !ccc.OK(&resp) {
		t.Fatalf("create market failed: %s", resp.Message)
	}

	mkt, _ := pbl.UnmarshalMarket(resp.Payload)

	resp = ccc.MockInvokeWithString(stub, "query", mkt.Id)
	if !ccc.OK(&resp) {
		t.Fatal("query market failed")
	}

	mkt2, _ := pbl.UnmarshalMarket(resp.Payload)

	if !pbl.CmpMarket(mkt, mkt2) {
		t.Fatal("data error")
	}

	mb, existed, err := ccc.GetStateAndCheck(stub, member.Id)
	if err != nil {
		t.Fatal("query user failed")
	}

	if !existed {
		t.Fatal("user not found")
	}

	member2, _ := pbm.Unmarshal(mb)

	if member2.Balance != (member.Balance - mkt.Fund) {
		t.Fatal("member balance failed")
	}

	member.Balance = member.Balance - mkt.Fund
}

func TestSettle(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, "create", "liquidity", member.Id, event.Id, "100.0")
	mkt, _ := pbl.UnmarshalMarket(resp.Payload)

	resp = ccc.MockInvokeWithString(stub, "settle", mkt.Id)
	if !ccc.OK(&resp) {
		t.Fatal("settle failed")
	}

	mkt.Settled = true

	resp = ccc.MockInvokeWithString(stub, "query", mkt.Id)
	mkt1, _ := pbl.UnmarshalMarket(resp.Payload)
	if !mkt1.Settled {
		t.Fatal("settle status failed")
	}

	if !pbl.CmpMarket(mkt, mkt1) {
		t.Fatal("data not match")
	}

	resp = ccc.MockInvokeWithString(stub, "settle", mkt.Id)
	if ccc.OK(&resp) {
		t.Fatal("can not settle a settled market")
	}

}

func TestWrongData(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, "query", "abc")
	if ccc.OK(&resp) {
		t.Fatal("can not query non-existed market")
	}

	resp = ccc.MockInvokeWithString(stub, "query", "a", "b")
	if ccc.OK(&resp) {
		t.Fatal("can not query with wrong parameters")
	}

	resp = ccc.MockInvokeWithString(stub, "settle", "a")
	if ccc.OK(&resp) {
		t.Fatal("can not settle non-existed market")
	}

	resp = ccc.MockInvokeWithString(stub, "settle", "a", "b")
	if ccc.OK(&resp) {
		t.Fatal("can not settle with wrong parameters")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "abc", member.Id, event.Id, "100.0")
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market with wrong flag")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "fund", member.Id, event.Id, "10000000000000.0")
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market without enough balance")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "fund", member.Id, event.Id, "abc")
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market without float64 number")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "liquidity", member.Id, event.Id, "100.0", "A")
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market with wrong parameters")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "fund", "abc", event.Id, "100.0")
	if ccc.OK(&resp) {
		t.Fatal("can not create market with non-existed user")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "liquidity", member.Id, "abc", "100.0")
	if ccc.OK(&resp) {
		t.Fatal("can not create market with non-existed user")
	}

	resp = ccc.MockInvokeWithString(stub, "create", "abc", member.Id, event.Id, "100.0")
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market with wrong flag")
	}

	resp = ccc.MockInvokeWithString(stub, "aaa", "ab")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke with wrong function name")
	}
}
