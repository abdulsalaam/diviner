package chaincode

import (
	ccc "diviner/chaincode/common"
	pbm "diviner/protos/member"
	"fmt"

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
	logger := shim.NewLogger("Member")
	logger.SetLevel(shim.LogDebug)
	logger.Warning("Member init")
	fmt.Println("Member Init fmt")
	return shim.Success(nil)
}

// Invoke ...
func (cc *memberCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger := shim.NewLogger("Member")
	args := stub.GetArgs()
	logger.Warningf("Member Invoke: ", args)
	fmt.Println("Member Invoke fmt: ", args)

	if len(args) != 2 {
		return ccc.Errorf("member args length error: %v", len(args))
	}

	fcn := string(args[0])
	switch fcn {
	case "query":
		return cc.query(stub, string(args[1]))
	case "create":
		return cc.create(stub, args[1])
	case "update":
		return cc.update(stub, args[1])
	}

	return ccc.Errorf("unknown function: %s", fcn)
}
