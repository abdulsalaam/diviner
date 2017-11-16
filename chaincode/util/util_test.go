package util

import (
	"testing"

	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type testcc struct{}

func (cc *testcc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *testcc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

var (
	stub     = ccc.NewMockStub("test", new(testcc))
	title    = "gogogo"
	outcomes = []string{"yes", "no"}
	balance  = 10000.0
	member   *pbm.Member
	event    *pbl.Event
	market1  *pbl.Market
	market2  *pbl.Market
)

func TestMain(m *testing.M) {
	priv, _ := csp.GeneratePrivateTempKey()
	member, _ = pbm.NewMember(priv, balance)
	event, _ = pbl.NewEvent(member.Id, title, outcomes[0], outcomes[1])
	market1, _ = pbl.NewMarketWithFund(member.Id, event, 100.0)

	market2, _ = pbl.NewMarketWithFund(member.Id, event, 200.0)

	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	ccc.PutMessage(stub, member.Id, member)
	ccc.PutMessage(stub, event.Id, event)
	eid, mid, _ := pbl.SepMarketID(market1.Id)
	ccc.PutMessageWithCompositeKey(stub, market1, pbl.MarketKey, eid, mid)
	eid, mid, _ = pbl.SepMarketID(market2.Id)
	ccc.PutMessageWithCompositeKey(stub, market2, pbl.MarketKey, eid, mid)
	stub.MockTransactionEnd(txid)

	m.Run()
}

func TestFindMaket(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	m, err := FindMarket(stub, market1.Id)
	if err != nil {
		t.Fatalf("find market failed: %v", err)
	}

	if !pbl.CmpMarket(m, market1) {
		t.Fatal("market data not match")
	}
}

func TestPutAndFindMarket(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	event1, _ := pbl.NewEvent(member.Id, "test", "a", "b")
	m1, _ := pbl.NewMarketWithFund(member.Id, event1, 10.0)

	bytes, err := PutMarket(stub, m1)
	if err != nil {
		t.Fatal("put market failed: %v", err)
	}

	m2, _ := pbl.UnmarshalMarket(bytes)

	m3, _ := FindMarket(stub, m1.Id)

	if !pbl.CmpMarket(m1, m2) || !pbl.CmpMarket(m1, m3) {
		t.Fatal("data not match")
	}
}

func TestFindAllMarket(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	lst, err := FindAllMarkets(stub, event.Id)
	if err != nil {
		t.Fatalf("find all markets of event (%s) failed: %v", event.Id, err)
	}

	if len(lst.List) != 2 {
		t.Fatal("length of list failed: %d", len(lst.List))
	}

	if !(pbl.CmpMarket(market1, lst.List[0]) ||
		pbl.CmpMarket(market1, lst.List[1]) ||
		pbl.CmpMarket(market2, lst.List[0]) ||
		pbl.CmpMarket(market2, lst.List[1])) {
		t.Fatal("data not match")
	}

}
