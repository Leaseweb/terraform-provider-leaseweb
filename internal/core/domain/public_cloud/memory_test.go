package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemory(t *testing.T) {
	memory := NewMemory(1, "unit")

	assert.Equal(t, float64(1), memory.Value)
	assert.Equal(t, "unit", memory.Unit)
}
