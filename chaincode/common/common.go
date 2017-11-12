package common

import (
	"diviner/common/cast"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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
	tmp := cast.StringsToByteArray(args...)
	return MockInit(stub, tmp...)
}

// MockInvoke ...
func MockInvoke(stub *shim.MockStub, args ...[]byte) pb.Response {
	return stub.MockInvoke(uuid.New().String(), args)
}

// MockInvokeWithString ...
func MockInvokeWithString(stub *shim.MockStub, args ...string) pb.Response {
	tmp := cast.StringsToByteArray(args...)
	return MockInvoke(stub, tmp...)
}
