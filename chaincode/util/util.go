package util

import (
	ccc "diviner/chaincode/common"
	pbl "diviner/protos/lmsr"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// PutMarket ...
func PutMarket(stub shim.ChaincodeStubInterface, market *pbl.Market) ([]byte, error) {
	evtId, mktId, ok := pbl.SepMarketID(market.Id)
	if !ok {
		return nil, fmt.Errorf("market id format error: %s", market.Id)
	}

	return ccc.PutMessageWithCompositeKey(stub, market, pbl.MarketKey, evtId, mktId)
}

func PutAsset(stub shim.ChaincodeStubInterface, asset *pbl.Asset) ([]byte, error) {
	event, market, outcome, member, ok := pbl.SepAssetID(asset.Id)
	if !ok {
		return nil, fmt.Errorf("asset id format error: %s", asset.Id)
	}

	return ccc.PutMessageWithCompositeKey(stub, asset, pbl.AssetKey, event, market, outcome, member)
}

func FindMarket(stub shim.ChaincodeStubInterface, market string) (*pbl.Market, error) {
	eid, mid, ok := pbl.SepMarketID(market)
	if !ok {
		return nil, fmt.Errorf("market id format error")
	}

	bytes, err := ccc.FindByPartialCompositeKey(stub, pbl.MarketKey, eid, mid)
	if err != nil {
		return nil, fmt.Errorf("query market (%s) error: %v", market, err)
	}

	return pbl.UnmarshalMarket(bytes)
}

func FindAllMarkets(stub shim.ChaincodeStubInterface, event string) (*pbl.Markets, error) {
	lst, err := ccc.FindAllByPartialCompositeKey(stub, pbl.MarketKey, event)
	if err != nil {
		return nil, fmt.Errorf("find all markets errors: %v", err)
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
