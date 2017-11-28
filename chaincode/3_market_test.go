package main

import (
	"testing"

	ccc "diviner/chaincode/common"
	ccu "diviner/chaincode/util"
	"diviner/common/cast"
	"diviner/common/csp"
	pbl "diviner/protos/lmsr"
	pbm "diviner/protos/member"
)

var (
	mktCreator *pbm.Member
	mktEvent   *pbl.Event
)

func initEvent(t *testing.T) {
	if mktCreator == nil {
		key, _ := csp.GeneratePrivateTempKey()
		mktCreator, _ = pbm.NewMember(key, balance)
	}

	if _, ok, err := ccc.GetStateAndCheck(stub, mktCreator.Id); err != nil {
		t.Fatal(err)
	} else if !ok {
		b0, _ := pbm.Marshal(mktCreator)
		ccc.MockInvoke(stub, fcnCreate, ccMember, b0)
	}

	if mktEvent == nil {
		mktEvent, _ = pbl.NewEvent(mktCreator.Id, title, outcomes[0], outcomes[1])
	}

	if _, ok, err := ccc.GetStateAndCheck(stub, mktEvent.Id); err != nil {
		t.Fatal(err)
	} else if !ok {
		b0, _ := pbl.MarshalEvent(mktEvent)
		ccc.MockInvoke(stub, fcnCreate, ccEvent, b0)
	}
}

func TestCreateWithFund(t *testing.T) {
	initEvent(t)
	num, _ := cast.ToBytes(100.0)
	resp := ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte(mktCreator.Id), []byte(mktEvent.Id), num, []byte{1})

	if !ccc.OK(&resp) {
		t.Fatalf("create market failed: %s", resp.Message)
	}

	mkt, _ := pbl.UnmarshalMarket(resp.Payload)

	resp = ccc.MockInvoke(stub, fcnQuery, ccMarket, []byte(mkt.Id))
	if !ccc.OK(&resp) {
		t.Fatalf("query market failed: %s", resp.Message)
	}

	mkt2, _ := pbl.UnmarshalMarket(resp.Payload)

	if !pbl.CmpMarket(mkt, mkt2) {
		t.Fatal("data error")
	}

	if mkt2.Fund != 100.0 {
		t.Fatalf("fund failed: %v", mkt2.Fund)
	}

	if mkt2.Liquidity != pbl.Liquidity(100.0, len(mktEvent.Outcomes)) {
		t.Fatalf("liquidity failed: %v", mkt2.Liquidity)
	}

	member2, existed, err := ccu.GetMemberAndCheck(stub, mktCreator.Id)
	if err != nil {
		t.Fatalf("query user failed: %s, %v", mktCreator.Id, err)
	} else if !existed {
		t.Fatal("user not found: %s", mktCreator.Id)
	}

	if member2.Balance != (mktCreator.Balance - mkt.Fund) {
		t.Fatal("member balance failed")
	}

	mktCreator.Balance = mktCreator.Balance - mkt.Fund
}

func TestCreateWithLiquidity(t *testing.T) {
	num, _ := cast.ToBytes(100.0)
	resp := ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte(mktCreator.Id), []byte(mktEvent.Id), num, []byte{0})
	if !ccc.OK(&resp) {
		t.Fatalf("create market failed: %s", resp.Message)
	}

	mkt, _ := pbl.UnmarshalMarket(resp.Payload)

	resp = ccc.MockInvoke(stub, fcnQuery, ccMarket, []byte(mkt.Id))
	if !ccc.OK(&resp) {
		t.Fatalf("query market failed: %s", resp.Message)
	}

	mkt2, _ := pbl.UnmarshalMarket(resp.Payload)

	if !pbl.CmpMarket(mkt, mkt2) {
		t.Fatal("data error")
	}

	if mkt2.Liquidity != 100.0 {
		t.Fatalf("liquidity failed: %v", mkt2.Liquidity)
	}

	if mkt2.Fund != pbl.Fund(100.0, len(mktEvent.Outcomes)) {
		t.Fatalf("fund failed: %v", mkt2.Fund)
	}

	member2, existed, err := ccu.GetMemberAndCheck(stub, mktCreator.Id)
	if err != nil {
		t.Fatal("query user failed")
	} else if !existed {
		t.Fatal("user not found")
	}

	if member2.Balance != (mktCreator.Balance - mkt.Fund) {
		t.Fatal("member balance failed")
	}

	mktCreator.Balance = mktCreator.Balance - mkt.Fund
}

func TestSettle(t *testing.T) {
	num, _ := cast.ToBytes(100.0)
	resp := ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte(mktCreator.Id), []byte(mktEvent.Id), num, []byte{0})
	mkt, _ := pbl.UnmarshalMarket(resp.Payload)

	resp = ccc.MockInvoke(stub, fcnSettle, ccMarket, []byte(mkt.Id))
	if !ccc.OK(&resp) {
		t.Fatal("settle failed")
	}

	mkt.Settled = true

	resp = ccc.MockInvoke(stub, fcnQuery, ccMarket, []byte(mkt.Id))
	mkt1, _ := pbl.UnmarshalMarket(resp.Payload)
	if !mkt1.Settled {
		t.Fatal("settle status failed")
	}

	if !pbl.CmpMarket(mkt, mkt1) {
		t.Fatal("data not match")
	}

	resp = ccc.MockInvoke(stub, fcnSettle, ccMarket, []byte(mkt.Id))
	if ccc.OK(&resp) {
		t.Fatal("can not settle a settled market")
	}

}

func TestWrongData(t *testing.T) {
	tmp, _ := pbl.NewMarketWithFund(mktCreator.Id, mktEvent, 100.0)
	resp := ccc.MockInvoke(stub, fcnQuery, ccMarket, []byte(tmp.Id))
	if ccc.OK(&resp) {
		t.Fatal("can not query non-existed market")
	}
	//fmt.Println(resp)

	resp = ccc.MockInvokeWithString(stub, string(fcnQuery), string(ccMarket), "a", "b")
	if ccc.OK(&resp) {
		t.Fatal("can not query with wrong parameters")
	}
	//fmt.Println(resp)

	resp = ccc.MockInvokeWithString(stub, string(fcnSettle), string(ccMarket), tmp.Id)
	if ccc.OK(&resp) {
		t.Fatal("can not settle non-existed market")
	}
	//fmt.Println(resp)

	resp = ccc.MockInvokeWithString(stub, string(fcnSettle), string(ccMarket), "a", "b")
	if ccc.OK(&resp) {
		t.Fatal("can not settle with wrong parameters")
	}
	//fmt.Println(resp)

	resp = ccc.MockInvokeWithString(stub, string(fcnCreate), string(ccMarket), mktCreator.Id, mktEvent.Id, "100.0", "abc")
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market with wrong flag")
	}
	//fmt.Println(resp)

	num, _ := cast.ToBytes(100000000000000.0)
	resp = ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte(mktCreator.Id), []byte(mktEvent.Id), num, []byte{1})
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market without enough balance")
	}
	//fmt.Println(resp)

	resp = ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte(mktCreator.Id), []byte(mktEvent.Id), []byte("abc"), []byte{1})
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market without float64 number")
	}

	//num, _ = cast.ToBytes(100000000000000.0)
	resp = ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte(mktCreator.Id), []byte(mktEvent.Id), []byte("100.0"), []byte{0}, []byte{1})
	if ccc.OK(&resp) {
		t.Fatalf("can not create a market with wrong parameters")
	}

	num, _ = cast.ToBytes(100.0)
	resp = ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte("abc"), []byte(mktEvent.Id), num, []byte{0})
	if ccc.OK(&resp) {
		t.Fatal("can not create market with non-existed user")
	}
	//fmt.Println(resp)

	resp = ccc.MockInvoke(stub, fcnCreate, ccMarket, []byte("abc"), []byte(mktEvent.Id), num, []byte{1})
	if ccc.OK(&resp) {
		t.Fatal("can not create market with non-existed user")
	}

	resp = ccc.MockInvokeWithString(stub, "aaa", string(ccMarket), "aaabbb")
	if ccc.OK(&resp) {
		t.Fatal("can not invoke with wrong function name")
	}
	//fmt.Println(resp)
}
