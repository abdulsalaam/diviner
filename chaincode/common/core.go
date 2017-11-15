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

// PutMessageAndReturn
func PutMessage(stub shim.ChaincodeStubInterface, key string, msg proto.Message) error {
	if bytes, err := proto.Marshal(msg); err != nil {
		return err
	} else {
		return stub.PutState(key, bytes)
	}
}

func PutMessageWithCompositeKey(stub shim.ChaincodeStubInterface, msg proto.Message, name string, keys ...string) error {
	if key, err := stub.CreateCompositeKey(name, keys); err != nil {
		return err
	} else {
		fmt.Println("put ckey: ", key)
		a, b, _ := stub.SplitCompositeKey(key)
		fmt.Println("split ckey: ", a, b)
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

// OK ...
func OK(resp *pb.Response) bool {
	return resp.Status == shim.OK
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

// FindByPartialCompositeKey ...
func FindByPartialCompositeKey(stub shim.ChaincodeStubInterface, name string, keys ...string) ([]byte, error) {
	it, err := stub.GetStateByPartialCompositeKey(name, keys)

	if err != nil {
		return nil, err
	}

	defer it.Close()

	if !it.HasNext() {
		return nil, fmt.Errorf("data not found", strings.Join(keys, ""))
	}

	result, err := it.Next()
	if err != nil {
		return nil, fmt.Errorf("next error: %v", err)
	}

	return result.Value, nil
}

// FindAllByPartialCompositeKey ...
func FindAllByPartialCompositeKey(stub shim.ChaincodeStubInterface, name string, keys ...string) (map[string][]byte, error) {
	ckey, _ := stub.CreateCompositeKey(name, keys)
	fmt.Println("find by ckey: ", ckey)
	it, err := stub.GetStateByPartialCompositeKey(name, keys)

	if err != nil {
		return nil, err
	}

	defer it.Close()

	if !it.HasNext() {
		return nil, fmt.Errorf("data not found: %s", strings.Join(keys, ""))
	}

	result := make(map[string][]byte)
	for i := 0; it.HasNext(); i++ {
		tmp, err := it.Next()
		if err != nil {
			return nil, fmt.Errorf("next error: %v", err)
		}
		result[tmp.Key] = tmp.Value
	}

	return result, nil
}
