package oracle

import (
	"fmt"

	ccc "diviner/chaincode/common"

	pbl "diviner/protos/lmsr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type oracleCC struct{}

// NewOracleChaincode ...
func NewOracleChaincode() shim.Chaincode {
	return new(oracleCC)
}

func (cc *oracleCC) findMarkets(stub shim.ChaincodeStubInterface, evtId string) (*pbl.Markets, error) {
	lst, err := ccc.FindAllByPartialCompositeKey(stub, pbl.MarketKey, evtId)
	if err != nil {
		return nil, fmt.Errorf("find markets errors: %v", err)
	}

	var markets []*pbl.Market

	for k, v := range lst {
		if m, err := pbl.UnmarshalMarket(v); err != nil {
			return nil, fmt.Errorf("unmarshal market error at %s: %v", k, err)
		} else {
			markets = append(markets, m)
		}
	}

	return &pbl.Markets{
		List: markets,
	}, nil

}

func (cc *oracleCC) findAssets(stub shim.ChaincodeStubInterface, mktId string) (*pbl.Assets, error) {
	lst, err := ccc.FindAllByPartialCompositeKey(stub, pbl.AssetKey, mktId)
	if err != nil {
		return nil, fmt.Errorf("find assets errors: %v", err)
	}

	var assets []*pbl.Asset

	for k, v := range lst {
		if x, err := pbl.UnmarshalAsset(v); err != nil {
			return nil, fmt.Errorf("unmarshal asset error at %s: %v", k, err)
		} else {
			assets = append(assets, x)
		}
	}

	return &pbl.Assets{
		List: assets,
	}, nil

}

func (cc *oracleCC) markets(stub shim.ChaincodeStubInterface, evtId string) pb.Response {
	result, err := cc.findMarkets(stub, evtId)
	if err != nil {
		return ccc.Errore(err)
	}

	if bytes, err := pbl.MarshalMarkets(result); err != nil {
		return ccc.Errorf("marshal markets error: %v", err)
	} else {
		return shim.Success(bytes)
	}
}

func (cc *oracleCC) assets(stub shim.ChaincodeStubInterface, mktId string) pb.Response {
	result, err := cc.findAssets(stub, mktId)
	if err != nil {
		return ccc.Errore(err)
	}

	if bytes, err := pbl.MarshalAssets(result); err != nil {
		return ccc.Errorf("marshal assets error: %v", err)
	} else {
		return shim.Success(bytes)
	}
}

func (cc *oracleCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *oracleCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, args := stub.GetFunctionAndParameters()
	len := len(args)

	switch fcn {
	case "markets":
		if len != 1 {
			return ccc.Errorf("args length error for markets: %v", len)
		}
		return cc.markets(stub, args[0])
	case "assets":
		if len != 1 {
			return ccc.Errorf("args length error for assets: %v", len)
		}
		return cc.assets(stub, args[0])
	case "approve":
	case "settle":
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(NewOracleChaincode())

	if err != nil {
		fmt.Printf("creating oracle chaincode failed: %v", err)
	}
}
