package oracle

import (
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
)

func TestMain(m *testing.M) {
	priv, _ := csp.GeneratePrivateTempKey()
	member, _ := pbm.NewMember(priv, balance)
	event, _ := pbl.NewEvent(member.Id, title, outcomes[0], outcomes[1])
	market1, _ := pbl.NewMarketWithFund(member.Id, event, 100.0)
	market2, _ := pbl.NewMarketWithFund(member.Id, event, 200.0)

	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	ccc.PutMessage(stub, member.Id, member)
	ccc.PutMessage(stub, event.Id, event)
	stub.MockTransactionEnd(txid)
	m.Run()
}
