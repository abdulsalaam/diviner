package member

import (
	"fmt"

	ccc "diviner/chaincode/common"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// MemberCC ...
type memberCC struct{}

func (mem *memberCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, existed, err := ccc.GetStateAndCheck(stub, id)
	if err != nil {
		return ccc.Errorf("get member with id (%s) error: %v", id, err)
	} else if !existed {
		return ccc.Errorf("id (%s) is not existed", id)
	}

	_, err = pbm.Unmarshal(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal failed %v", err)
	}

	return shim.Success(bytes)
}

func (mem *memberCC) create(stub shim.ChaincodeStubInterface, data []byte) pb.Response {
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

	return ccc.PutStateAndReturn(stub, member.Id, data, nil)
}

func (mem *memberCC) update(stub shim.ChaincodeStubInterface, data []byte) pb.Response {
	member, err := pbm.Unmarshal(data)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	_, existed, err := ccc.GetStateAndCheck(stub, member.Id)
	if err != nil {
		return ccc.Errorf("get member with id (%s) error: %v", member.Id, err)
	} else if !existed {
		return ccc.Errorf("id (%s) is not existed", member.Id)
	}

	return ccc.PutStateAndReturn(stub, member.Id, data, nil)
}

// Init ...
func (mem *memberCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (mem *memberCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) != 2 {
		return ccc.Errorf("args length error: %v", len(args))
	}

	fcn := string(args[0])
	switch fcn {
	case "query":
		return mem.query(stub, string(args[1]))
	case "create":
		return mem.create(stub, args[1])
	case "update":
		return mem.update(stub, args[1])
	}

	return ccc.Errorf("unknown function: %s", fcn)
}

func main() {
	err := shim.Start(new(memberCC))

	if err != nil {
		fmt.Printf("creating member chaincode failed: %v", err)
	}
}
