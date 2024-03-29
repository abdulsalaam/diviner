package common

import (
	"diviner/common/cast"
	"fmt"
	"strings"

	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	perrors "github.com/pkg/errors"
)

// Errorf ...
func Errorf(format string, a ...interface{}) pb.Response {
	return shim.Error(fmt.Sprintf(format, a...))
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

// PutMessage ...
func PutMessage(stub shim.ChaincodeStubInterface, key string, msg proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return bytes, stub.PutState(key, bytes)
}

// PutMessageWithCompositeKey ...
func PutMessageWithCompositeKey(stub shim.ChaincodeStubInterface, msg proto.Message, name string, keys ...string) ([]byte, error) {
	key, err := stub.CreateCompositeKey(name, keys)
	if err != nil {
		return nil, err
	}
	return PutMessage(stub, key, msg)
}

// PutStateByCompositeKey ...
func PutStateByCompositeKey(stub shim.ChaincodeStubInterface, value []byte, name string, keys ...string) error {
	key, err := stub.CreateCompositeKey(name, keys)

	if err != nil {
		return err
	}
	return stub.PutState(key, value)
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
	bytes, err := PutMessage(stub, key, msg)
	if err != nil {
		return Errore(err)
	}
	return shim.Success(bytes)
}

// MarshalAndReturn ...
func MarshalAndReturn(msg proto.Message) pb.Response {
	bytes, err := proto.Marshal(msg)

	if err != nil {
		return Errore(err)
	}
	return shim.Success(bytes)
}

// OK ...
func OK(resp *pb.Response) bool {
	return resp.Status == shim.OK
}

// NotOK ...
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

// GetOneValue ...
func GetOneValue(x map[string][]byte) []byte {
	if x == nil || len(x) > 0 {
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

// SetEventAndReturn ...
func SetEventAndReturn(stub shim.ChaincodeStubInterface, name string, resp pb.Response) pb.Response {
	if OK(&resp) && resp.Payload != nil {
		evtID := chaincodeEventID(stub, name)

		if err := stub.SetEvent(evtID, resp.Payload); err != nil {
			return Errorf("set event %s error: %v", name, err)
		}
	}

	return resp
}

func chaincodeEventID(stub shim.ChaincodeStubInterface, name string) string {
	return name + stub.GetTxID()
}

// InvokeChaincodeWithString ...
func InvokeChaincodeWithString(stub shim.ChaincodeStubInterface, chaincodeName, channel string, args ...string) pb.Response {
	tmp := cast.StringsToByteArray(args...)
	return stub.InvokeChaincode(chaincodeName, tmp, channel)
}
