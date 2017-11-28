package main

import (
	ccc "diviner/chaincode/common"

	pbl "diviner/protos/lmsr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type eventCC struct {
}

// NewEventChaincode ...
func NewEventChaincode() Mychaincode {
	return new(eventCC)
}

func (cc *eventCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {

	bytes, err := ccc.Find(stub, id)
	if err != nil {
		return ccc.Errore(err)
	}

	_, err = pbl.UnmarshalEvent(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal event error: %v", err)
	}

	return shim.Success(bytes)
}

func (cc *eventCC) create(stub shim.ChaincodeStubInterface, data []byte) pb.Response {
	event, err := pbl.UnmarshalEvent(data)
	if err != nil {
		ccc.Errorf("unmarshal event failure: %v", err)
	}

	if event.Id == "" || event.Title == "" || event.User == "" {
		ccc.Errorf("event id, title or user is empty")
	}

	if event.Approved {
		ccc.Errorf("event is approved")
	}

	if len(event.Outcomes) < 2 {
		return ccc.Errorf("length of outcomes must be larer than 1: %v", len(event.Outcomes))
	}

	_, err = ccc.Find(stub, event.User)
	if err != nil {
		return ccc.Errorf("find event creator error: %v", err)
	}

	return ccc.PutMessageAndReturn(stub, event.Id, event)
}

func (cc *eventCC) approve(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, err := ccc.Find(stub, id)
	if err != nil {
		return ccc.Errore(err)
	}

	event, err := pbl.UnmarshalEvent(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal event error: %v", err)
	}

	if event.Approved {
		return ccc.Errorf("event has been approved")
	}

	event.Approved = true

	if _, err := ccc.PutMessage(stub, event.Id, event); err != nil {
		return ccc.Errorf("put event error: %v", err)
	}

	return shim.Success(nil)
}

func (cc *eventCC) Invoke(stub shim.ChaincodeStubInterface, fcn string, args [][]byte) pb.Response {
	if len(args) != 1 {
		return ccc.Errorf("args length for event invoke error: %d", len(args))
	}

	switch fcn {
	case "query":
		return cc.query(stub, string(args[0]))
	case "create":
		return ccc.SetEventAndReturn(stub, ccc.ChaincodeEventID(stub, "event"), cc.create(stub, args[0]))
	case "approve":
		return ccc.SetEventAndReturn(stub, ccc.ChaincodeEventID(stub, "event"), cc.approve(stub, string(args[0])))
	}

	return ccc.Errorf("event unknown function: %s", fcn)
}
