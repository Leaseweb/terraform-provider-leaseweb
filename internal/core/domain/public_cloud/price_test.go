package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrice(t *testing.T) {
	got := NewPrice("hourlyPrice", "monthlyPrice")
	want := Price{HourlyPrice: "hourlyPrice", MonthlyPrice: "monthlyPrice"}

	assert.Equal(t, want, got)
}
