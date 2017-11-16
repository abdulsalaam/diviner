package lmsr

import (
	"testing"

	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/google/uuid"
)

var (
	stub     = ccc.NewMockStub("oracle", NewLMSRChaincode())
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

func TestMarkets(t *testing.T) {

	resp := ccc.MockInvokeWithString(stub, "markets", "a", "b")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke with wrong data")
	}

	resp = ccc.MockInvokeWithString(stub, "markets", event.Id)
	if !ccc.OK(&resp) {
		t.Fatalf("invoke markets failed: %s", resp.Message)
	}

	markets, _ := pbl.UnmarshalMarkets(resp.Payload)
	if len(markets.List) != 2 {
		t.Fatal("list length failed")
	}

	if !(pbl.CmpMarket(market1, markets.List[0]) ||
		pbl.CmpMarket(market1, markets.List[1]) ||
		pbl.CmpMarket(market2, markets.List[0]) ||
		pbl.CmpMarket(market2, markets.List[1])) {
		t.Fatal("data not match")
	}

}
