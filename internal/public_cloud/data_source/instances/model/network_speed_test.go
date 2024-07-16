package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func Test_newNetworkSpeed(t *testing.T) {
	entityNetworkSpeed := entity.NewNetworkSpeed(23, "unit")

	networkSpeed := newNetworkSpeed(entityNetworkSpeed)

	assert.Equal(
		t,
		"unit",
		networkSpeed.Unit.ValueString(),
		"unit should be set",
	)
	assert.Equal(
		t,
		int64(23),
		networkSpeed.Value.ValueInt64(),
		"value should be set",
	)
}
