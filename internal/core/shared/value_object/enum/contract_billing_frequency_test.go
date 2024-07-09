package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractBillingFrequency_Value(t *testing.T) {
	got := ContractBillingFrequencySix.Value()

	assert.Equal(t, int64(6), got)

}
