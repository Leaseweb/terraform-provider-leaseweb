package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newCpu(t *testing.T) {
	cpu := domain.NewCpu(1, "unit")
	got := newCpu(cpu)

	assert.Equal(t, int64(1), got.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}
