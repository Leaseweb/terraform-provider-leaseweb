package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStickySession(t *testing.T) {
	got := NewStickySession(true, 5)

	assert.True(t, got.Enabled)
	assert.Equal(t, int64(5), got.MaxLifeTime)
}
