package lmsr

import (
	"fmt"
	"math"
	"os"
	"testing"

	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"diviner/common/cast"
	"diviner/common/csp"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"

	"github.com/google/uuid"
)

var (
	stub     = ccc.NewMockStub("oracle", NewLMSRChaincode())
	title    = "gogogo"
	outcomes = []string{"yes", "no"}
	balance  = 10000.0
	member   *pbm.Member
	event0   *pbl.Event
	market1  *pbl.Market
	market2  *pbl.Market
)

func TestMain(m *testing.M) {
	priv, _ := csp.GeneratePrivateTempKey()
	member, _ = pbm.NewMember(priv, balance)
	event0, _ = pbl.NewEvent(member.Id, title, outcomes[0], outcomes[1])
	market1, _ = pbl.NewMarketWithFund(member.Id, event0, 100.0)
	market2, _ = pbl.NewMarketWithFund(member.Id, event0, 200.0)

	txid := uuid.New().String()

	stub.MockTransactionStart(txid)

	ccc.PutMessage(stub, member.Id, member)

	ccc.PutMessage(stub, event0.Id, event0)

	eid, mid, _ := pbl.SepMarketID(market1.Id)
	ccc.PutMessageWithCompositeKey(stub, market1, pbl.MarketKey, eid, mid)

	eid, mid, _ = pbl.SepMarketID(market2.Id)
	ccc.PutMessageWithCompositeKey(stub, market2, pbl.MarketKey, eid, mid)

	stub.MockTransactionEnd(txid)

	m.Run()
}

func TestMarkets(t *testing.T) {

	resp := ccc.MockInvokeWithString(stub, "markets", "a", "b")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke with wrong data")
	}

	resp = ccc.MockInvokeWithString(stub, "markets", event0.Id)
	if !ccc.OK(&resp) {
		t.Fatalf("invoke markets failed: %s", resp.Message)
	}

	markets, _ := pbl.UnmarshalMarkets(resp.Payload)
	if len(markets.List) != 2 {
		t.Fatal("list length failed")
	}

	if !(pbl.CmpMarket(market1, markets.List[0]) ||
		pbl.CmpMarket(market1, markets.List[1]) ||
		pbl.CmpMarket(market2, markets.List[0]) ||
		pbl.CmpMarket(market2, markets.List[1])) {
		t.Fatal("data not match")
	}

}

func TestSimulator(t *testing.T) {
	txid := uuid.New().String()
	priv0, _ := csp.GeneratePrivateTempKey()
	priv1, _ := csp.GeneratePrivateTempKey()
	priv2, _ := csp.GeneratePrivateTempKey()
	priv3, _ := csp.GeneratePrivateTempKey()
	priv4, _ := csp.GeneratePrivateTempKey()
	priv5, _ := csp.GeneratePrivateTempKey()

	mem0, _ := pbm.NewMember(priv0, balance)
	mem1, _ := pbm.NewMember(priv1, balance)
	mem2, _ := pbm.NewMember(priv2, balance)
	mem3, _ := pbm.NewMember(priv3, balance)
	mem4, _ := pbm.NewMember(priv4, balance)
	mem5, _ := pbm.NewMember(priv5, 0.01)

	event, _ := pbl.NewEvent(mem0.Id, "GO", "Yes", "No")
	market, _ := pbl.NewMarketWithLiquidity(mem0.Id, event, 100.0)
	yes := pbl.ShareID(market.Id, event.Outcomes[0].Id)
	no := pbl.ShareID(market.Id, event.Outcomes[1].Id)

	if _, ok := market.Shares[yes]; !ok {
		fmt.Println("share id failed")
		os.Exit(-1)
	}

	if _, ok := market.Shares[no]; !ok {
		fmt.Println("share id failed")
		os.Exit(-1)
	}

	stub.MockTransactionStart(txid)

	// Member
	ccc.PutMessage(stub, mem0.Id, mem0)
	ccc.PutMessage(stub, mem1.Id, mem1)
	ccc.PutMessage(stub, mem2.Id, mem2)
	ccc.PutMessage(stub, mem3.Id, mem3)
	ccc.PutMessage(stub, mem4.Id, mem4)
	ccc.PutMessage(stub, mem5.Id, mem5)

	// event and market

	ccc.PutMessage(stub, event.Id, event)
	ccu.PutMarket(stub, market)
	mem0.Balance -= market.Fund
	ccc.PutMessage(stub, mem0.Id, mem0)
	stub.MockTransactionEnd(txid)

	fmt.Println("step 0: market")
	dumpMarket(market)
	yesVolume := 0.0
	noVolume := 0.0
	price := 0.0
	checkShares(t, market, yes, no, yesVolume, noVolume)

	// step 1: mem1 buy yes 100
	txVolume := 100.0
	yesVolume += txVolume
	mem1, market, price = myCheck(t, mem1, market, "buy", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 1: mem1 buy yes 100: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)

	// step 2: mem2 buy yes 40
	txVolume = 40
	yesVolume += txVolume
	mem2, market, price = myCheck(t, mem2, market, "buy", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 2: mem2 buy yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem2)

	// step 3: mem3 buy no 20
	txVolume = 20
	noVolume += txVolume
	mem3, market, price = myCheck(t, mem3, market, "buy", no, txVolume)
	fmt.Printf("\nstep 3: mem3 buy no 20: %v\n", price)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	dumpMarket(market)
	dumpMember(mem3)

	// step 4: mem1 buy yes 50
	txVolume = 50
	yesVolume += txVolume
	mem1, market, price = myCheck(t, mem1, market, "buy", yes, txVolume)
	fmt.Printf("\nstep 4: mem1 buy yes 50: %v\n", price)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	dumpMarket(market)
	dumpMember(mem1)

	// step 5: mem2 buy yes 100
	txVolume = 100
	yesVolume += txVolume
	mem2, market, price = myCheck(t, mem2, market, "buy", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 5: mem2 buy yes 100: %v\n", price)
	dumpMarket(market)
	dumpMember(mem2)

	// step 6: mem4 buy no 50
	txVolume = 50
	noVolume += txVolume
	mem4, market, price = myCheck(t, mem4, market, "buy", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 6: mem4 buy no 50: %v\n", price)
	dumpMarket(market)
	dumpMember(mem4)

	// step 7: mem1 sell yes 40
	txVolume = 40
	yesVolume -= txVolume
	mem1, market, price = myCheck(t, mem1, market, "sell", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 7: mem 1 sell yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)

	// step 8: mem3 buy no 30
	txVolume = 30
	noVolume += txVolume
	mem3, market, price = myCheck(t, mem3, market, "buy", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 8: mem3 buy no 30: %v\n", price)
	dumpMarket(market)
	dumpMember(mem3)

	// step 9: mem2 buy yes 40
	txVolume = 40
	yesVolume += txVolume
	mem2, market, price = myCheck(t, mem2, market, "buy", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 9: mem2 buy yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem2)

	// step 10: mem4 buy no 300
	txVolume = 300
	noVolume += txVolume
	mem4, market, price = myCheck(t, mem4, market, "buy", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 10: mem4 buy no 300: %v\n", price)
	dumpMarket(market)
	dumpMember(mem4)

	// step 11: mem3 sell no 10
	txVolume = 10
	noVolume -= txVolume
	mem3, market, price = myCheck(t, mem3, market, "sell", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 11: mem3 sell no 10: %v\n", price)
	dumpMarket(market)
	dumpMember(mem3)

	// step 12: mem4 buy no 150
	txVolume = 150
	noVolume += txVolume
	mem4, market, price = myCheck(t, mem4, market, "buy", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 12: mem4 buy no 150: %v\n", price)
	dumpMarket(market)
	dumpMember(mem4)

	// step 13: mem2 sell yes 40
	txVolume = 40
	yesVolume -= txVolume
	mem2, market, price = myCheck(t, mem2, market, "sell", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 13: mem2 sell yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem2)

	// step 14: mem3 buy no 20
	txVolume = 20
	noVolume += txVolume
	mem3, market, price = myCheck(t, mem3, market, "buy", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 14: mem3 buy no 20: %v\n", price)
	dumpMarket(market)
	dumpMember(mem3)

	// step 15: mem1 buy yes 40
	txVolume = 40
	yesVolume += txVolume
	mem1, market, price = myCheck(t, mem1, market, "buy", yes, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 15: mem1 buy yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)

	// step 16: mem3 buy no 200
	txVolume = 200
	noVolume += txVolume
	mem3, market, price = myCheck(t, mem3, market, "buy", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 16: mem3 buy no 200: %v\n", price)
	dumpMarket(market)
	dumpMember(mem3)

	// step 17: mem4 sell no 100
	txVolume = 100
	noVolume -= txVolume
	mem4, market, price = myCheck(t, mem4, market, "sell", no, txVolume)
	checkShares(t, market, yes, no, yesVolume, noVolume)
	fmt.Printf("\nstep 17: mem4 sell no 100: %v\n", price)
	dumpMarket(market)
	dumpMember(mem4)

	members := make(map[string]*pbm.Member)
	members[mem1.Id] = mem1
	members[mem2.Id] = mem2
	members[mem3.Id] = mem3
	members[mem4.Id] = mem4
	members[mem5.Id] = mem5

	my := 0.0
	mn := 0.0

	totalShares1 := 0
	totalShares2 := 0

	// find asset in ledger by member assets
	for _, v := range members {
		for x := range v.Assets {
			_, existed, err := ccu.GetAssetAndCheck(stub, x)
			if err != nil {
				t.Fatalf("find asset failed: %s, %v", x, err)
			}

			if !existed {
				t.Fatalf("asset not found %s", x)
			}
			totalShares1 += 1
		}
	}

	// sum share volume by member assets
	for _, v := range members {
		k1 := pbl.AssetID(v.Id, yes)
		k2 := pbl.AssetID(v.Id, no)

		my += v.Assets[k1]
		mn += v.Assets[k2]
	}

	if yesVolume != my || noVolume != mn {
		t.Fatal("member volume falied: %v, %v, %v, %v", yesVolume, my, noVolume, mn)
	}

	// find all assets by event id
	resp := ccc.MockInvokeWithString(stub, "assets", event.Id)
	if !ccc.OK(&resp) {
		t.Fatalf("get assets by event (%s) failed: %v", event.Id, resp.Message)
	}

	assets, _ := pbl.UnMarshalAssets(resp.Payload)
	totalVolume := 0.0
	totalShares2 = len(assets.List)

	if totalShares1 != totalShares2 {
		t.Fatal("asset count failed: %v, %v", totalShares1, totalShares2)
	}

	for _, x := range assets.List {
		_, _, _, user, _ := pbl.SepAssetID(x.Id)

		v1 := members[user].Assets[x.Id]
		v2 := x.Volume

		if v1 != v2 {
			t.Fatal("assert in ledger failed: %s, %s, %v, %v", user, x.Id, v1, v2)
		}

		totalVolume += v2
	}

	if totalVolume != yesVolume+noVolume {
		t.Fatalf("total volume failed: %v, %v, %v", totalVolume, yesVolume, noVolume)
	}

	// find all assets by market id
	evt, mkt, _ := pbl.SepMarketID(market.Id)
	resp = ccc.MockInvokeWithString(stub, "assets", evt, mkt)
	if !ccc.OK(&resp) {
		t.Fatalf("get assets by event (%s) failed: %v", event.Id, resp.Message)
	}

	assets, _ = pbl.UnMarshalAssets(resp.Payload)
	totalVolume = 0.0
	totalShares2 = len(assets.List)

	if totalShares1 != totalShares2 {
		t.Fatal("asset count failed: %v, %v", totalShares1, totalShares2)
	}

	for _, x := range assets.List {
		_, _, _, user, _ := pbl.SepAssetID(x.Id)

		v1 := members[user].Assets[x.Id]
		v2 := x.Volume

		if v1 != v2 {
			t.Fatal("assert in ledger failed: %s, %s, %v, %v", user, x.Id, v1, v2)
		}

		totalVolume += v2
	}

	if totalVolume != yesVolume+noVolume {
		t.Fatalf("total volume failed: %v, %v, %v", totalVolume, yesVolume, noVolume)
	}

	// approve event and check all
	result := event.Outcomes[0].Id
	resultShare := pbl.ShareID(market.Id, result)
	resp = ccc.MockInvokeWithString(stub, "approve", event.Id, result)
	if !ccc.OK(&resp) {
		t.Fatalf("appove event %s failed: %s", event.Id, resp.Message)
	}

	eventResult, err := ccu.FindEvent(stub, event.Id)
	if err != nil {
		t.Fatal(err)
	}
	if !eventResult.Approved {
		t.Fatal("event is not approved")
	}

	marketResult, _, _ := ccu.GetMarketAndCheck(stub, market.Id)
	if !marketResult.Settled {
		t.Fatal("market is not settled")
	}

	owner, _, _ := ccu.GetMemberAndCheck(stub, market.User)
	fmt.Printf("owner balance %v, %v, %v, %v\n", owner.Balance, mem0.Balance, marketResult.Cost, marketResult.Shares[resultShare])
	if owner.Balance != mem0.Balance+marketResult.Cost-market.Shares[resultShare] {
		t.Fatalf("owner balance failed: %v, %v, %v, %v", owner.Balance, mem0.Balance, marketResult.Cost, marketResult.Shares[resultShare])
	}

	for _, x := range members {
		x1, _, _ := ccu.GetMemberAndCheck(stub, x.Id)
		if len(x1.Assets) > 0 {
			t.Fatalf("member (%s) assets must be empty after settled", x1.Id)
		}

		aid := pbl.AssetID(x1.Id, resultShare)
		fmt.Printf("member balance %v, %v, %v\n", x1.Balance, x.Balance, x.Assets[aid])
		if x1.Balance != x.Balance+x.Assets[aid] {
			t.Fatalf("member (%s) balance failed: %v, %v, %v", x1.Id, x1.Balance, x.Balance, x.Assets[aid])
		}

		/*for s := range market.Shares {
			a := pbl.AssetID(x1.Id, s)
			if _, ok := x1.Assets[a]; ok {
				t.Fatalf("asset (%s) must be removed after settled", a)
			}
		}*/

	}

	// test market failed
	resp = ccc.MockInvokeWithString(stub, "settle", market.Id, result)
	if ccc.OK(&resp) {
		t.Fatalf("can not settle a settled market")
	}

	// test wrong with mem5
	resp = ccc.MockInvokeWithString(stub, "buy", mem5.Id, yes, "100")
	if ccc.OK(&resp) {
		t.Fatalf("can not buy without enough balance")
	}

	resp = ccc.MockInvokeWithString(stub, "sell", mem5.Id, yes, "100")
	if ccc.OK(&resp) {
		t.Fatalf("can not sell without enough asset volume")
	}

}

func TestWrongData(t *testing.T) {
	resp := ccc.MockInvokeWithString(stub, "assets")
	if ccc.OK(&resp) {
		t.Fatal("can not find any assets with wrong data")
	}

	resp = ccc.MockInvokeWithString(stub, "buy", "a", "b", "c", "d", "e")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke buy with wrong data")
	}

	resp = ccc.MockInvokeWithString(stub, "sell", "a", "b", "c", "d", "e")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke sell with wrong data")
	}

	resp = ccc.MockInvokeWithString(stub, "sell", "a", "b", "c")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke buy with wrong type for volume")
	}

	resp = ccc.MockInvokeWithString(stub, "sell", "a", "b", "0")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke buy with volume <= 0")
	}

	memId := "abc"
	evtId := pbl.EventID()
	mktId := pbl.MarketID(evtId)
	outcomeId := pbl.OutcomeID(evtId, 10)
	shareId := pbl.ShareID(mktId, outcomeId)

	resp = ccc.MockInvokeWithString(stub, "buy", memId, shareId, "10")
	if ccc.OK(&resp) {
		t.Fatalf("can not invoke with non-existed data")
	}

	shareId = pbl.ShareID(market1.Id, event0.Outcomes[0].Id)

	resp = ccc.MockInvokeWithString(stub, "buy", memId, shareId+"1", "10")
	if ccc.OK(&resp) {
		t.Fatalf("can not invoke with non-existed data")
	}

	resp = ccc.MockInvokeWithString(stub, "buy", memId, shareId, "10")
	if ccc.OK(&resp) {
		t.Fatalf("can not invoke with non-existed data")
	}

}

func checkTx(t *testing.T, cmd, user, share, volume string) float64 {
	resp := ccc.MockInvokeWithString(stub, cmd, user, share, volume)
	if !ccc.OK(&resp) {
		t.Fatalf("buy failed: %s", resp.Message)
	}

	price, err := cast.BytesToFloat64(resp.Payload)
	if err != nil {
		t.Fatalf("price bytes failed: %v", err)
	}

	return price
}

func checkMember(t *testing.T, old *pbm.Member, share string, volume, price float64) *pbm.Member {
	memtmp, _, _ := ccu.GetMemberAndCheck(stub, old.Id)
	asstmp := pbl.AssetID(old.Id, share)
	if memtmp.Assets[asstmp] != old.Assets[asstmp]+volume {
		t.Fatalf("asset volume failed: %v, %v, %v", memtmp.Assets[asstmp], old.Assets[asstmp], volume)
	}

	if memtmp.Balance != old.Balance-price {
		t.Fatal("member balance failed: %v, %v, %v", memtmp.Balance, old.Balance, price)
	}

	return memtmp
}

func checkMarket(t *testing.T, old *pbl.Market, share string, volume, price float64) *pbl.Market {
	mkttmp, _, _ := ccu.GetMarketAndCheck(stub, old.Id)
	if mkttmp.Shares[share] != old.Shares[share]+volume {
		t.Fatal("market share failed: %v, %v, %v", old.Shares[share], old.Shares[share], volume)
	}

	if price != mkttmp.Cost-old.Cost {
		t.Fatal("price and cost failed: %v, %v, %v", mkttmp.Cost, old.Cost, price)
	}

	return mkttmp
}

func checkShares(t *testing.T, market *pbl.Market, yes, no string, y, n float64) {
	if market.Shares[yes] != y || market.Shares[no] != n {
		t.Fatal("share amount failed: %v, %v, %v, %v", market.Shares[yes], y, market.Shares[no], n)
	}
}

func myCheck(t *testing.T, mem *pbm.Member, market *pbl.Market, cmd, share string, volume float64) (*pbm.Member, *pbl.Market, float64) {
	price := checkTx(t, cmd, mem.Id, share, fmt.Sprintf("%v", volume))

	if cmd == "sell" {
		price = -price
		volume = -volume
	}
	mkt := checkMarket(t, market, share, volume, price)
	user := checkMember(t, mem, share, volume, price)

	return user, mkt, math.Abs(price)
}

func dumpMarket(m *pbl.Market) {
	fmt.Println("--- market ---")
	fmt.Println("id: ", m.Id)
	fmt.Println("liquidity: ", m.Liquidity)
	fmt.Println("Fund: ", m.Fund)
	fmt.Println("Cost: ", m.Cost)
	fmt.Println("shares and prices: ")
	for k, v := range m.Shares {
		fmt.Println(k, v, m.Prices[k])
	}
}

func dumpMember(m *pbm.Member) {
	fmt.Println("--- member ---")
	fmt.Println("id: ", m.Id)
	fmt.Println("balance: ", m.Balance)
	fmt.Println("assets: ")
	for k, v := range m.Assets {
		fmt.Println(k, v)
	}
}
