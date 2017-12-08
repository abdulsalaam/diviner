package util

import (
	ccc "diviner/chaincode/common"
	pbc "diviner/protos/common"
	pbl "diviner/protos/market"
	pbm "diviner/protos/member"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetMemberAndCheck ...
func GetMemberAndCheck(stub shim.ChaincodeStubInterface, id string) (*pbm.Member, bool, error) {
	bytes, existed, err := ccc.GetStateAndCheck(stub, id)
	if err != nil {
		return nil, false, err
	} else if !existed {
		return nil, false, nil
	}

	x, err := pbm.Unmarshal(bytes)
	if err != nil {
		return nil, false, err
	}

	return x, true, nil
}

// FindEvent ...
func FindEvent(stub shim.ChaincodeStubInterface, id string) (*pbl.Event, error) {
	if bytes, err := ccc.Find(stub, id); err != nil {
		return nil, err
	} else if evt, err := pbl.UnmarshalEvent(bytes); err != nil {
		return nil, err
	} else {
		return evt, nil
	}
}

// PutMarket ...
func PutMarket(stub shim.ChaincodeStubInterface, market *pbl.Market) ([]byte, error) {
	evtID, mktID, ok := pbl.SepMarketID(market.Id)
	if !ok {
		return nil, fmt.Errorf("market id format error: %s", market.Id)
	}

	return ccc.PutMessageWithCompositeKey(stub, market, pbl.MarketKey, evtID, mktID)
}

// GetMarketAndCheck ...
func GetMarketAndCheck(stub shim.ChaincodeStubInterface, market string) (*pbl.Market, bool, error) {
	eid, mid, ok := pbl.SepMarketID(market)
	if !ok {
		return nil, false, fmt.Errorf("market id format error")
	}

	result, existed, err := ccc.GetStateByCompositeKeyAndCheck(stub, pbl.MarketKey, eid, mid)
	if err != nil {
		return nil, false, fmt.Errorf("query market (%s) error: %v", market, err)
	} else if !existed {
		return nil, false, nil
	}

	bytes := ccc.GetOneValue(result)
	if err != nil {
		return nil, false, fmt.Errorf("list content error")
	}

	tmp, err := pbl.UnmarshalMarket(bytes)
	if err != nil {
		return nil, false, err
	}
	return tmp, true, nil
}

// FindAllMarkets ...
func FindAllMarkets(stub shim.ChaincodeStubInterface, event string) (*pbl.Markets, error) {
	result, err := ccc.FindAllByPartialCompositeKey(stub, pbl.MarketKey, event)
	if err != nil {
		return nil, fmt.Errorf("find all markets errors: %v", err)
	}

	var markets []*pbl.Market

	for k, v := range result {
		m, err := pbl.UnmarshalMarket(v)
		if err != nil {
			return nil, fmt.Errorf("unmarshal market error at %s: %v", k, err)
		}
		markets = append(markets, m)
	}

	return &pbl.Markets{
		List: markets,
	}, nil
}

// PutAsset ...
func PutAsset(stub shim.ChaincodeStubInterface, asset *pbl.Asset) ([]byte, error) {
	event, market, outcome, member, ok := pbl.SepAssetID(asset.Id)
	if !ok {
		return nil, fmt.Errorf("asset id format error: %s", asset.Id)
	}

	return ccc.PutMessageWithCompositeKey(stub, asset, pbl.AssetKey, event, market, outcome, member)
}

// GetAssetAndCheck ...
func GetAssetAndCheck(stub shim.ChaincodeStubInterface, asset string) (*pbl.Asset, bool, error) {
	event, market, outcome, member, ok := pbl.SepAssetID(asset)
	if !ok {
		return nil, false, fmt.Errorf("asset id format error: %s", asset)
	}

	result, existed, err := ccc.GetStateByCompositeKeyAndCheck(stub, pbl.AssetKey, event, market, outcome, member)
	if err != nil {
		return nil, false, fmt.Errorf("find asset (%s) error: %v", asset, err)
	} else if !existed {
		return nil, false, nil
	}

	bytes := ccc.GetOneValue(result)
	if bytes == nil {
		return nil, false, fmt.Errorf("list content error")
	}

	tmp, err := pbl.UnmarshalAsset(bytes)
	if err != nil {
		return nil, false, err
	}

	return tmp, true, nil
}

// FindAllAssets ...
func FindAllAssets(stub shim.ChaincodeStubInterface, keys ...string) (*pbl.Assets, error) {
	result, err := ccc.FindAllByPartialCompositeKey(stub, pbl.AssetKey, keys...)

	if err != nil {
		return nil, fmt.Errorf("find all assets errors: %v", err)
	}

	var assets []*pbl.Asset

	for k, v := range result {
		x, err := pbl.UnmarshalAsset(v)
		if err != nil {
			return nil, fmt.Errorf("unmarshal asset error at %s: %v", k, err)
		}
		assets = append(assets, x)
	}

	return &pbl.Assets{
		List: assets,
	}, nil
}

// CheckAndPutVerfication ...
func CheckAndPutVerfication(stub shim.ChaincodeStubInterface, in, check []byte, expired int64) (*pbc.Verification, bool, error) {
	v, err := pbc.Unmarshal(check)
	if err != nil {
		return nil, false, err
	}

	if ok, err := pbc.Verify(v, in, expired); err != nil {
		return nil, false, err
	} else if !ok {
		return nil, false, nil
	}

	if err := stub.PutState("chk"+stub.GetTxID(), check); err != nil {
		return v, false, err
	}

	return v, true, nil
}
