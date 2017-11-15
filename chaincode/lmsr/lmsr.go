package lmsr

import (
	ccc "diviner/chaincode/common"
	"fmt"
	"math"
	"strconv"

	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type lmsrCC struct{}

// NewLMSRChaincode ...
func NewLMSRChaincode() shim.Chaincode {
	return new(lmsrCC)
}

func (cc *lmsrCC) updateAsset(stub shim.ChaincodeStubInterface, user, share string, volume float64) (*pbl.Asset, error) {
	id := pbl.AssetID(user, share)
	var asset *pbl.Asset

	bytes, err := ccc.FindByPartialCompositeKey(stub, pbl.AssetKey, id)
	if err != nil {
		if asset, err = pbl.NewAsset(user, share, 0.0); err != nil {
			return nil, err
		}
	} else {
		if asset, err = pbl.UnmarshalAsset(bytes); err != nil {
			return nil, err
		}
	}

	if asset.Volume+volume < 0 {
		return nil, fmt.Errorf("asset volume is not enough, need %v but %v", math.Abs(volume), asset.Volume)
	}

	asset.Volume += volume

	if assetKey, err := stub.CreateCompositeKey(pbl.AssetKey, []string{asset.Id}); err != nil {
		return nil, err
	} else if err2 := ccc.PutMessage(stub, assetKey, asset); err2 != nil {
		return nil, err2
	}

	return asset, nil
}

func (cc *lmsrCC) tx(stub shim.ChaincodeStubInterface, user, share string, volume float64) pb.Response {
	mktId, _, ok := pbl.SepShareID(share)
	if !ok {
		return ccc.Errorf("share id format error")
	}

	tmp, err := ccc.FindByPartialCompositeKey(stub, pbl.MarketKey, mktId)
	if err != nil {
		return ccc.Errore(err)
	}

	market, err := pbl.UnmarshalMarket(tmp)
	if err != nil {
		return ccc.Errorf("unmarshal market error: %v", err)
	}

	_, ok = market.Shares[share]
	if !ok {
		return ccc.Errorf("share (%s) not found in market (%s)", share, mktId)
	}

	tmp, err = ccc.Find(stub, user)
	if err != nil {
		return ccc.Errore(err)
	}

	member, err := pbm.Unmarshal(tmp)
	if err != nil {
		return ccc.Errorf("unmarshal member error: %v", err)
	}

	price, err := pbl.EstimateMarket(market, share, volume)
	if err != nil {
		return ccc.Errore(err)
	}

	if member.Balance-price < 0 {
		return ccc.Errorf("member balance is not enough, need %v but %v", math.Abs(price), member.Balance)
	}

	pbl.UpdateMarket(market, share, volume)
	member.Balance -= price
	asset, err := cc.updateAsset(stub, user, share, volume)
	if err != nil {
		return ccc.Errorf("put asset error: %v", err)
	}
	member.Assets[asset.Id] = asset.Volume

	if err := ccc.PutMessage(stub, member.Id, member); err != nil {
		return ccc.Errorf("put member error: %v", err)
	}

	if mktbytes, err := pbl.MarshalMarket(market); err != nil {
		return ccc.Errorf("marshal market error")
	} else if err := ccc.PutStateByCompositeKey(stub, pbl.MarketKey, mktbytes, market.Id); err != nil {
		return ccc.Errorf("put market error")
	}

	return shim.Success(nil)
}

func (cc *lmsrCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (cc *lmsrCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, args := stub.GetFunctionAndParameters()
	len := len(args)

	if len != 3 {
		return ccc.Errorf("args length error for buying: %v", len)
	}

	volume, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return ccc.Errorf("volume must be float64: %s", args[2])
	} else if volume <= 0 {
		return ccc.Errorf("volume must be larger than 0: %v", volume)
	}
	user := args[0]
	share := args[1]

	switch fcn {
	case "buy":
		return cc.tx(stub, user, share, volume)
	case "sell":
		return cc.tx(stub, user, share, -volume)
	}

	return ccc.Errorf("unknown function: %s", fcn)
}

func main() {
	err := shim.Start(NewLMSRChaincode())

	if err != nil {
		fmt.Printf("creating lmsr chaincode failed: %v", err)
	}
}
