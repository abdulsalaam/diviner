package lmsr

import (
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
	sum := foldLeftFloat64(shares, 0.0, func(a float64, b interface{}) float64 {
		return a + Exp(liquidity, b.(*Share).Volume)
	})

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
