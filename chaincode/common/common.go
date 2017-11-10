package common

import (
	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// StringsToByteArray ...
func StringsToByteArray(args ...string) [][]byte {
	result := make([][]byte, len(args))

	for i, x := range args {
		result[i] = []byte(x)
	}

	return result
}

// NewMockStub ...
func NewMockStub(name string, cc shim.Chaincode) *shim.MockStub {
	return shim.NewMockStub(name, cc)
}

// MockInit ...
func MockInit(stub *shim.MockStub, args ...[]byte) pb.Response {
	return stub.MockInit(uuid.New().String(), args)
}

// MockInitWithString ...
func MockInitWithString(stub *shim.MockStub, args ...string) pb.Response {
	tmp := StringsToByteArray(args...)
	return MockInit(stub, tmp...)
}

// MockInvoke ...
func MockInvoke(stub shim.MockStub, args ...[]byte) pb.Response {
	return stub.MockInvoke(uuid.New().String(), args)
}

// MockInvokeWithString ...
func MockInvokeWithString(stub *shim.MockStub, args ...string) pb.Response {
	tmp := StringsToByteArray(args...)
	return stub.MockInvoke(uuid.New().String(), tmp)
}
