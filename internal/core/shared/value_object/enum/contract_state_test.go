package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractState_String(t *testing.T) {
	got := ContractStateActive.String()

	assert.Equal(t, "ACTIVE", got)

}

func TestNewContractState(t *testing.T) {
	want := ContractStateDeleteScheduled
	got, err := NewContractState("DELETE_SCHEDULED")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
