package event

import (
	ccc "diviner/chaincode/common"
	"diviner/chaincode/member"
	"diviner/common/csp"
	pbc "diviner/protos/common"
	pbmk "diviner/protos/market"
	pbm "diviner/protos/member"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	stub    *shim.MockStub
	memStub *shim.MockStub

	priv1 bccsp.Key
	priv2 bccsp.Key
	m1    *pbm.Member
	m2    *pbm.Member
)

func TestCreate(t *testing.T) {
	curr := time.Now()

	event, err := pbmk.NewEvent(m1.Address, "gogo", curr.AddDate(1, 0, 0), "yes", "no")
	if err != nil {
		t.Fatal(err)
	}

	eb, _ := pbmk.MarshalEvent(event)
	valid, _ := pbc.NewVerification(priv1, eb)
	vb, _ := pbc.Marshal(valid)

	resp := ccc.MockInvoke(stub, []byte("create"), eb, vb)
	if ccc.NotOK(&resp) {
		t.Fatalf("create event failed: %s", resp.Message)
	}

	resp = ccc.MockInvoke(stub, []byte("create"), eb, vb)
	if resp.Status != ccc.CONFLICT {
		t.Fatalf("can not create an existed event. code: %d", resp.Status)
	}

	event, _ = pbmk.NewEvent(m2.Address, "gogo", curr.AddDate(1, 0, 0), "yes", "no")
	eb, _ = pbmk.MarshalEvent(event)
	valid, _ = pbc.NewVerification(priv2, eb)
	vb, _ = pbc.Marshal(valid)
	resp = ccc.MockInvoke(stub, []byte("create"), eb, vb)

	if resp.Status != ccc.NOTFOUND {
		t.Fatalf("can not create an event with non-existed member: %d, %s", resp.Status, resp.Message)
	}

	event, _ = pbmk.NewEvent(m2.Address, "gogo", curr.AddDate(1, 0, 0), "yes", "no")
	eb, _ = pbmk.MarshalEvent(event)
	valid, _ = pbc.NewVerification(priv1, eb)
	vb, _ = pbc.Marshal(valid)
	resp = ccc.MockInvoke(stub, []byte("create"), eb, vb)
	//t.Logf("%d, %s\n", resp.Status, resp.Message)
	if resp.Status != ccc.BADREQUEST {
		t.Fatalf("can not create an event with wrong creator: %d, %s", resp.Status, resp.Message)
	}

	priv2pub, _ := csp.GetPublicKeyBytes(priv2)
	valid.PublicKey = priv2pub
	vb, _ = pbc.Marshal(valid)
	resp = ccc.MockInvoke(stub, []byte("create"), eb, vb)
	//t.Logf("%d, %s\n", resp.Status, resp.Message)
	if resp.Status != ccc.UNAUTHORIZED {
		t.Fatalf("can not create an event without authorization: %d, %s", resp.Status, resp.Message)
	}

}

func TestQuery(t *testing.T) {
	curr := time.Now()

	event, err := pbmk.NewEvent(m1.Address, "gogo", curr.AddDate(1, 0, 0), "yes", "no")
	if err != nil {
		t.Fatal(err)
	}

	eb, _ := pbmk.MarshalEvent(event)
	valid, _ := pbc.NewVerification(priv1, eb)
	vb, _ := pbc.Marshal(valid)

	resp := ccc.MockInvoke(stub, []byte("create"), eb, vb)
	if ccc.NotOK(&resp) {
		t.Fatalf("create event failed: %s", resp.Message)
	}

	resp = ccc.MockInvoke(stub, []byte("query"), []byte(event.Id))
	if ccc.NotOK(&resp) {
		t.Fatalf("query event failed: %d, %s", resp.Status, resp.Message)
	}

	evt2, err := pbmk.UnmarshalEvent(resp.Payload)
	//t.Logf("%s, %s\n", evt2.String(), event.String())
	if evt2.String() != event.String() {
		t.Fatalf("data not match: %s, %s", evt2.String(), event.String())
	}
}

func TestMain(m *testing.M) {
	stub = ccc.NewMockStub("event", NewEventChaincode())
	resp := ccc.MockInit(stub)
	if ccc.NotOK(&resp) {
		fmt.Printf("init event chaincode failed, %d, %s\n", resp.Status, resp.Message)
		os.Exit(-1)
	}

	memStub = ccc.NewMockStub("member", member.NewMemberChaincode())
	stub.MockPeerChaincode("member", memStub)

	priv1, _ = csp.GeneratePrivateTempKey()
	m1, _ = pbm.NewMemberWithPrivateKey(priv1, 0.0)
	b1, _ := pbm.Marshal(m1)
	v1, _ := pbc.NewVerification(priv1, b1)
	vb1, _ := pbc.Marshal(v1)
	resp = ccc.MockInvoke(memStub, []byte("register"), b1, vb1)
	if ccc.NotOK(&resp) {
		fmt.Println("create member failed")
		os.Exit(-1)
	}
	m1, _ = pbm.Unmarshal(resp.Payload)

	priv2, _ = csp.GeneratePrivateTempKey()
	m2, _ = pbm.NewMemberWithPrivateKey(priv2, 0.0)

	os.Exit(m.Run())
}
