package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInstanceType(t *testing.T) {
	want := InstanceType{Name: "name"}
	got := NewInstanceType("name")

	assert.Equal(t, want, got)
}
