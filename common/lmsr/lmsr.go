package lmsr

import "math"

// Exp ...
func Exp(liquidity, volume float64) float64 {
	return math.Exp(volume / liquidity)
}

// Sum ...
func Sum(liquidity float64, volumes map[string]float64) float64 {
	sum := 0.0

	for _, x := range volumes {
		sum += Exp(liquidity, x)
	}

	return sum
}

// Cost ...
func Cost(liquidity float64, volumes map[string]float64) float64 {
	return liquidity * math.Log(Sum(liquidity, volumes))
}

// Fund ...
func Fund(liquidity float64, n int) float64 {
	return liquidity * math.Log(float64(n))
}

// FundWithVolumes ...
func FundWithVolumes(liquidity float64, volumes map[string]float64) float64 {
	return Cost(liquidity, volumes)
}

// Liquidity ...
func Liquidity(fund float64, n int) float64 {
	return fund / math.Log(float64(n))
}

/*
func LiquidityWithVolumes(fund float64, volumes []float64) float64 {
	return fund / math.Log(Sum())
}
*/

// Price ...
func Price(liquidity, volume, sum float64) float64 {
	return Exp(liquidity, volume) / sum
}

// Return ...
func Return(liquidity, answer float64, volumes map[string]float64) float64 {
	return Cost(liquidity, volumes) - answer
}

// Revenue ...
func Revenue(liquidity, fund, answer float64, volumes map[string]float64) float64 {
	return Return(liquidity, answer, volumes) - fund
}

// BuyWithLimitAmount ...
func BuyWithLimitAmount(liquidity, amount, old float64, idx int, volumes []float64) float64 {
	if idx >= len(volumes) {
		panic("index out of bound")
	}

	x := Exp(liquidity, amount+old)

	for i, v := range volumes {
		if i != idx {
			x -= Exp(liquidity, v)
		}
	}

	return liquidity*math.Log(x) - volumes[idx]
}
