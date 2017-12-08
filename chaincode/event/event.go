package event

import (
	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"diviner/common/base58"

	"github.com/golang/protobuf/ptypes"

	pbmk "diviner/protos/market"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("event")

type eventCC struct {
	expired int64
}

// NewEventChaincode ...
func NewEventChaincode() shim.Chaincode {
	return &eventCC{
		expired: 300,
	}
}

func (cc *eventCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, existed, err := ccc.GetStateAndCheck(stub, id)
	if err != nil {
		return ccc.Errore(err)
	} else if !existed {
		return ccc.NotFound(id)
	}

	if _, err := pbmk.UnmarshalEvent(bytes); err != nil {
		return ccc.Errore(err)
	}

	return shim.Success(bytes)
}

func (cc *eventCC) create(stub shim.ChaincodeStubInterface, event *pbmk.Event) pb.Response {
	if event.Id == "" || event.Title == "" || event.User == "" {
		return ccc.BadRequest("data error about id (%s), title (%s), user (%s)", event.Id, event.Title, event.User)
	}

	_, existed, err := ccc.GetStateAndCheck(stub, event.Id)
	if err != nil {
		return ccc.Errore(err)
	} else if existed {
		return ccc.Conflict(event.Id)
	}

	// check event and put
	len := len(event.Outcomes)
	if len < 2 {
		return ccc.BadRequest("size of outcomes must be more than or equal 2, but %d", len)
	}

	curr := ptypes.TimestampNow()

	if event.End.Seconds <= curr.Seconds {
		return ccc.BadRequest("event expired: %v", event.End)
	}

	resp := stub.InvokeChaincode("member", [][]byte{
		[]byte("find"),
		[]byte(event.User),
	}, "")

	if ccc.NotOK(&resp) {
		return resp
	}

	event.Allowed = true
	event.Approved = false
	event.Result = ""

	return ccc.PutMessageAndReturn(stub, event.Id, event)

}

// Init ...
func (cc *eventCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("init event chaincode")
	return shim.Success(nil)
}

// Invoke ...
func (cc *eventCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	len := len(args)

	if len < 2 {
		return ccc.BadRequest("args length must be more than or equal 2, but %d", len)
	}

	cmd := string(args[0])

	switch cmd {
	case "query":
		return cc.query(stub, string(args[1]))

	case "create":
		if len != 3 {
			return ccc.BadRequest("args length must be 3, but %d", len)
		}

		v, ok, err := ccu.CheckAndPutVerfication(stub, args[1], args[2], cc.expired)
		if err != nil {
			return ccc.Errore(err)
		} else if !ok {
			return ccc.Unauthorized()
		}

		event, err := pbmk.UnmarshalEvent(args[1])
		if err != nil {
			return ccc.Errore(err)
		}

		me := base58.Encode(v.PublicKey)
		if event.User != me {
			return ccc.BadRequest("event creator error, must be %s, but %s", me, event.User)
		}

		return cc.create(stub, event)

	}
	return ccc.NotImplemented(cmd)
}

/*
func main() {
	logger.SetLevel(shim.LogDebug)
	logger.Debug("start event chaincode")
	err := shim.Start(NewEventChaincode())
	if err != nil {
		logger.Errorf("creating event chaincode failed: %v\n", err)
	}
}
*/
