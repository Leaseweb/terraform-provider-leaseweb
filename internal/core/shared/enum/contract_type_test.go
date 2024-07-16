package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractType_String(t *testing.T) {
	got := ContractTypeHourly.String()

	assert.Equal(t, "HOURLY", got)
}

func TestNewContractType(t *testing.T) {
	want := ContractTypeMonthly
	got, err := NewContractType("MONTHLY")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
