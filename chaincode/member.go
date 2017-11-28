package main

import (
	ccc "diviner/chaincode/common"

	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type memberCC struct {
}

func NewMemberChaincode() Mychaincode {
	return &memberCC{}
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

// Invoke ...
func (cc *memberCC) Invoke(stub shim.ChaincodeStubInterface, fcn string, args [][]byte) pb.Response {
	if len(args) != 1 {
		return ccc.Errorf("args length for member invoke error: %d", len(args))
	}

	switch fcn {
	case "query":
		return cc.query(stub, string(args[0]))
	case "create":
		return ccc.SetEventAndReturn(stub, ccc.ChaincodeEventID(stub, "member"), cc.create(stub, args[0]))
	case "update":
		return ccc.SetEventAndReturn(stub, ccc.ChaincodeEventID(stub, "member"), cc.update(stub, args[0]))
	}

	return ccc.Errorf("member unknown function: %s", fcn)
}
