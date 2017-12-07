package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type manageCC struct {
}

// NewManageChaincode ...
func NewManageChaincode() shim.Chaincode {
	return new(manageCC)
}

var logger = shim.NewLogger("manage")

func (cc *manageCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("init manage chaincode")
	return shim.Success(nil)
}

// Invoke ...
func (cc *manageCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func main() {
	logger.SetLevel(shim.LogDebug)
	logger.Debug("start manage chaincode")
	err := shim.Start(NewManageChaincode())
	if err != nil {
		logger.Errorf("creating manage chaincode failed: %v\n", err)
	}
}
