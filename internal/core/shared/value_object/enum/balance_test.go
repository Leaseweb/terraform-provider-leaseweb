package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance_String(t *testing.T) {
	got := BalanceRoundRobin.String()

	assert.Equal(t, "ROUNDROBIN", got)

}
