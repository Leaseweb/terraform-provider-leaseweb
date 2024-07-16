package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractBillingFrequency_Value(t *testing.T) {
	got := ContractBillingFrequencySix.Value()

	assert.Equal(t, 6, got)
}

func TestNewContractBillingFrequency(t *testing.T) {
	want := ContractBillingFrequencySix
	got, err := NewContractBillingFrequency(6)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
