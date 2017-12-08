package common

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (

	// BADREQUEST ...
	BADREQUEST = shim.ERRORTHRESHOLD

	// UNAUTHORIZED ...
	UNAUTHORIZED = 401

	// FORBIDDEN ...
	FORBIDDEN = 403

	// NOTFOUND ...
	NOTFOUND = 404

	// METHODNOTALLOW ...
	METHODNOTALLOW = 405

	// NOTACCEPTABLE ...
	NOTACCEPTABLE = 406

	// CONFLICT ...
	CONFLICT = 409

	// INTERNALSERVERERROR ...
	INTERNALSERVERERROR = shim.ERROR

	// NOTIMPLEMENTED ...
	NOTIMPLEMENTED = 501
)

// Error ...
func Error(status int32, format string, a ...interface{}) pb.Response {
	msg := fmt.Sprintf(format, a...)

	return pb.Response{
		Status:  status,
		Message: msg,
	}
}

// Errorf ...
func Errorf(format string, a ...interface{}) pb.Response {
	return shim.Error(fmt.Sprintf(format, a...))
}

// Errore ...
func Errore(err error) pb.Response {
	return shim.Error(err.Error())
}

// BadRequest ...
func BadRequest(format string, a ...interface{}) pb.Response {
	return Error(BADREQUEST, format, a...)
}

// Unauthorized ...
func Unauthorized() pb.Response {
	return Error(UNAUTHORIZED, "validation error")
}

// NotFound ...
func NotFound(key string) pb.Response {
	return Error(NOTFOUND, "key (%s) not found", key)
}

// Forbidden ...
func Forbidden(a, b string) pb.Response {
	return Error(FORBIDDEN, "%s can not read %s", a, b)
}

// NotAcceptable ...
func NotAcceptable(key string) pb.Response {
	return Error(NOTACCEPTABLE, "%s is not acceptable", key)
}

// Conflict ...
func Conflict(key string) pb.Response {
	return Error(CONFLICT, "%s has been existed", key)
}

// NotImplemented ...
func NotImplemented(cmd string) pb.Response {
	return Error(NOTIMPLEMENTED, "%s is not implemented", cmd)
}
