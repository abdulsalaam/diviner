package lmsr

import (
	ccc "diviner/chaincode/common"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type lmsrCC struct{}

// NewLMSRChaincode ...
func NewLMSRChaincode() shim.Chaincode {
	return new(lmsrCC)
}

func (cc *lmsrCC) buy(stub shim.ChaincodeStubInterface, user, share string, volume float64) pb.Response {

}

func (cc *lmsrCC) sell(stub shim.ChaincodeStubInterface, user, asset string, volume float64) pb.Response {
}

func (cc *lmsrCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (cc *lmsrCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, args := stub.GetFunctionAndParameters()
	len := len(args)

	switch fcn {
	case "buy":
		if len != 3 {
			return ccc.Errorf("args length error for buying: %v", len)
		}

		volume, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return ccc.Errorf("volume must be float64: %s", args[1])
		}
		user := args[0]
		share := args[1]
		return cc.buy(stub, user, share, volume)

	case "sell":
		if len != 3 {
			return ccc.Errorf("args length error for selling: %v", len)
		}
		volume, err := strconv.ParseFloat(args[2], 64)
		user := args[0]
		asset := args[1]
		return cc.sell(stub, user, asset, volume)

	}

	return ccc.Errorf("unknown function: %s", fcn)
}
