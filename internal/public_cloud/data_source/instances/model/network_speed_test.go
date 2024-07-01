package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newNetworkSpeed(t *testing.T) {
	sdkNetworkSpeed := publicCloud.NewNetworkSpeed(23, "unit")

	networkSpeed := newNetworkSpeed(*sdkNetworkSpeed)

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
