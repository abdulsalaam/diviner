package market

import (
	"fmt"
	"math"
	"time"

	"diviner/common/lmsr"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// InitMarket ...
func InitMarket(user string, event *Event, start, end time.Time) (*Market, error) {
	if event.Approved {
		return nil, fmt.Errorf("event is approved")
	}

	if !event.Allowed {
		return nil, fmt.Errorf("event is not allowed")
	}

	if len(event.Outcomes) <= 1 {
		return nil, fmt.Errorf("event outcomes must have more than one outcome: %d", len(event.Outcomes))
	}

	if start.After(end) {
		return nil, fmt.Errorf("start time must be before end time: %v, %v", start, end)
	}

	if end.Before(time.Now()) {
		return nil, fmt.Errorf("end time must be after now: %v", end)
	}

	sts, err := ptypes.TimestampProto(start)
	if err != nil {
		return nil, err
	}

	ets, err := ptypes.TimestampProto(end)
	if err != nil {
		return nil, err
	}

	if event.End.Seconds <= ets.Seconds {
		return nil, fmt.Errorf("event end time must be after market end time: %v, %v", event.End, ets)
	}

	mkt := &Market{
		Id:      MarketID(event.Id),
		Event:   event.Id,
		User:    user,
		Shares:  make(map[string]float64),
		Settled: false,
		Allowed: true,
		Start:   sts,
		End:     ets,
	}

	for _, x := range event.Outcomes {
		id := ShareID(mkt.Id, x.Id)
		mkt.Shares[id] = 0.0
	}

	return mkt, nil
}

// NewMarketWithFund ...
func NewMarketWithFund(user string, event *Event, start, end time.Time, fund float64) (*Market, error) {
	if fund <= 0 {
		return nil, fmt.Errorf("fund must larger than 0: %v", fund)
	}

	mkt, err := InitMarket(user, event, start, end)
	if err != nil {
		return nil, err
	}
	mkt.Fund = fund
	mkt.Liquidity = lmsr.Liquidity(fund, len(event.Outcomes))
	mkt.Cost = fund
	Reprices(mkt)

	return mkt, nil
}

// NewMarketWithLiquidity ...
func NewMarketWithLiquidity(user string, event *Event, start, end time.Time, liquidity float64) (*Market, error) {
	if liquidity <= 0 {
		return nil, fmt.Errorf("liquidity must larger than 0: %v", liquidity)
	}

	fund := lmsr.Fund(liquidity, len(event.Outcomes))

	mkt, err := InitMarket(user, event, start, end)
	if err != nil {
		return nil, err
	}
	mkt.Fund = fund
	mkt.Liquidity = liquidity
	mkt.Cost = fund
	Reprices(mkt)

	return mkt, nil
}

// EstimateMarket ...
func EstimateMarket(market *Market, share string, volume float64) (float64, error) {
	if volume == 0 {
		return 0, fmt.Errorf("volume can not equal 0")
	}

	orgVol, ok := market.Shares[share]
	if !ok {
		return 0, fmt.Errorf("can not find share (%s) in market (%s)", share, market.Id)
	}

	if orgVol+volume < 0 {
		return 0, fmt.Errorf("volume is not enough to sell. org %v, but %v", orgVol, math.Abs(volume))
	}

	newshares := make(map[string]float64)

	for k, v := range market.Shares {
		newshares[k] = v
	}

	newshares[share] += volume

	cost := lmsr.Cost(market.Liquidity, newshares)

	return cost - market.Cost, nil
}

// UpdateMarket ...
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

// Reprices ...
func Reprices(market *Market) {
	sum := lmsr.Sum(market.Liquidity, market.Shares)

	market.Prices = make(map[string]float64)
	for k, v := range market.Shares {
		market.Prices[k] = lmsr.Price(market.Liquidity, v, sum)
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

// MarshalMarkets ...
func MarshalMarkets(lst *Markets) ([]byte, error) {
	return proto.Marshal(lst)
}
