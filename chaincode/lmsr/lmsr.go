package lmsr

import (
	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"fmt"
	"math"
	"strconv"
	"strings"

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

func (cc *lmsrCC) markets(stub shim.ChaincodeStubInterface, evtId string) pb.Response {
	result, err := ccu.FindAllMarkets(stub, evtId)
	if err != nil {
		return ccc.Errore(err)
	}

	return ccc.MarshalAndReturn(result)
}

func (cc *lmsrCC) assets(stub shim.ChaincodeStubInterface, keys []string) pb.Response {
	result, err := ccu.FindAllAssets(stub, keys...)
	if err != nil {
		return ccc.Errorf("find all assets of (%s) error: %v", strings.Join(keys, pbl.Sep), err)
	}

	return ccc.MarshalAndReturn(result)
}

func (cc *lmsrCC) updateAsset(stub shim.ChaincodeStubInterface, user, share string, volume float64) (*pbl.Asset, error) {
	id := pbl.AssetID(user, share)
	asset, existed, err := ccu.GetAssetAndCheck(stub, id)
	if err != nil {
		return nil, err
	}

	if !existed {
		asset = &pbl.Asset{
			Id:     id,
			Volume: 0.0,
		}
	}

	if asset.Volume+volume < 0 {
		return nil, fmt.Errorf("asset volume is not enough, need %v but %v", math.Abs(volume), asset.Volume)
	}

	asset.Volume += volume

	if _, err := ccu.PutAsset(stub, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (cc *lmsrCC) tx(stub shim.ChaincodeStubInterface, user, share string, volume float64) pb.Response {
	mktId, _, ok := pbl.SepShareID(share)
	if !ok {
		return ccc.Errorf("share id format error")
	}

	market, existed, err := ccu.GetMarketAndCheck(stub, mktId)
	if err != nil {
		return ccc.Errore(err)
	} else if !existed {
		return ccc.Errorf("market (%s) not found", mktId)
	}

	_, ok = market.Shares[share]
	if !ok {
		return ccc.Errorf("share (%s) not found in market (%s)", share, mktId)
	}

	tmp, err := ccc.Find(stub, user)
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

	if _, err := ccc.PutMessage(stub, member.Id, member); err != nil {
		return ccc.Errorf("put member error: %v", err)
	}

	if _, err := ccu.PutMarket(stub, market); err != nil {
		return ccc.Errorf("put market error: %v", err)
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

	switch fcn {
	case "buy", "sell":
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

		if fcn == "sell" {
			volume = -volume
		}

		return cc.tx(stub, user, share, volume)

	case "markets":
		if len != 1 {
			return ccc.Errorf("args length error for markets: %v", len)
		}
		return cc.markets(stub, args[0])
	case "assets":
		if len < 1 {
			return ccc.Errorf("args length error for assets: %v", len)
		}

		return cc.assets(stub, args)
	}

	return ccc.Errorf("unknown function: %s", fcn)
}

func main() {
	err := shim.Start(NewLMSRChaincode())

	if err != nil {
		fmt.Printf("creating lmsr chaincode failed: %v", err)
	}
}
