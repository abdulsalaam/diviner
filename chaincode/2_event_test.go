package main

import (
	ccc "diviner/chaincode/common"

	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"testing"
)

func TestCreateEvent(t *testing.T) {
	if _, ok, err := ccc.GetStateAndCheck(stub, m0.Id); err != nil {
		t.Fatal(err)
	} else if !ok {
		b0, _ := pbm.Marshal(m0)
		ccc.MockInvoke(stub, fcnCreate, ccMember, b0)
	}

	evt, _ := pbl.NewEvent(m0.Id, title, outcomes[0], outcomes[1])
	bytes, _ := pbl.MarshalEvent(evt)

	resp := ccc.MockInvoke(stub, fcnCreate, ccEvent, bytes)

	if !ccc.OK(&resp) {
		t.Fatalf("create event failed: %s", resp.Message)
	}

	event, err := pbl.UnmarshalEvent(resp.Payload)
	if err != nil {
		t.Fatal("data structure error: %v", err)
	}

	if event.User != m0.Id {
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

	resp = ccc.MockInvoke(stub, fcnQuery, ccEvent, []byte(event.Id))
	if !ccc.OK(&resp) {
		t.Fatal("query event failed")
	}

	resp = ccc.MockInvoke(stub, fcnQuery, ccEvent, []byte("abc"))
	if ccc.OK(&resp) {
		t.Fatal("can not query non-existed event")
	}

	resp = ccc.MockInvoke(stub, fcnApprove, ccEvent, []byte(event.Id))
	if !ccc.OK(&resp) {
		t.Fatal("approve failed")
	}

	resp = ccc.MockInvoke(stub, fcnApprove, ccEvent, []byte(event.Id))
	if ccc.OK(&resp) {
		t.Fatal("can not approve an approved event")
	}

	resp = ccc.MockInvoke(stub, fcnApprove, ccEvent, []byte("abc"))
	if ccc.OK(&resp) {
		t.Fatal("can not approve an non-existed event")
	}

	resp = ccc.MockInvoke(stub, fcnQuery, ccEvent, []byte(event.Id))
	if !ccc.OK(&resp) {
		t.Fatal("query event failed")
	}

	event, _ = pbl.UnmarshalEvent(resp.Payload)
	if !event.Approved {
		t.Fatal("event is not approved")
	}

}

func TestWrongParameter(t *testing.T) {
	resp := ccc.MockInvoke(stub, fcnQuery, ccEvent)
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvoke(stub, fcnCreate, ccEvent, []byte("a"), []byte("b"), []byte("c"), []byte("d"))
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvoke(stub, fcnCreate, ccEvent, []byte("a"), []byte("b"), []byte("c"))
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}

	resp = ccc.MockInvoke(stub, fcnApprove, ccEvent, []byte("a"), []byte("b"), []byte("c"))
	if ccc.OK(&resp) {
		t.Error("wrong parameters must be failed")
	}
}
