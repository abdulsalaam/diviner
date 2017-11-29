package main

import (
	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"diviner/common/cast"
	pbl "diviner/protos/lmsr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type marketCC struct{}

// NewMarketChaincode ...
func NewMarketChaincode() Mychaincode {
	return new(marketCC)
}

func (cc *marketCC) create(stub shim.ChaincodeStubInterface, user, event string, num float64, isFund bool) pb.Response {
	mem, existed, err := ccu.GetMemberAndCheck(stub, user)
	if err != nil {
		return ccc.Errore(err)
	} else if !existed {
		return ccc.Errorf("market creator not found: %s", user)
	}

	evt, err := ccu.FindEvent(stub, event)
	if err != nil {
		return ccc.Errore(err)
	}

	var fund float64

	if isFund {
		fund = num
	} else {
		fund = pbl.Fund(num, len(evt.Outcomes))
	}

	if mem.Balance < fund {
		return ccc.Errorf("usre balance is not enough, need %v but %v", fund, mem.Balance)
	}

	var market *pbl.Market

	if isFund {
		market, err = pbl.NewMarketWithFund(mem.Id, evt, num)
	} else {
		market, err = pbl.NewMarketWithLiquidity(mem.Id, evt, num)
	}

	if err != nil {
		return ccc.Errorf("new market error: %v", err)
	}

	mem.Balance -= market.Fund

	if _, err := ccc.PutMessage(stub, mem.Id, mem); err != nil {
		return ccc.Errorf("put member error: %v", err)
	}

	if bytes, err := ccu.PutMarket(stub, market); err != nil {
		return ccc.Errorf("put market error: %v", err)
	} else {
		return shim.Success(bytes)
	}
}

func (cc *marketCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {
	if m, existed, err := ccu.GetMarketAndCheck(stub, id); err != nil {
		return ccc.Errorf("query market (%s) error: %v", id, err)
	} else if !existed {
		return ccc.Errorf("market (%s) not found", id)
	} else {
		return ccc.MarshalAndReturn(m)
	}
}

func (cc *marketCC) settle(stub shim.ChaincodeStubInterface, id string) pb.Response {
	market, existed, err := ccu.GetMarketAndCheck(stub, id)
	if err != nil {
		return ccc.Errorf("find market (%s) error: %v", id, err)
	} else if !existed {
		return ccc.Errorf("market (%s) not found", id)
	}

	if market.Settled {
		return ccc.Errorf("can not settle a settled market")
	}

	market.Settled = true

	if _, err := ccu.PutMarket(stub, market); err != nil {
		return ccc.Errorf("put market error: %v", err)
	}

	return shim.Success(nil)
}

// Invoke ...
func (cc *marketCC) Invoke(stub shim.ChaincodeStubInterface, fcn string, args [][]byte) pb.Response {
	len := len(args)
	switch fcn {
	case "create":
		if len != 4 {
			return ccc.Errorf("args length for market create error: %v", len)
		}

		flag, err := cast.BytesToBool(args[3])
		if err != nil {
			return ccc.Errorf("isFund must be boolean: %v, %v", args[3], err)
		}

		num, err := cast.BytesToFloat64(args[2])
		if err != nil {
			return ccc.Errorf("num must be float64: %v, %v", args[2], err)
		}

		user := string(args[0])
		event := string(args[1])

		return ccc.SetEventAndReturn(stub, "market", cc.create(stub, user, event, num, flag))

	case "query":
		if len != 1 {
			return ccc.Errorf("args length error for query: %v", len)
		}

		return cc.query(stub, string(args[0]))

	case "settle":
		if len != 1 {
			return ccc.Errorf("args length error for settle: %v", len)
		}

		return cc.settle(stub, string(args[0]))
	}

	return ccc.Errorf("unknown function: %s", fcn)
}
