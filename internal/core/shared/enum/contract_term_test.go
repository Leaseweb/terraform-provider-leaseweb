package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractTerm_Value(t *testing.T) {
	got := ContractTermSix.Value()

	assert.Equal(t, 6, got)
}

func TestNewContractTerm(t *testing.T) {
	want := ContractTermOne
	got, err := NewContractTerm(1)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestContractTerm_Values(t *testing.T) {
	want := []int{0, 1, 3, 6, 12}
	got := ContractTermThree.Values()

	assert.Equal(t, want, got)
}
