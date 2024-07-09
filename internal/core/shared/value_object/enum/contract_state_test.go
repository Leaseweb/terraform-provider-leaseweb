package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractState_String(t *testing.T) {
	got := ContractStateActive.String()

	assert.Equal(t, "ACTIVE", got)

}
