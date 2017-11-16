package common

import (
	"fmt"
	"strings"

	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	perrors "github.com/pkg/errors"
)

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

// GetStateByCompositeKeyAndCheck ...
func GetStateByCompositeKeyAndCheck(stub shim.ChaincodeStubInterface, name string, keys ...string) (map[string][]byte, bool, error) {
	it, err := stub.GetStateByPartialCompositeKey(name, keys)

	if err != nil {
		return nil, false, err
	}

	defer it.Close()

	if !it.HasNext() {
		return nil, false, nil
	}

	result := make(map[string][]byte)
	for i := 0; it.HasNext(); i++ {
		tmp, err := it.Next()
		if err != nil {
			return nil, false, fmt.Errorf("next error: %v", err)
		}
		result[tmp.Key] = tmp.Value
	}

	return result, true, nil
}

// PutMessageAndReturn
func PutMessage(stub shim.ChaincodeStubInterface, key string, msg proto.Message) ([]byte, error) {
	if bytes, err := proto.Marshal(msg); err != nil {
		return nil, err
	} else {
		return bytes, stub.PutState(key, bytes)
	}
}

// PutMessageWithCompositeKey ...
func PutMessageWithCompositeKey(stub shim.ChaincodeStubInterface, msg proto.Message, name string, keys ...string) ([]byte, error) {
	if key, err := stub.CreateCompositeKey(name, keys); err != nil {
		return nil, err
	} else {
		return PutMessage(stub, key, msg)
	}
}

// PutStateByCompositeKey
func PutStateByCompositeKey(stub shim.ChaincodeStubInterface, value []byte, name string, keys ...string) error {
	if key, err := stub.CreateCompositeKey(name, keys); err != nil {
		return err
	} else {
		return stub.PutState(key, value)
	}
}

// PutStateAndReturn ...
func PutStateAndReturn(stub shim.ChaincodeStubInterface, key string, value, payload []byte) pb.Response {
	err := stub.PutState(key, value)
	if err != nil {
		return Errorf("put key (%s) error: %v", key, value)
	}

	return shim.Success(payload)
}

// PutMessageAndReturn ..
func PutMessageAndReturn(stub shim.ChaincodeStubInterface, key string, msg proto.Message) pb.Response {
	if bytes, err := PutMessage(stub, key, msg); err != nil {
		return Errore(err)
	} else {
		return shim.Success(bytes)
	}
}

// MarshalAndReturn ...
func MarshalAndReturn(msg proto.Message) pb.Response {
	if bytes, err := proto.Marshal(msg); err != nil {
		return Errore(err)
	} else {
		return shim.Success(bytes)
	}
}

// OK ...
func OK(resp *pb.Response) bool {
	return resp.Status == shim.OK
}

// NotOK
func NotOK(resp *pb.Response) bool {
	return !(resp.Status == shim.OK)
}

// Find ...
func Find(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	bytes, existed, err := GetStateAndCheck(stub, id)
	if err != nil {
		return nil, perrors.Errorf("find id (%s) error: %v", id, err)
	} else if !existed {
		return nil, perrors.Errorf("id (%s) is not existed", id)
	}

	return bytes, nil
}

func GetOneValue(x map[string][]byte) []byte {
	if len(x) > 0 {
		for _, v := range x {
			return v
		}
	}
	return nil
}

// FindByPartialCompositeKey ...
func FindByPartialCompositeKey(stub shim.ChaincodeStubInterface, name string, keys ...string) ([]byte, error) {
	result, existed, err := GetStateByCompositeKeyAndCheck(stub, name, keys...)
	if existed {
		return GetOneValue(result), nil
	} else if err == nil {
		return nil, fmt.Errorf("data not found: %s", strings.Join(keys, ""))
	}

	return nil, err
}

// FindAllByPartialCompositeKey ...
func FindAllByPartialCompositeKey(stub shim.ChaincodeStubInterface, name string, keys ...string) (map[string][]byte, error) {
	result, existed, err := GetStateByCompositeKeyAndCheck(stub, name, keys...)
	if existed {
		return result, nil
	} else if err == nil {
		return nil, fmt.Errorf("data not found: %s", strings.Join(keys, ""))
	}

	return nil, err
}
