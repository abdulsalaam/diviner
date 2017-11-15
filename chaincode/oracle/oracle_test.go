package oracle

import (
	"fmt"
	"testing"

	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/google/uuid"
)

var (
	stub     = ccc.NewMockStub("oracle", NewOracleChaincode())
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
	fmt.Println("m1: ", market1.Id)
	market2, _ = pbl.NewMarketWithFund(member.Id, event, 200.0)

	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	ccc.PutMessage(stub, member.Id, member)
	ccc.PutMessage(stub, event.Id, event)
	ccc.PutMessageWithCompositeKey(stub, market1, pbl.MarketKey, market1.Id)
	ccc.PutMessageWithCompositeKey(stub, market2, pbl.MarketKey, market2.Id)
	stub.MockTransactionEnd(txid)
	m.Run()
}

func TestMarkets(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, "markets", market1.Id)
	if !ccc.OK(&resp) {
		t.Fatalf("invoke markets failed: %s", resp.Message)
	}
}
