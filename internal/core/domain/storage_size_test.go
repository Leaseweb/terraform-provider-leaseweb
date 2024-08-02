package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageSize(t *testing.T) {
	got := NewStorageSize(1, "unit")
	want := StorageSize{Size: 1, Unit: "unit"}

	assert.Equal(t, want, got)
}
