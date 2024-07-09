package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_String(t *testing.T) {
	got := StateCreating.String()

	assert.Equal(t, "CREATING", got)
}
