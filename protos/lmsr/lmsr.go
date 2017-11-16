package lmsr

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	// Sep ...
	Sep = "#"

	// MarketKey ...
	MarketKey = "market"
	AssetKey  = "asset"
)

func id() string {
	return uuid.New().String()
}

func prefix(p string) string {
	return p + id()
}

func combine(with string, a, b interface{}) string {
	return fmt.Sprintf("%v"+with+"%v", a, b)
}

// EventID ...
func EventID() string {
	return prefix("")
}

// MarketID ...
func MarketID(evt string) string {
	mkt := id()
	return combine(Sep, evt, mkt)
}

// OutcomeID ...
func OutcomeID(evt string, index int) string {
	return combine("@", evt, index)
}

// ShareID ...
func ShareID(mkt, outcome string) string {
	return combine(Sep, mkt, outcome)
}

// AssetID ...
func AssetID(user, share string) string {
	return combine(Sep, share, user)
}

func Split(str, sep string) (string, string, bool) {
	tmp := strings.Split(str, sep)
	len := len(tmp)
	if len > 1 {
		return strings.Join(tmp[0:len-1], sep), tmp[len-1], true
	}
	return "", "", false
}

// SepAssetID ...
func SepAssetID(id string) (string, string, string, string, bool) {
	tmp := strings.SplitN(id, Sep, 4)

	if len(tmp) != 4 {
		return "", "", "", "", false
	}

	return tmp[0], tmp[1], tmp[2], tmp[3], true // evt, mkt, outcome, user, true

}

// SepMarketID
func SepMarketID(id string) (string, string, bool) {
	return Split(id, Sep)
}

// SepShareID ...
func SepShareID(id string) (string, string, bool) {
	return Split(id, Sep)
}

// SepOutcomeID ...
func SepOutcomeID(id string) (string, int, bool) {
	s1, s2, ok := Split(id, "@")
	if !ok {
		return "", 0, false
	}

	idx, err := strconv.ParseInt(s2, 10, 32)
	if err != nil {
		return "", 0, false
	}

	return s1, int(idx), true
}

// Fund ...
func Fund(liquidity float64, n int) float64 {
	return liquidity * math.Log(float64(n))
}

// Liquidity ...
func Liquidity(fund float64, n int) float64 {
	return fund / math.Log(float64(n))
}

// Exp ...
func Exp(liquidity, share float64) float64 {
	return math.Exp(share / liquidity)
}

// ExpSum ...
func ExpSum(liquidity float64, shares []float64) float64 {
	sum := 0.0

	for _, x := range shares {
		sum += Exp(liquidity, x)
	}

	return sum
}

// SharePrice ...
func SharePrice(liquidity, share, sum float64) float64 {
	return Exp(liquidity, share) / sum
}

// Cost ...
func Cost(liquidity float64, shares []float64) float64 {
	return liquidity * math.Log(ExpSum(liquidity, shares))
}

// Return ...
func Return(liquidity, answer float64, shares []float64) float64 {
	return Cost(liquidity, shares) - answer
}

// Revenue ...
func Revenue(liquidity, fund, answer float64, shares []float64) float64 {
	return Return(liquidity, answer, shares) - fund
}
