package main

import (
	"diviner/common/csp"
	"os"
	"testing"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	ccc "diviner/chaincode/common"
	pbm "diviner/protos/member"
)

var (
	fcnQuery   = []byte("query")
	fcnCreate  = []byte("create")
	fcnUpdate  = []byte("update")
	fcnApprove = []byte("approve")
	fcnSettle  = []byte("settle")

	ccMember = []byte("member")
	ccEvent  = []byte("event")
	ccMarket = []byte("market")
	ccTx     = []byte("tx")

	balance = 100000.0

	stub *shim.MockStub

	priv0 bccsp.Key
	priv1 bccsp.Key
	priv2 bccsp.Key
	priv3 bccsp.Key

	m0 *pbm.Member
	m1 *pbm.Member
	m2 *pbm.Member
	m3 *pbm.Member

	title    = "gogogo"
	outcomes = []string{"yes", "no"}
)

func TestMain(m *testing.M) {
	priv0, _ = csp.GeneratePrivateTempKey()
	priv1, _ = csp.GeneratePrivateTempKey()
	priv2, _ = csp.GeneratePrivateTempKey()
	priv3, _ = csp.GeneratePrivateTempKey()

	m0, _ = pbm.NewMember(priv0, balance)
	m1, _ = pbm.NewMember(priv1, balance)
	m2, _ = pbm.NewMember(priv2, balance)
	m3, _ = pbm.NewMember(priv3, balance)

	stub = ccc.NewMockStub("lmsr_test", NewLMSRChaincode())
	ccc.MockInit(stub)

	os.Exit(m.Run())
}
