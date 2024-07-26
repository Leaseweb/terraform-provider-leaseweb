package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	got := NewStorage(Price{HourlyPrice: "1"}, Price{MonthlyPrice: "2"})
	want := Storage{
		Local:   Price{HourlyPrice: "1"},
		Central: Price{MonthlyPrice: "2"},
	}

	assert.Equal(t, want, got)

}
