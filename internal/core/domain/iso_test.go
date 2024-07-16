package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIso(t *testing.T) {
	iso := NewIso("id", "name")

	assert.Equal(t, "id", iso.Id)
	assert.Equal(t, "name", iso.Name)
}
