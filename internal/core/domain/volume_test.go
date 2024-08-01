package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVolume(t *testing.T) {
	got := NewVolume(1.2, "unit")
	want := Volume{Size: 1.2, Unit: "unit"}

	assert.Equal(t, want, got)
}
