package common

import (
	"diviner/common/cast"
	"fmt"

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

// Errorf ...
func Errorf(format string, a ...interface{}) pb.Response {
	return shim.Error(fmt.Sprintf(format, a))
}

// Errore ...
func Errore(err error) pb.Response {
	return shim.Error(err.Error())
}

// GetStateAndCheck ...
func GetStateAndCheck(stub shim.ChaincodeStubInterface, key string) ([]byte, bool, error) {
	bytes, err := stub.GetState(key)

	return bytes, bytes != nil && err == nil, err
}

// PutStateAndReturn ...
func PutStateAndReturn(stub shim.ChaincodeStubInterface, key string, value, payload []byte) pb.Response {
	err := stub.PutState(key, value)
	if err != nil {
		return Errorf("put key (%s) error: %v", key, value)
	}

	return shim.Success(payload)
}

// OK ...
func OK(resp *pb.Response) bool {
	return resp.Status == shim.OK
}
