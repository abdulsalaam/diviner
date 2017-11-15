package oracle

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type oracleCC struct{}

// NewOracleChaincode ...
func NewOracleChaincode() shim.Chaincode {
	return new(oracleCC)
}

func (cc *oracleCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *oracleCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func main() {
	err := shim.Start(NewOracleChaincode())

	if err != nil {
		fmt.Printf("creating oracle chaincode failed: %v", err)
	}
}
