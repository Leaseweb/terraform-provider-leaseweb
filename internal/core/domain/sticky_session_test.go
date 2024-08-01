package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStickySession(t *testing.T) {
	got := NewStickySession(true, 5)

	assert.True(t, got.Enabled)
	assert.Equal(t, 5, got.MaxLifeTime)
}
