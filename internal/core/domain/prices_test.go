package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrices(t *testing.T) {
	got := NewPrices(
		"currency",
		"symbol",
		Price{HourlyPrice: "1"},
		Storage{Central: Price{MonthlyPrice: "2"}},
	)

	want := Prices{
		Currency:       "currency",
		CurrencySymbol: "symbol",
		Compute:        Price{HourlyPrice: "1"},
		Storage:        Storage{Central: Price{MonthlyPrice: "2"}},
	}

	assert.Equal(t, want, got)
}
