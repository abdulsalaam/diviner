package event

import (
	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbe "diviner/protos/lmsr"
	pbm "diviner/protos/member"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	query    = "query"
	create   = "create"
	approve  = "approve"
	stub     *shim.MockStub
	user     = "user1"
	title    = "gogogo"
	outcomes = []string{"yes", "no"}
)

func TestMain(m *testing.M) {
	stub = ccc.NewMockStub("event", NewEventChaincode())
	ccc.MockInit(stub, nil)

	priv, _ := csp.GeneratePrivateTempKey()
	m1, _ := pbm.NewMember(priv, 0.0)
	b1, _ := pbm.Marshal(m1)

	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	stub.PutState(m1.Id, b1)
	stub.MockTransactionEnd(txid)

	user = m1.Id

	m.Run()
}

func TestCreate(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, create, user, title, outcomes[0], outcomes[1])
	if !ccc.OK(&resp) {
		t.Fatalf("create event failed: %s", resp.Message)
	}

	event, err := pbe.UnmarshalEvent(resp.Payload)
	if err != nil {
		t.Fatal("data structure error: %v", err)
	}

	if event.User != user {
		t.Fatal("user not match")
	}

	if event.Title != title {
		t.Fatal("title not match")
	}

	if len(outcomes) != len(event.Outcomes) {
		t.Fatal("length of outcomes not match")
	}

	for i, x := range outcomes {
		if x != event.Outcomes[i].Title {
			t.Fatal("outcome title not match at %d", i)
		}
	}

	resp = ccc.MockInvokeWithString(stub, query, event.Id)
	if !ccc.OK(&resp) {
		t.Fatal("query event failed")
	}

	resp = ccc.MockInvokeWithString(stub, query, "abc")
	if ccc.OK(&resp) {
		t.Fatal("can not query non-existed event")
	}

	resp = ccc.MockInvokeWithString(stub, approve, event.Id)
	if !ccc.OK(&resp) {
		t.Fatal("approve failed")
	}

	resp = ccc.MockInvokeWithString(stub, approve, event.Id)
	if resp.Status == shim.OK {
		t.Fatal("can not approve an approved event")
	}

	resp = ccc.MockInvokeWithString(stub, approve, "abc")
	if resp.Status == shim.OK {
		t.Fatal("can not approve an non-existed event")
	}

	resp = ccc.MockInvokeWithString(stub, query, event.Id)
	if !ccc.OK(&resp) {
		t.Fatal("query event failed")
	}

	event, _ = pbe.UnmarshalEvent(resp.Payload)
	if !event.Approved {
		t.Fatal("event is not approved")
	}

}

func TestWrongParameter(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, query)
	if resp.Status == shim.OK {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(stub, create, "a", "b", "c", "d")
	if resp.Status == shim.OK {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(stub, create, "a", "b", "c")
	if resp.Status == shim.OK {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(stub, approve, "a", "b", "c")
	if resp.Status == shim.OK {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(stub, "test", "a", "b", "c")
	if resp.Status == shim.OK {
		t.Error("wrong parameters must be failed")
	}
}
