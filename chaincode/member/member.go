package main

import (
	"diviner/common/base58"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	pbc "diviner/protos/common"
	pbm "diviner/protos/member"

	pb "github.com/hyperledger/fabric/protos/peer"
)

type memberCC struct {
	expired int64
	fee     float64
	balance float64
}

// NewMemberChaincode ...
func NewMemberChaincode() shim.Chaincode {
	return &memberCC{
		expired: 300,
		fee:     0.001,
		balance: 10000.0,
	}
}

var logger = shim.NewLogger("member")

func (cc *memberCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("init member chaincode")
	return shim.Success(nil)
}

func (cc *memberCC) amount(in float64) float64 {
	return in * (1.0 + cc.fee)
}

func (cc *memberCC) register(stub shim.ChaincodeStubInterface, member *pbm.Member) pb.Response {
	_, existed, err := ccc.GetStateAndCheck(stub, member.Address)
	if err != nil {
		return ccc.Errorf("get member address (%s) error: %v", member.Address, err)
	} else if existed {
		return ccc.NotFound(member.Address)
	}

	member = pbm.NewMember(member.Address, cc.balance)

	return ccc.PutMessageAndReturn(stub, member.Address, member)
}

func (cc *memberCC) transfer(stub shim.ChaincodeStubInterface, tx *pbm.Transfer, from string) pb.Response {
	target, existed, err := ccu.GetMemberAndCheck(stub, tx.Target)
	if err != nil {
		return ccc.Errorf("find target (%s) error: %v", tx.Target, err)
	} else if !existed {
		return ccc.NotFound(tx.Target)
	}

	if target.Blocked {
		return ccc.NotAcceptable(tx.Target)
	}

	source, existed, err := ccu.GetMemberAndCheck(stub, from)
	if err != nil {
		return ccc.Errorf("find source (%s) error: %v", from, err)
	} else if !existed {
		return ccc.NotFound(from)
	}

	if source.Blocked {
		return ccc.NotAcceptable(from)
	}

	amount := cc.amount(tx.Amount)
	if source.Balance < amount {
		return ccc.Errorf("balance is not enough. need (%v) but (%v)", tx.Amount+cc.fee, source.Balance)
	}

	source.Balance -= amount
	target.Balance += tx.Amount

	if _, err := ccc.PutMessage(stub, target.Address, target); err != nil {
		return ccc.Errorf("update target error: %v", err)
	}

	return ccc.PutMessageAndReturn(stub, source.Address, source)
}

// Invoke ...
func (cc *memberCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	len := len(args)

	if len != 3 {
		return ccc.BadRequest("args len error: %d", len)
	}

	cmd := string(args[0])
	if err := ccu.CheckAndPutVerfication(stub, args[1], args[2], cc.expired); err != nil {
		return ccc.Errore(err)
	}

	switch cmd {
	case "query":
		id := string(args[1])
		v, _ := pbc.Unmarshal(args[2])
		me := base58.Encode(v.PublicKey)

		if id != me {
			return ccc.Forbidden(me, id)
		}

		bytes, existed, err := ccc.GetStateAndCheck(stub, id)
		if err != nil {
			return ccc.Errore(err)
		} else if !existed {
			return ccc.NotFound(id)
		}
		return shim.Success(bytes)

	case "register":
		mem, err := pbm.Unmarshal(args[1])
		if err != nil {
			ccc.Errore(err)
		}
		return cc.register(stub, mem)

	case "transfer":
		tx, err := pbm.UnmarshalTransfer(args[1])
		if err != nil {
			return ccc.Errore(err)
		}
		v, _ := pbc.Unmarshal(args[2])
		from := base58.Encode(v.PublicKey)

		return cc.transfer(stub, tx, from)
	}

	return ccc.NotImplemented(cmd)
}

func main() {
	logger.SetLevel(shim.LogDebug)
	logger.Debug("start member chaincode")
	err := shim.Start(NewMemberChaincode())
	if err != nil {
		logger.Errorf("creating member chaincode failed: %v\n", err)
	}
}
