package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("test1")

type test1CC struct{}

func (cc *test1CC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	stub.PutState("a", []byte("1000"))
	stub.PutState("b", []byte("1000"))
	return shim.Success(nil)
}

func (cc *test1CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, args := stub.GetFunctionAndParameters()

	switch fcn {
	case "query":
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("args length error: %d", len(args)))
		}

		bytes, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		} else if bytes == nil {
			shim.Error(fmt.Sprintf("%s not found", args[0]))
		}
		return shim.Success(bytes)
	case "update":
		if len(args) != 2 {
			return shim.Error(fmt.Sprintf("args length error: %d", len(args)))
		}

		err := stub.PutState(args[0], []byte(args[1]))
		if err != nil {
			return shim.Error(err.Error())
		}

		msg := fmt.Sprintf("{\"id\":%q, \"value\":%q}", args[0], args[1])
		return shim.Success([]byte(msg))
	}

	return shim.Error(fmt.Sprintf("fcn error: %q", fcn))
}

func main() {
	err := shim.Start(new(test1CC))
	if err != nil {
		logger.Errorf("creating member chaincode failed: %v\n", err)
	}
}
