package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCpu(t *testing.T) {
	cpu := NewCpu(1, "unit")

	assert.Equal(t, int64(1), cpu.Value)
	assert.Equal(t, "unit", cpu.Unit)
}
