package chaincode

import (
	ccc "diviner/chaincode/common"

	pbl "diviner/protos/lmsr"

	"testing"
)

func TestCreate(t *testing.T) {
	evt, _ := pbl.NewEvent(m5.Id, title, outcomes[0], outcomes[1])
	bytes, _ := pbl.MarshalEvent(evt)

	resp := ccc.MockInvoke(evtStub, []byte(create), bytes)

	if !ccc.OK(&resp) {
		t.Fatalf("create event failed: %s", resp.Message)
	}

	event, err := pbl.UnmarshalEvent(resp.Payload)
	if err != nil {
		t.Fatal("data structure error: %v", err)
	}

	if event.User != m5.Id {
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

	resp = ccc.MockInvokeWithString(evtStub, string(query), event.Id)
	if !ccc.OK(&resp) {
		t.Fatal("query event failed")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(query), "abc")
	if ccc.OK(&resp) {
		t.Fatal("can not query non-existed event")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(approve), event.Id)
	if !ccc.OK(&resp) {
		t.Fatal("approve failed")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(approve), event.Id)
	if ccc.OK(&resp) {
		t.Fatal("can not approve an approved event")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(approve), "abc")
	if ccc.OK(&resp) {
		t.Fatal("can not approve an non-existed event")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(query), event.Id)
	if !ccc.OK(&resp) {
		t.Fatal("query event failed")
	}

	event, _ = pbl.UnmarshalEvent(resp.Payload)
	if !event.Approved {
		t.Fatal("event is not approved")
	}

}

func TestWrongParameter(t *testing.T) {
	resp := ccc.MockInvokeWithString(evtStub, string(query))
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(create), "a", "b", "c", "d")
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(create), "a", "b", "c")
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(evtStub, string(approve), "a", "b", "c")
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvokeWithString(evtStub, "test", "a", "b", "c")
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}
}
