package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newNetworkSpeed(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(23, "unit")

	got := newNetworkSpeed(networkSpeed)

	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
	assert.Equal(
		t,
		int64(23),
		got.Value.ValueInt64(),
		"value should be set",
	)
}
