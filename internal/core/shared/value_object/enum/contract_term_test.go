package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractTerm_Value(t *testing.T) {
	got := ContractTermSix.Value()

	assert.Equal(t, int64(6), got)
}
