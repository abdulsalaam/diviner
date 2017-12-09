package market

import (
	"diviner/common/lmsr"
	"testing"
	"time"
)

var (
	user     = "user1"
	title    = "gogogo"
	outcomes = []string{"yes", "no"}
)

func checkMarket(m *Market, u string, e *Event, num float64, flag bool, t *testing.T) {
	var fund float64
	var liq float64

	if flag {
		fund = num
		liq = lmsr.Liquidity(fund, len(e.Outcomes))
	} else {
		liq = num
		fund = lmsr.Fund(liq, len(e.Outcomes))
	}

	if m.GetUser() != u {
		t.Error("user not match")
	}

	if m.GetEvent() != e.Id {
		t.Error("event not match")
	}

	if m.GetLiquidity() != liq {
		t.Error("liquidity not match")
	}

	if m.GetFund() != fund {
		t.Error("fund not match")
	}

	if m.GetCost() != fund {
		t.Error("cost wrong")
	}

	if m.GetSettled() {
		t.Error("settled wrong")
	}

	var price float64 = 1.0 / float64(len(e.Outcomes))

	for k, v := range m.Shares {
		if v != 0 {
			t.Errorf("volume wrong")
		}

		if p, ok := m.Prices[k]; !ok {
			t.Errorf("can not find share price: %s", k)
		} else if p != price {
			t.Error("price wrong")
		}

	}
}

func TestNewMarket(t *testing.T) {

	nextYear := time.Now().AddDate(1, 0, 0)
	start := nextYear.AddDate(0, -2, 0)
	end := nextYear.AddDate(0, 0, -1)

	evt, _ := NewEvent(user, title, nextYear, outcomes[0], outcomes[1])

	mkt1, err := NewMarketWithFund(user, evt, start, end, 100.0)
	if err != nil {
		t.Fatal("create market with fund failed")
	}

	checkMarket(mkt1, user, evt, 100.0, true, t)

	mkt2, err := NewMarketWithLiquidity(user, evt, start, end, 100.0)
	if err != nil {
		t.Fatal("create market with liquidity failed")
	}

	checkMarket(mkt2, user, evt, 100.0, false, t)

	if mkt1.String() == mkt2.String() {
		t.Fatal("two makets must not equal")
	}
}

func TestWrongData(t *testing.T) {
	nextYear := time.Now().AddDate(1, 0, 0)
	start := nextYear.AddDate(0, -2, 0)
	end := nextYear.AddDate(0, 0, -1)

	evt, _ := NewEvent(user, title, nextYear, outcomes[0], outcomes[1])
	_, err := NewMarketWithFund(user, evt, start, end, 0)
	if err == nil {
		t.Fatal("fund = 0 must return error")
	}

	_, err = NewMarketWithLiquidity(user, evt, start, end, 0)
	if err == nil {
		t.Fatal("liquidity = 0 must return error")
	}

	evt.Approved = true
	_, err = NewMarketWithFund(user, evt, start, end, 100.0)
	if err == nil {
		t.Fatal("can not create market with approved event")
	}

	evt.Approved = false
	evt.Allowed = false
	_, err = NewMarketWithFund(user, evt, start, end, 100.0)
	if err == nil {
		t.Fatal("can not create market with rejected event")
	}

	_, err = NewMarketWithLiquidity(user, evt, start, end, 100.0)
	if err == nil {
		t.Fatal("can not create market with approved event")
	}
}

func TestMarshal(t *testing.T) {
	nextYear := time.Now().AddDate(1, 0, 0)
	start := nextYear.AddDate(0, -2, 0)
	end := nextYear.AddDate(0, 0, -1)

	evt, _ := NewEvent(user, title, nextYear, outcomes[0], outcomes[1])

	m1, _ := NewMarketWithFund(user, evt, start, end, 100.0)

	byte1, err := MarshalMarket(m1)
	if err != nil {
		t.Fatal("marshal market failed")
	}

	m2, err := UnmarshalMarket(byte1)
	if err != nil {
		t.Fatal("unmarshal market failed")
	}

	if m1.String() != m2.String() {
		t.Fatal("data not match")
	}
}
