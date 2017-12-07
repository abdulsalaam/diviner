package util

import (
	"fmt"
	"os"
	"testing"

	ccc "diviner/chaincode/common"
	"diviner/common/csp"
	pbl "diviner/protos/market"
	pbm "diviner/protos/member"

	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type testcc struct{}

func (cc *testcc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *testcc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

var (
	stub     = ccc.NewMockStub("test", new(testcc))
	title    = "gogogo"
	outcomes = []string{"yes", "no"}
	balance  = 10000.0
	member   *pbm.Member
	event    *pbl.Event
	market1  *pbl.Market
	market2  *pbl.Market
	asset1   *pbl.Asset
	asset2   *pbl.Asset
)

func getOneKey(x map[string]float64) (string, bool) {

	for k := range x {
		return k, true
	}

	return "", false
}

func TestMain(m *testing.M) {
	priv, _ := csp.GeneratePrivateTempKey()
	member, _ = pbm.NewMember(priv, balance)
	event, _ = pbl.NewEvent(member.Id, title, outcomes[0], outcomes[1])
	market1, _ = pbl.NewMarketWithFund(member.Id, event, 100.0)
	market2, _ = pbl.NewMarketWithFund(member.Id, event, 200.0)

	s1, _ := getOneKey(market1.Shares)
	asset1, _ = pbl.NewAsset(member.Id, s1, 10)
	s2, _ := getOneKey(market2.Shares)
	asset2, _ = pbl.NewAsset(member.Id, s2, 20)

	txid := uuid.New().String()
	stub.MockTransactionStart(txid)

	ccc.PutMessage(stub, member.Id, member)

	ccc.PutMessage(stub, event.Id, event)

	eid, mid, _ := pbl.SepMarketID(market1.Id)
	ccc.PutMessageWithCompositeKey(stub, market1, pbl.MarketKey, eid, mid)

	eid, mid, _ = pbl.SepMarketID(market2.Id)
	ccc.PutMessageWithCompositeKey(stub, market2, pbl.MarketKey, eid, mid)

	if _, err := PutAsset(stub, asset1); err != nil {
		fmt.Printf("put assert error: %v\n", asset1)
		os.Exit(-1)
	}

	if _, err := PutAsset(stub, asset2); err != nil {
		fmt.Printf("put asset error: %v\n", asset2)
	}

	stub.MockTransactionEnd(txid)

	m.Run()
}

func TestFindMember(t *testing.T) {
	m, existed, err := GetMemberAndCheck(stub, member.Id)
	if err != nil {
		t.Fatalf("get member (%s) failed: %v", member.Id, err)
	}

	if !existed {
		t.Fatalf("must find an existed member (%s)", member.Id)
	}

	if m.Id != member.Id || m.Balance != member.Balance || m.Address != member.Address {
		t.Fatalf("data not match: %v, %v", m, member)
	}

	m, existed, err = GetMemberAndCheck(stub, "abcdef")
	if err != nil {
		t.Fatalf("get non-existed member failed: %v", err)
	}

	if existed {
		t.Fatal("can not find non-existed member")
	}

	if m != nil {
		t.Fatal("non-existed member must return nil")
	}

}

func TestFindMaket(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	m, existed, err := GetMarketAndCheck(stub, market1.Id)
	if err != nil {
		t.Fatalf("find market failed: %v", err)
	} else if !existed {
		t.Fatal("data not found")
	}

	if !pbl.CmpMarket(m, market1) {
		t.Fatalf("market data not match: %v, %v", m, market1)
	}
}

func TestPutAndFindMarket(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	event1, _ := pbl.NewEvent(member.Id, "test", "a", "b")
	m1, _ := pbl.NewMarketWithFund(member.Id, event1, 10.0)

	bytes, err := PutMarket(stub, m1)
	if err != nil {
		t.Fatalf("put market failed: %v\n", err)
	}

	m2, _ := pbl.UnmarshalMarket(bytes)

	m3, _, _ := GetMarketAndCheck(stub, m1.Id)

	if !pbl.CmpMarket(m1, m2) || !pbl.CmpMarket(m1, m3) {
		t.Fatalf("data not match: \n%v\n%v\n%v\n", m1, m2, m3)
	}
}

func TestFindAllMarket(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	lst, err := FindAllMarkets(stub, event.Id)
	if err != nil {
		t.Fatalf("find all markets of event (%s) failed: %v", event.Id, err)
	}

	if len(lst.List) != 2 {
		t.Fatal("length of list failed: %d", len(lst.List))
	}

	if !(pbl.CmpMarket(market1, lst.List[0]) ||
		pbl.CmpMarket(market1, lst.List[1]) ||
		pbl.CmpMarket(market2, lst.List[0]) ||
		pbl.CmpMarket(market2, lst.List[1])) {
		t.Fatal("data not match")
	}

}

func TestGetAndFindAsset(t *testing.T) {
	txid := uuid.New().String()
	stub.MockTransactionStart(txid)
	defer stub.MockTransactionEnd(txid)

	a1, existed, err := GetAssetAndCheck(stub, asset1.Id)
	if err != nil {
		t.Fatalf("get an existed asset failed: %v", err)
	}

	if !existed {
		t.Fatalf("can not find an existed asset: %v", asset1.Id)
	}

	if *a1 != *asset1 {
		t.Fatal("data not match")
	}

	evt, mkt1, _ := pbl.SepMarketID(market1.Id)
	ma1, err := FindAllAssets(stub, evt, mkt1)
	if err != nil {
		t.Fatalf("find assets for market failed: %v", err)
	}

	if len(ma1.List) != 1 {
		t.Fatal("asset list length failed")
	}

	if *asset1 != *ma1.List[0] {
		t.Fatal("data not match")
	}

	all, err := FindAllAssets(stub, event.Id)
	if err != nil {
		t.Fatal("find assets for event failed: %v", err)
	}

	if len(all.List) != 2 {
		t.Fatal("length of list failed")
	}

	if !(*asset1 == *all.List[0] || *asset1 == *all.List[1] ||
		*asset2 == *all.List[0] || *asset2 == *all.List[1]) {
		t.Fatal("data not match")
	}

	_, existed, err = GetAssetAndCheck(stub, "abc")
	if existed || err == nil {
		t.Fatal("can not find asset with wrong data")
	}

	_, existed, err = GetAssetAndCheck(stub, "a#b#c#d")
	if existed || err != nil {
		t.Fatal("can not find asset with wrong data")
	}

}
