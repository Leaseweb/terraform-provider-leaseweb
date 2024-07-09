package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractType_String(t *testing.T) {
	got := ContractTypeHourly.String()

	assert.Equal(t, "HOURLY", got)
}
