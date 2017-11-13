package lmsr

import (
	"diviner/common/util"

	"github.com/gogo/protobuf/proto"
	perrors "github.com/pkg/errors"
)

// InitShares ...
func InitShares(mkt string, outcomes []*Outcome) []*Share {
	shares := make([]*Share, len(outcomes))

	for i, x := range outcomes {
		s := &Share{
			Id:      ShareID(mkt, x.Id),
			Market:  mkt,
			Outcome: x.Id,
			Volume:  0.0,
		}

		shares[i] = s
	}

	return shares
}

// InitPrices ...
func InitPrices(liquidity float64, shares []*Share) []*Price {
	sum, ok := util.FoldLeftFloat64(shares, 0.0, func(a float64, b interface{}) (float64, bool) {
		switch b.(type) {
		case *Share:
			return a + Exp(liquidity, b.(*Share).Volume), true
		default:
			return 0.0, false
		}
	})

	if !ok {
		return nil
	}

	result := make([]*Price, len(shares))

	for i, x := range shares {
		result[i] = &Price{
			Share: x.Id,
			Price: Exp(liquidity, x.Volume) / sum,
		}
	}

	return result
}

// InitMarket ...
func InitMarket(user string, event *Event) *Market {
	mkt := &Market{
		Id:      MarketID(),
		Event:   event.Id,
		User:    user,
		Settled: false,
	}

	mkt.Shares = InitShares(mkt.Id, event.Outcomes)
	return mkt
}

// NewMarketWithFund ...
func NewMarketWithFund(user string, event *Event, fund float64) (*Market, error) {
	if fund <= 0 {
		return nil, perrors.Errorf("fund must larger than 0: %v", fund)
	}

	mkt := InitMarket(user, event)
	mkt.Fund = fund
	mkt.Liquidity = Liquidity(fund, len(event.Outcomes))
	mkt.Cost = fund
	mkt.Prices = InitPrices(mkt.Liquidity, mkt.Shares)

	return mkt, nil
}

// NewMarketWithLiquidity ...
func NewMarketWithLiquidity(user string, event *Event, liquidity float64) (*Market, error) {
	if liquidity <= 0 {
		return nil, perrors.Errorf("liquidity must larger than 0: %v", liquidity)
	}

	len := len(event.Outcomes)
	if len < 2 {
		return nil, perrors.Errorf("length of outcomes must larger than 1: %v", len)
	}

	fund := Fund(liquidity, len)

	mkt := InitMarket(user, event)
	mkt.Fund = fund
	mkt.Liquidity = liquidity
	mkt.Prices = InitPrices(liquidity, mkt.Shares)
	mkt.Cost = fund

	return mkt, nil
}

// UnmarshalMarket ...
func UnmarshalMarket(data []byte) (*Market, error) {
	mkt := &Market{}
	err := proto.Unmarshal(data, mkt)
	return mkt, err
}

// MarshalMarket ...
func MarshalMarket(m *Market) ([]byte, error) {
	return proto.Marshal(m)
}
