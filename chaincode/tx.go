package main

import (
	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"diviner/common/cast"
	"fmt"
	"math"
	"strings"

	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type txCC struct{}

// NewLMSRChaincode ...
func NewTxChaincode() Mychaincode {
	return new(txCC)
}

func (cc *txCC) returnFloat(value float64) pb.Response {
	if bytes, err := cast.ToBytes(value); err != nil {
		return ccc.Errore(err)
	} else {
		return shim.Success(bytes)
	}
}

func (cc *txCC) markets(stub shim.ChaincodeStubInterface, evtId string) pb.Response {
	result, err := ccu.FindAllMarkets(stub, evtId)
	if err != nil {
		return ccc.Errore(err)
	}

	return ccc.MarshalAndReturn(result)
}

func (cc *txCC) assets(stub shim.ChaincodeStubInterface, keys []string) pb.Response {
	result, err := ccu.FindAllAssets(stub, keys...)
	if err != nil {
		return ccc.Errorf("find all assets of (%s) error: %v", strings.Join(keys, pbl.Sep), err)
	}

	return ccc.MarshalAndReturn(result)
}

func (cc *txCC) updateAsset(stub shim.ChaincodeStubInterface, user, share string, volume float64) (*pbl.Asset, error) {

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

func (cc *txCC) tx(stub shim.ChaincodeStubInterface, user, share string, volume float64) pb.Response {
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

	if member.Assets == nil {
		member.Assets = make(map[string]float64)
	}
	member.Assets[asset.Id] = asset.Volume

	if _, err := ccc.PutMessage(stub, member.Id, member); err != nil {
		return ccc.Errorf("put member error: %v", err)
	}

	if _, err := ccu.PutMarket(stub, market); err != nil {
		return ccc.Errorf("put market error: %v", err)
	}

	return cc.returnFloat(math.Abs(price))

}

func (cc *txCC) settleAssets(stub shim.ChaincodeStubInterface, marketId, result string) error {

	evt, mkt, ok := pbl.SepMarketID(marketId)
	if !ok {
		return fmt.Errorf("id (%s) format error", marketId)
	}

	assets, err := ccu.FindAllAssets(stub, evt, mkt)

	if err != nil {
		return fmt.Errorf("find assets of (%s) error: %v", marketId, err)
	}

	updateMembers := make(map[string]*pbm.Member)

	for _, x := range assets.List {
		aevt, amkt, aoutcome, auser, ok := pbl.SepAssetID(x.Id)
		if aevt != evt || amkt != mkt {
			return fmt.Errorf("data error, need (%s, %s) but (%s, %s)", evt, mkt, aevt, amkt)
		}
		if !ok {
			return fmt.Errorf("asset id (%s) error", x.Id)
		}

		mem, existed, err := ccu.GetMemberAndCheck(stub, auser)
		if err != nil {
			return fmt.Errorf("get member (%s) error: %v", auser, err)
		} else if !existed {
			return fmt.Errorf("member (%s) not found", auser)
		}

		if aoutcome == result {
			mem.Balance += mem.Assets[x.Id]
		}

		delete(mem.Assets, x.Id)
		updateMembers[mem.Id] = mem
		if err := stub.DelState(x.Id); err != nil {
			return fmt.Errorf("delete asset (%s) error: %v", x.Id, err)
		}
	}

	for _, x := range updateMembers {
		if _, err := ccc.PutMessage(stub, x.Id, x); err != nil {
			return fmt.Errorf("update member (%s) error: %v", x.Id, err)
		}
	}

	return nil
}

func (cc *txCC) setteMarket(stub shim.ChaincodeStubInterface, market *pbl.Market, result string) (float64, error) {
	if market.Settled {
		return 0.0, fmt.Errorf("market (%s) is settled", market.Id)
	}

	owner, existed, err := ccu.GetMemberAndCheck(stub, market.User)
	if err != nil {
		return 0.0, fmt.Errorf("find market owner (%s) error: %v", market.User, err)
	} else if !existed {
		return 0.0, fmt.Errorf("market owner (%s) not found", market.User)
	}

	share := pbl.ShareID(market.Id, result)
	volume, ok := market.Shares[share]
	if !ok {
		return 0.0, fmt.Errorf("share (%s) not found in market (%s)", share, market.Id)
	}

	returns := market.Cost - volume

	owner.Balance += returns

	market.Settled = true

	if _, err := ccu.PutMarket(stub, market); err != nil {
		return 0.0, fmt.Errorf("put market (%s) error: %v", market.Id, err)
	}

	if _, err := ccc.PutMessage(stub, owner.Id, owner); err != nil {
		return 0.0, fmt.Errorf("put member (%s) error: %v", owner.Id, err)
	}

	if err := cc.settleAssets(stub, market.Id, result); err != nil {
		return 0.0, err
	}

	return returns, nil
}

func (cc *txCC) settle(stub shim.ChaincodeStubInterface, id, result string) pb.Response {
	market, existed, err := ccu.GetMarketAndCheck(stub, id)
	if err != nil {
		return ccc.Errorf("get market (%s) error: %v", id, err)
	} else if !existed {
		return ccc.Errorf("market (%s) not found", id)
	}

	if a, err := cc.setteMarket(stub, market, result); err != nil {
		return ccc.Errore(err)
	} else {
		return cc.returnFloat(a)
	}
}

func (cc txCC) approve(stub shim.ChaincodeStubInterface, id, result string) pb.Response {
	event, err := ccu.FindEvent(stub, id)
	if err != nil {
		return ccc.Errorf("find event (%s) error: %v", id, err)
	}

	if event.Approved {
		return ccc.Errorf("event (%s) is approved", id)
	}

	if pbl.FindOutcome(event, result) < 0 {
		return ccc.Errorf("result (%s) is not in event (%s)", result, id)
	}

	event.Approved = true

	if _, err := ccc.PutMessage(stub, event.Id, event); err != nil {
		return ccc.Errorf("put event (%s) error: %v", event.Id, err)
	}

	markets, err := ccu.FindAllMarkets(stub, event.Id)
	if err != nil {
		ccc.Errorf("find martets of event (%s) error: %v", event.Id, err)
	}

	for _, x := range markets.List {
		if x.Settled {
			continue
		} else if _, err := cc.setteMarket(stub, x, result); err != nil {
			return ccc.Errore(err)
		}
	}

	return shim.Success(nil)

}

// Invoke ...
func (cc *txCC) Invoke(stub shim.ChaincodeStubInterface, fcn string, args [][]byte) pb.Response {
	len := len(args)

	switch fcn {
	case "buy", "sell":
		if len != 3 {
			return ccc.Errorf("args length error for buying: %v", len)
		}

		volume, err := cast.BytesToFloat64(args[2])
		if err != nil {
			return ccc.Errorf("volume must be float64: %s", args[2])
		} else if volume <= 0 {
			return ccc.Errorf("volume must be larger than 0: %v", volume)
		}
		user := string(args[0])
		share := string(args[1])

		if fcn == "sell" {
			volume = -volume
		}

		return ccc.SetEventAndReturn(stub, "tx", cc.tx(stub, user, share, volume))

	case "markets":
		if len != 1 {
			return ccc.Errorf("args length error for markets: %v", len)
		}
		return cc.markets(stub, string(args[0]))
	case "assets":
		if len < 1 {
			return ccc.Errorf("args length error for assets: %v", len)
		}

		ids := cast.ByteArrayToStrings(args)
		return cc.assets(stub, ids)

	case "settle":
		if len != 2 {
			return ccc.Errorf("args length error for settle: %v", len)
		}

		return cc.settle(stub, string(args[0]), string(args[1]))

	case "approve":
		if len != 2 {
			return ccc.Errorf("args length error for approve: %v", len)
		}

		return cc.approve(stub, string(args[0]), string(args[1]))
	}

	return ccc.Errorf("unknown function: %s", fcn)
}
