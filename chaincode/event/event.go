package main

import (
	"fmt"

	ccc "diviner/chaincode/common"

	pbe "diviner/protos/lmsr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type eventCC struct{}

// NewEventChaincode ...
func NewEventChaincode() shim.Chaincode {
	return new(eventCC)
}

func (evt *eventCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {

	bytes, err := ccc.Find(stub, id)
	if err != nil {
		return ccc.Errore(err)
	}

	_, err = pbe.UnmarshalEvent(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal event error: %v", err)
	}

	return shim.Success(bytes)
}

func (evt *eventCC) create(stub shim.ChaincodeStubInterface, user, title string, outcomes []string) pb.Response {
	if len(outcomes) < 2 {
		return ccc.Errorf("length of outcomes must be larer than 1: %v", len(outcomes))
	}

	if title == "" {
		return ccc.Errorf("title is empty")
	}

	//_, err := ccc.Find(stub, user)
	//if err != nil {
	//	return ccc.Errore(err)
	//}
	resp := stub.InvokeChaincode("member", [][]byte{[]byte("query"), []byte(user)}, "")
	if resp.Status != shim.OK {
		return resp
	}

	event, err := pbe.NewEvent(user, title, outcomes...)
	if err != nil {
		return ccc.Errorf("create event error: %v", err)
	}

	return ccc.PutMessageAndReturn(stub, event.Id, event)
}

func (evt *eventCC) approve(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, err := ccc.Find(stub, id)
	if err != nil {
		return ccc.Errore(err)
	}

	event, err := pbe.UnmarshalEvent(bytes)
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

func (evt *eventCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (evt *eventCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, params := stub.GetFunctionAndParameters()
	switch fcn {
	case "query":
		if len(params) != 1 {
			return ccc.Errorf("args length error for query: %v", len(params))
		}
		return evt.query(stub, params[0])

	case "create":
		if len(params) < 4 {
			return ccc.Errorf("args length error for create: %v", len(params))
		}
		return ccc.SetEventAndReturn(stub, "testEvent", evt.create(stub, params[0], params[1], params[2:]))
	case "approve":
		if len(params) != 1 {
			return ccc.Errorf("args length error for query: %v", len(params))
		}
		return evt.approve(stub, params[0])
	}

	return ccc.Errorf("unknown function: %s", fcn)
}

func main() {
	err := shim.Start(NewEventChaincode())

	if err != nil {
		fmt.Printf("creating event chaincode failed: %v", err)
	}
}
