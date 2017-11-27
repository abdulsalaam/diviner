package main

import (
	ccc "diviner/chaincode/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Mychaincode interface {
	Invoke(shim.ChaincodeStubInterface, string, [][]byte) pb.Response
}

type lmsrCC struct{}

func (cc *lmsrCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (cc *lmsrCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()

	if len(args) < 2 {
		ccc.Errorf("lmsr invoke args length error: %v", len(args))
	}

	fcn := string(args[0])
	module := string(args[1])
	var mycc Mychaincode
	switch module {
	case "member":
		mycc = NewMemberChaincode()
	case "event":
	case "market":
	case "tx":
	}

	if mycc != nil {
		return mycc.Invoke(stub, fcn, args[2:])
	}

	return ccc.Errorf("lmsr invoke unknown module: %q", module)
}

func main() {

}
