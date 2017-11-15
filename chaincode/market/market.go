package market

import (
	"fmt"
	"strconv"

	ccc "diviner/chaincode/common"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type marketCC struct{}

// NewMarketChaincode ...
func NewMarketChaincode() shim.Chaincode {
	return new(marketCC)
}

func (cc *marketCC) create(stub shim.ChaincodeStubInterface, user, event string, num float64, isFund bool) pb.Response {
	mb, err := ccc.Find(stub, user)
	if err != nil {
		return ccc.Errore(err)
	}

	mem, err := pbm.Unmarshal(mb)
	if err != nil {
		return ccc.Errorf("unmarshal member error: %v", err)
	}

	eb, err := ccc.Find(stub, event)
	if err != nil {
		return ccc.Errore(err)
	}

	evt, err := pbl.UnmarshalEvent(eb)
	if err != nil {
		return ccc.Errorf("unmarshal event error: %v", err)
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

	mb, err = pbm.Marshal(mem)
	if err != nil {
		return ccc.Errorf("marshal member error: %v", err)
	}

	if err := stub.PutState(mem.Id, mb); err != nil {
		return ccc.Errorf("update member balance error: %v", err)
	}

	bytes, err := pbl.MarshalMarket(market)
	if err != nil {
		return ccc.Errorf("marshal market error: %v", err)
	}

	if err := ccc.PutStateByCompositeKey(stub, pbl.MarketKey, bytes, market.Id); err != nil {
		return ccc.Errorf("put market error: %v", err)
	}
	return shim.Success(bytes)
}

func (cc *marketCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {

	bytes, err := ccc.FindByPartialCompositeKey(stub, pbl.MarketKey, id)
	if err != nil {
		return ccc.Errorf("query market (%s) error: %v", id, err)
	}

	_, err = pbl.UnmarshalMarket(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	return shim.Success(bytes)
}

func (cc *marketCC) settle(stub shim.ChaincodeStubInterface, id string) pb.Response {

	bytes, err := ccc.FindByPartialCompositeKey(stub, pbl.MarketKey, id)
	if err != nil {
		return ccc.Errorf("query market (%s) error: %v", id, err)
	}

	market, err := pbl.UnmarshalMarket(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	if market.Settled {
		return ccc.Errorf("can not settle a settled market")
	}

	market.Settled = true

	bytes2, err := pbl.MarshalMarket(market)
	if err != nil {
		return ccc.Errorf("marshal data error: %v", err)
	}

	if err := ccc.PutStateByCompositeKey(stub, pbl.MarketKey, bytes2, market.Id); err != nil {
		return ccc.Errorf("put market error: %v", err)
	}
	return shim.Success(nil)
}

// Init ...
func (cc *marketCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (cc *marketCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, params := stub.GetFunctionAndParameters()
	len := len(params)
	switch fcn {
	case "create":
		if len != 4 {
			return ccc.Errorf("args length error for create: %v", len)
		}

		flag := params[0]
		if flag != "fund" && flag != "liquidity" {
			return ccc.Errorf("flag must be `fund` or `liquidity`: %s", flag)
		}

		num, err := strconv.ParseFloat(params[3], 64)
		if err != nil {
			return ccc.Errorf("num must be float64: %s", params[3])
		}

		user := params[1]
		event := params[2]

		return cc.create(stub, user, event, num, flag == "fund")

	case "query":
		if len != 1 {
			return ccc.Errorf("args length error for query: %v", len)
		}
		return cc.query(stub, params[0])
	case "settle":
		if len != 1 {
			return ccc.Errorf("args length error for settle: %v", len)
		}

		return cc.settle(stub, params[0])
	}

	return ccc.Errorf("unknown function: %s", fcn)
}

func main() {
	err := shim.Start(NewMarketChaincode())

	if err != nil {
		fmt.Printf("creating member chaincode failed: %v", err)
	}
}
