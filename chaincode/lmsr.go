package main

import (
	ccc "diviner/chaincode/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Mychaincode interface {
	Invoke(shim.ChaincodeStubInterface, string, [][]byte) pb.Response
}

var logger = shim.NewLogger("lmsr")

type lmsrCC struct{}

// NewLMSRChaincode ...
func NewLMSRChaincode() shim.Chaincode {
	return new(lmsrCC)
}

func (cc *lmsrCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("init lmsr chaincode")
	return shim.Success(nil)
}

// Invoke ...
func (cc *lmsrCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	logger.Debugf("invoke args = %v\n", args)
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
		mycc = NewEventChaincode()
	case "market":
		mycc = NewMarketChaincode()
	case "tx":
		mycc = NewTxChaincode()
	}

	if mycc != nil {
		return mycc.Invoke(stub, fcn, args[2:])
	}

	return ccc.Errorf("lmsr invoke unknown module: %q", module)
}

func main() {
	logger.SetLevel(shim.LogDebug)
	logger.Debug("start lmsr chaincode")
	err := shim.Start(NewLMSRChaincode())
	if err != nil {
		//fmt.Printf("creating lmsr chaincode failed: %v", err)
		logger.Errorf("creating lmsr chaincode failed: %v\n", err)
	}
}
