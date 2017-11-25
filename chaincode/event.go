package chaincode

import (
	ccc "diviner/chaincode/common"

	pbl "diviner/protos/lmsr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type eventCC struct{}

// NewEventChaincode ...
func NewEventChaincode() shim.Chaincode {
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

	resp := ccc.InvokeChaincodeWithString(stub, "member", "", "query", event.User)
	if !ccc.OK(&resp) {
		return resp
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

func (cc *eventCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *eventCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) != 2 {
		return ccc.Errorf("args length error for event invoke: %v", len(args))
	}

	fcn := string(args[0])
	data := args[1]

	switch fcn {
	case "query":
		return cc.query(stub, string(data))
	case "create":
		return cc.create(stub, data)
	case "approve":
		return cc.approve(stub, string(data))
	}

	return ccc.Errorf("unknown function: %s", fcn)
}
