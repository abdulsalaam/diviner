package lmsr

import (
	"math"

	"github.com/golang/protobuf/proto"
	perrors "github.com/pkg/errors"
)

// InitMarket ...
func InitMarket(user string, event *Event) *Market {
	mkt := &Market{
		Id:      MarketID(event.Id),
		Event:   event.Id,
		User:    user,
		Shares:  make(map[string]float64),
		Settled: false,
	}

	for _, x := range event.Outcomes {
		id := ShareID(mkt.Id, x.Id)
		mkt.Shares[id] = 0.0
	}

	return mkt
}

// NewMarketWithFund ...
func NewMarketWithFund(user string, event *Event, fund float64) (*Market, error) {
	if fund <= 0 {
		return nil, perrors.Errorf("fund must larger than 0: %v", fund)
	}

	if event.Approved {
		return nil, perrors.Errorf("event is approved")
	}

	mkt := InitMarket(user, event)
	mkt.Fund = fund
	mkt.Liquidity = Liquidity(fund, len(event.Outcomes))
	mkt.Cost = fund
	Reprices(mkt)

	return mkt, nil
}

// NewMarketWithLiquidity ...
func NewMarketWithLiquidity(user string, event *Event, liquidity float64) (*Market, error) {
	if liquidity <= 0 {
		return nil, perrors.Errorf("liquidity must larger than 0: %v", liquidity)
	}

	if event.Approved {
		return nil, perrors.Errorf("event is approved")
	}

	len := len(event.Outcomes)
	if len < 2 {
		return nil, perrors.Errorf("length of outcomes must larger than 1: %v", len)
	}

	fund := Fund(liquidity, len)

	mkt := InitMarket(user, event)
	mkt.Fund = fund
	mkt.Liquidity = liquidity
	mkt.Cost = fund
	Reprices(mkt)

	return mkt, nil
}

// CmpMarket ...
func CmpMarket(m1, m2 *Market) bool {
	if m1.Id != m2.Id || m1.User != m2.User || m1.Event != m2.Event ||
		m1.Liquidity != m2.Liquidity || m1.Fund != m2.Fund || m1.Cost != m2.Cost || m1.Settled != m2.Settled {
		return false
	}

	if len(m1.Shares) != len(m2.Shares) {
		return false
	}

	if len(m1.Prices) != len(m2.Prices) {
		return false
	}

	for k, v1 := range m1.Shares {
		v2, ok := m2.Shares[k]
		if !ok {
			return false
		}

		if v1 != v2 {
			return false
		}
	}

	for k, v1 := range m1.Prices {
		v2, ok := m2.Prices[k]
		if !ok {
			return false
		}
		if v1 != v2 {
			return false
		}
	}

	return true
}

func EstimateMarket(market *Market, share string, volume float64) (float64, error) {
	if volume == 0 {
		return 0, perrors.Errorf("volume can not equal 0")
	}

	orgVol, ok := market.Shares[share]
	if !ok {
		return 0, perrors.Errorf("can not find share (%s) in market (%s)", share, market.Id)
	}

	if orgVol+volume < 0 {
		return 0, perrors.Errorf("volume is not enough to sell. org %v, but %v", orgVol, math.Abs(volume))
	}

	sum := 0.0

	for k, v := range market.Shares {
		if k == share {
			sum += Exp(market.Liquidity, v+volume)
		} else {
			sum += Exp(market.Liquidity, v)
		}
	}

	cost := market.Liquidity * math.Log(sum)

	return cost - market.Cost, nil
}

func UpdateMarket(market *Market, share string, volume float64) (float64, error) {
	price, err := EstimateMarket(market, share, volume)
	if err != nil {
		return 0, nil
	}

	market.Shares[share] += volume
	market.Cost += price
	Reprices(market)
	return price, nil
}

func Reprices(market *Market) {
	sum := 0.0
	for _, v := range market.Shares {
		sum += Exp(market.Liquidity, v)
	}

	market.Prices = make(map[string]float64)
	for k, v := range market.Shares {
		market.Prices[k] = Exp(market.Liquidity, v) / sum
	}
}

// UnmarshalMarket ...
func UnmarshalMarket(data []byte) (*Market, error) {
	mkt := &Market{}
	if err := proto.Unmarshal(data, mkt); err != nil {
		return nil, err
	}
	return mkt, nil
}

// MarshalMarket ...
func MarshalMarket(m *Market) ([]byte, error) {
	return proto.Marshal(m)
}

// UnmarshalMarkets ...
func UnmarshalMarkets(data []byte) (*Markets, error) {
	lst := &Markets{}
	if err := proto.Unmarshal(data, lst); err != nil {
		return nil, err
	}
	return lst, nil
}

// MarshalMarket ...
func MarshalMarkets(lst *Markets) ([]byte, error) {
	return proto.Marshal(lst)
}
