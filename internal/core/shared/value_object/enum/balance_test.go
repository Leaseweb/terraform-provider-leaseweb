package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance_String(t *testing.T) {
	got := BalanceRoundRobin.String()

	assert.Equal(t, "roundrobin", got)

}

func TestNewBalance(t *testing.T) {
	want := BalanceSource
	got, err := NewBalance("source")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
