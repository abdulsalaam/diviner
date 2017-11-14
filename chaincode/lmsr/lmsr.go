package lmsr

import (
	ccc "diviner/chaincode/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type lmsrCC struct{}

// NewLMSRChaincode ...
func NewLMSRChaincode() shim.Chaincode {
	return new(lmsrCC)
}

func (cc *lmsrCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (cc *lmsrCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, args := stub.GetFunctionAndParameters()

	switch fcn {
	case "buy":
		len := len(args)
		if len != 2 {
			return ccc.Errorf("args length error for buy: %v", len)
		}

		//share := args[1]
		//num := args[2]

	case "sell":
	}

	return ccc.Errorf("unknown function: %s", fcn)
}
