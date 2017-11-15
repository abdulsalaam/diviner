package lmsr

import "testing"

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
		liq = Liquidity(fund, len(e.Outcomes))
	} else {
		liq = num
		fund = Fund(liq, len(e.Outcomes))
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
	evt, _ := NewEvent(user, title, outcomes[0], outcomes[1])

	mkt1, err := NewMarketWithFund(user, evt, 100.0)
	if err != nil {
		t.Fatal("create market with fund failed")
	}

	checkMarket(mkt1, user, evt, 100.0, true, t)

	mkt2, err := NewMarketWithLiquidity(user, evt, 100.0)
	if err != nil {
		t.Fatal("create market with liquidity failed")
	}

	checkMarket(mkt2, user, evt, 100.0, false, t)

	if CmpMarket(mkt1, mkt2) {
		t.Fatal("two market must not equal")
	}
}

func TestWrongData(t *testing.T) {
	evt, _ := NewEvent(user, title, outcomes[0], outcomes[1])
	_, err := NewMarketWithFund(user, evt, 0)
	if err == nil {
		t.Fatal("fund = 0 must return error")
	}

	_, err = NewMarketWithLiquidity(user, evt, 0)
	if err == nil {
		t.Fatal("liquidity = 0 must return error")
	}

	evt.Approved = true
	_, err = NewMarketWithFund(user, evt, 100.0)
	if err == nil {
		t.Fatal("can not create market with approved event")
	}

	_, err = NewMarketWithLiquidity(user, evt, 100.0)
	if err == nil {
		t.Fatal("can not create market with approved event")
	}
}

func TestMarshal(t *testing.T) {
	evt, _ := NewEvent(user, title, outcomes[0], outcomes[1])

	m1, _ := NewMarketWithFund(user, evt, 100.0)

	byte1, err := MarshalMarket(m1)
	if err != nil {
		t.Fatal("marshal market failed")
	}

	m2, err := UnmarshalMarket(byte1)
	if err != nil {
		t.Fatal("unmarshal market failed")
	}

	if !CmpMarket(m1, m2) {
		t.Fatal("data not match")
	}
}
