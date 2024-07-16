package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newMemory(t *testing.T) {
	entityMemory := domain.NewMemory(1, "unit")

	got := newMemory(entityMemory)

	assert.Equal(
		t,
		float64(1),
		got.Value.ValueFloat64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
}
