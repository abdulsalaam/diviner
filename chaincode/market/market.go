package market

import (
	"fmt"

	ccc "diviner/chaincode/common"
	pbm "diviner/protos/lmsr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type marketCC struct{}

// NewMarketChaincode ...
func NewMarketChaincode() shim.Chaincode {
	return new(marketCC)
}

func (mkt *marketCC) create(user, event string, fund bool) pb.Response {
	return shim.Success(nil)
}

func (mkt *marketCC) query(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, existed, err := ccc.GetStateAndCheck(stub, id)
	if err != nil {
		return ccc.Errorf("query market (%s) error: %v", id, err)
	} else if !existed {
		return ccc.Errorf("market id (%s) is not existed", id)
	}

	_, err = pbm.UnmarshalMarket(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	return shim.Success(bytes)
}

func (mkt *marketCC) settled(stub shim.ChaincodeStubInterface, id string) pb.Response {
	bytes, existed, err := ccc.GetStateAndCheck(stub, id)
	if err != nil {
		return ccc.Errorf("query market (%s) error: %v", id, err)
	} else if !existed {
		return ccc.Errorf("market id (%s) is not existed", id)
	}

	market, err := pbm.UnmarshalMarket(bytes)
	if err != nil {
		return ccc.Errorf("unmarshal data error: %v", err)
	}

	if market.Settled {
		return ccc.Errorf("can not settle a settled market")
	}

	market.Settled = true

	bytes2, err := pbm.MarshalMarket(market)
	if err != nil {
		return ccc.Errorf("marshal data error: %v", err)
	}

	return ccc.PutStateAndReturn(stub, id, bytes2, nil)
}

// Init ...
func (mkt *marketCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (mkt *marketCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func main() {
	err := shim.Start(NewMarketChaincode())

	if err != nil {
		fmt.Printf("creating member chaincode failed: %v", err)
	}
}
