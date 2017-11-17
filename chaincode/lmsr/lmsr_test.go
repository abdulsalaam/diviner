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
	mem5, _ := pbm.NewMember(priv5, balance)

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
	stub.MockTransactionEnd(txid)

	fmt.Println("step 0: market")
	dumpMarket(market)
	yesVolume := 0.0
	noVolume := 0.0
	price := 0.0

	// step 1: mem1 buy yes 100
	txVolume := 100.0
	yesVolume += txVolume
	mem1, market, price = myCheck(t, mem1, market, "buy", yes, txVolume)
	/*price := checkTx(t, "buy", mem1.Id, yes, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, yes, txVolume, price)
	mem1 = checkMember(t, mem1, yes, txVolume, price)*/
	fmt.Printf("\nstep 1: mem1 buy yes 100: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)

	// step 2: mem2 buy yes 40
	txVolume = 40
	yesVolume += txVolume
	price = checkTx(t, "buy", mem2.Id, yes, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, yes, txVolume, price)
	mem2 = checkMember(t, mem2, yes, txVolume, price)
	fmt.Printf("\nstep 2: mem2 buy yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem2)

	// step 3: mem3 buy no 20
	txVolume = 20
	noVolume += txVolume
	price = checkTx(t, "buy", mem3.Id, no, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, no, txVolume, price)
	mem3 = checkMember(t, mem3, no, txVolume, price)
	fmt.Printf("\nstep 3: mem3 buy no 20: %v\n", price)
	dumpMarket(market)
	dumpMember(mem3)

	// step 4: mem1 buy yes 50
	txVolume = 50
	yesVolume += txVolume
	price = checkTx(t, "buy", mem1.Id, yes, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, yes, txVolume, price)
	mem1 = checkMember(t, mem1, yes, txVolume, price)
	fmt.Printf("\nstep 4: mem1 buy yes 50: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)

	// step 5: mem2 buy yes 100
	txVolume = 100
	yesVolume += txVolume
	price = checkTx(t, "buy", mem2.Id, yes, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, yes, txVolume, price)
	mem2 = checkMember(t, mem2, yes, txVolume, price)
	fmt.Printf("\nstep 5: mem2 buy yes 100: %v\n", price)
	dumpMarket(market)
	dumpMember(mem2)

	// step 6: mem4 buy no 50
	txVolume = 50
	noVolume += txVolume
	price = checkTx(t, "buy", mem4.Id, no, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, no, txVolume, price)
	mem4 = checkMember(t, mem4, no, txVolume, price)
	fmt.Printf("\nstep 6: mem4 buy no 50: %v\n", price)
	dumpMarket(market)
	dumpMember(mem4)

	// step 7: mem 1 sell yes 40
	txVolume = 40
	yesVolume -= txVolume
	price = checkTx(t, "sell", mem1.Id, yes, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, yes, -txVolume, -price)
	mem1 = checkMember(t, mem1, yes, -txVolume, -price)
	fmt.Printf("\nstep 7: mem 1 sell yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)

	// step 8: mem3 buy no 30
	txVolume = 30
	noVolume += txVolume
	price = checkTx(t, "buy", mem3.Id, no, fmt.Sprintf("%v", txVolume))
	market = checkMarket(t, market, no, txVolume, price)
	mem1 = checkMember(t, mem3, no, txVolume, price)
	fmt.Printf("\nstep 7: mem 1 sell yes 40: %v\n", price)
	dumpMarket(market)
	dumpMember(mem1)
	/*
		// step 9: buy yes 40
		order, _ = mock.users[1].Order(market, yes, 40, 0.99)
		addOrder(&orders, assets, order)

		// step 10: buy no 300
		order, _ = mock.users[3].Order(market, no, 300, 0.99)
		addOrder(&orders, assets, order)

		// step 11: sell no 10
		aid = core.AssetID(mock.users[2].Account, market.Shares[1].ID)
		order, _ = mock.users[2].Sell(assets[aid], 10)
		addOrder(&orders, assets, order)

		// step 12: buy no 150
		order, _ = mock.users[3].Order(market, no, 150, 0.99)
		addOrder(&orders, assets, order)

		// step 13: sell yes 40
		aid = core.AssetID(mock.users[1].Account, market.Shares[0].ID)
		order, _ = mock.users[1].Sell(assets[aid], 40)
		addOrder(&orders, assets, order)

		// step 14: buy no 20
		order, _ = mock.users[2].Order(market, no, 20, 0.99)
		addOrder(&orders, assets, order)

		// step 15: buy yes 40
		order, _ = mock.users[0].Order(market, yes, 40, 0.99)
		addOrder(&orders, assets, order)

		// step 16: buy no 200
		order, _ = mock.users[2].Order(market, no, 200, 0.99)
		addOrder(&orders, assets, order)

		// step 17: sell no 100
		aid = core.AssetID(mock.users[3].Account, market.Shares[1].ID)
		order, _ = mock.users[3].Sell(assets[aid], 100)
		addOrder(&orders, assets, order)
	*/
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
	memtmp, _, _ := ccu.GetMember(stub, old.Id)
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
