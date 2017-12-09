package prediction

import (
	fmt "fmt"

	"github.com/google/uuid"
)

const (
	// Sep ...
	Sep = "#"

	// MarketKey ...
	MarketKey = "market"

	// AssetKey ...
	AssetKey = "asset"
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

// OutcomeID ...
func OutcomeID(evt string, index int) string {
	return combine("@", evt, index)
}
