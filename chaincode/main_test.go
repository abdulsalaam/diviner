package chaincode

import (
	"diviner/common/csp"
	"os"
	"testing"

	ccc "diviner/chaincode/common"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	query   = []byte("query")
	create  = []byte("create")
	update  = []byte("update")
	approve = []byte("approve")

	priv1 bccsp.Key
	priv2 bccsp.Key
	priv3 bccsp.Key
	//priv5 bccsp.Key

	m1 *pbm.Member
	m2 *pbm.Member
	m3 *pbm.Member
	//m5 *pbm.Member

	memStub *shim.MockStub
	evtStub *shim.MockStub

	title    = "gogogo"
	outcomes = []string{"yes", "no"}
)

func TestMain(m *testing.M) {

	//fmt.Println(m)
	priv1, _ = csp.GeneratePrivateTempKey()
	priv2, _ = csp.GeneratePrivateTempKey()
	priv3, _ = csp.GeneratePrivateTempKey()
	//priv5, _ = csp.GeneratePrivateTempKey()

	m1, _ = pbm.NewMember(priv1, 0.0)
	m2, _ = pbm.NewMember(priv2, 10.0)
	m3, _ = pbm.NewMember(priv3, 1000.0)
	//m5, _ = pbm.NewMember(priv5, 100000.0)

	memStub = ccc.NewMockStub("member", NewMemberChaincode())
	ccc.MockInit(memStub)

	evtStub = ccc.NewMockStub("event", NewEventChaincode())
	evtStub.MockPeerChaincode("member", memStub)
	ccc.MockInit(evtStub)

	//tmp, _ := pbm.Marshal(m5)
	//txid := uuid.New().String()
	//evtStub.MockTransactionStart(txid)
	//evtStub.InvokeChaincode("member", [][]byte{create, tmp}, "")
	//evtStub.MockTransactionEnd(txid)

	r := m.Run()
	os.Exit(r)
}
