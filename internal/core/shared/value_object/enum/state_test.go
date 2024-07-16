package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_String(t *testing.T) {
	got := StateCreating.String()

	assert.Equal(t, "CREATING", got)
}

func TestNewState(t *testing.T) {
	want := StateCreating
	got, err := NewState("CREATING")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
