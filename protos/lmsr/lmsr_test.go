package lmsr

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/uuid"
)

func TestSplit(t *testing.T) {
	var (
		k1 = "A"
		k2 = "B"
		k3 = "C"
		k4 = "D"
	)

	key1 := fmt.Sprintf("%s%s%s%s%s%s%s", k1, Sep, k2, Sep, k3, Sep, k4)
	p1 := fmt.Sprintf("%s%s%s%s%s", k1, Sep, k2, Sep, k3)

	a, b, ok := Split(key1, Sep)
	if !ok {
		t.Fatal("split failed")
	}

	if a != p1 || b != k4 {
		t.Fatalf("data not match, %s, %s, %s, %s", a, b, p1, k4)
	}

	key2 := fmt.Sprintf("%s%s%s", k1, Sep, k2)

	a, b, ok = Split(key2, Sep)
	if !ok {
		t.Fatal("split failed")
	}

	if a != k1 || b != k2 {
		t.Fatalf("data not match, %s, %s, %s, %s", a, b, k1, k2)
	}

	if _, _, ok = Split(k1, Sep); ok {
		t.Fatalf("can not split with wrong data")
	}
}

func TestId(t *testing.T) {
	user := uuid.New().String()
	event := EventID()
	outcome := OutcomeID(event, 1)

	a, idx, ok := SepOutcomeID(outcome)
	if !ok || a != event || idx != 1 {
		t.Fatalf("data not match %s, %s, %d", outcome, a, idx)
	}

	market := MarketID(event)

	a, b, ok := SepMarketID(market)
	if !ok || a != event {
		t.Fatalf("data not match %s, %s", market, a, b)
	}

	share := ShareID(market, outcome)

	a, b, ok = SepShareID(share)
	if !ok || a != market || b != outcome {
		t.Fatalf("data not match %s, %s, %s", share, a, b)
	}

	asset := AssetID(user, share)

	a, b, ok = Split(asset, Sep)
	if !ok || a != share || b != user {
		t.Fatalf("data not match %s, %s, %s", asset, a, b)
	}

	w := "aaa@a123"
	a, idx, ok = SepOutcomeID(w)
	if ok {
		t.Fatal("can not seperate wrong format: %s", w)
	}

	a, idx, ok = SepOutcomeID(market)
	if ok {
		t.Fatal("can not seperate wrong format: %s", market)
	}

}

func TestLMSR(t *testing.T) {
	liq := 100.0
	data := []float64{0.0, 0.0, 0.0}
	sum := ExpSum(liq, data)
	if sum != float64(len(data)) {
		t.Fatalf("ExpSum failed: %v, %v", sum, float64(len(data)))
	}

	p1 := SharePrice(liq, data[0], sum)
	if p1 != (1.0 / sum) {
		t.Fatalf("share price failed: %v, %v", p1, 1.0/sum)
	}

	cost := Cost(liq, data)
	if cost != liq*math.Log(sum) {
		t.Fatalf("cost failed; %v, %v", cost, liq*math.Log(sum))
	}

}
