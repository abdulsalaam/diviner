package main

import (
	"fmt"

	ccc "diviner/chaincode/common"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// MemberCC ...
type memberCC struct{}

// NewMemberChaincode ...
func NewMemberChaincode() shim.Chaincode {
	return new(memberCC)
}

func (cc *memberCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, err := ccc.Find(stub, id)
	if err != nil {
		return ccc.Errore(err)
	}

	_, err = pbm.Unmarshal(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal failed %v", err)
	}

	return shim.Success(bytes)
}

func (cc *memberCC) create(stub shim.ChaincodeStubInterface, data []byte) pb.Response {
	member, err := pbm.Unmarshal(data)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	_, existed, err := ccc.GetStateAndCheck(stub, member.Id)
	if err != nil {
		return ccc.Errorf("get member with id (%s) error: %v", member.Id, err)
	} else if existed {
		return ccc.Errorf("id (%s) is existed", member.Id)
	}

	return ccc.PutStateAndReturn(stub, member.Id, data, data)
}

func (cc *memberCC) update(stub shim.ChaincodeStubInterface, data []byte) pb.Response {
	member, err := pbm.Unmarshal(data)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	_, err = ccc.Find(stub, member.Id)
	if err != nil {
		return ccc.Errore(err)
	}

	return ccc.PutStateAndReturn(stub, member.Id, data, data)
}

// Init ...
func (cc *memberCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (cc *memberCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) != 2 {
		return ccc.Errorf("args length error: %v", len(args))
	}

	fcn := string(args[0])
	switch fcn {
	case "query":
		return cc.query(stub, string(args[1]))
	case "create":
		if err := stub.SetEvent("testEvent", []byte("testEvent")); err != nil {
			return shim.Error(err.Error())
		}
		return cc.create(stub, args[1])
	case "update":
		return ccc.SetEventAndReturn(stub, "testEvent", cc.update(stub, args[1]))
	}

	return ccc.Errorf("unknown function: %s", fcn)
}

func main() {
	err := shim.Start(NewMemberChaincode())

	if err != nil {
		fmt.Printf("creating member chaincode failed: %v", err)
	}
}
