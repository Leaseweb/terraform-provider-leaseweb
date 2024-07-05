package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInstanceType(t *testing.T) {
	instanceType := NewInstanceType("name")

	assert.Equal(t, "name", instanceType.Name)
}
