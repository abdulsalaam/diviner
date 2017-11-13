package event

import (
	ccc "diviner/chaincode/common"
	pbe "diviner/protos/lmsr"
	"testing"

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
	stub = ccc.NewMockStub("event", new(eventCC))
	ccc.MockInit(stub)
	m.Run()
}

func TestCreate(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, create, user, title, outcomes[0], outcomes[1])
	if resp.Status != shim.OK {
		t.Fatal("create event failed: %s", resp.Message)
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
	if resp.Status != shim.OK {
		t.Fatal("query event failed")
	}

	resp = ccc.MockInvokeWithString(stub, query, "abc")
	if resp.Status == shim.OK {
		t.Fatal("can not query non-existed event")
	}

	resp = ccc.MockInvokeWithString(stub, approve, event.Id)
	if resp.Status != shim.OK {
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
	if resp.Status != shim.OK {
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
