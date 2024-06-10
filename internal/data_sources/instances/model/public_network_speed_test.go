package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newPublicNetworkSpeed(t *testing.T) {
	sdkPublicNetworkSpeed := publicCloud.NewPublicNetworkSpeed()
	sdkPublicNetworkSpeed.SetUnit("unit")
	sdkPublicNetworkSpeed.SetValue(23)

	publicNetworkSpeed := newPublicNetworkSpeed(*sdkPublicNetworkSpeed)

	assert.Equal(t, "unit", publicNetworkSpeed.Unit.ValueString(), "unit should be set")
	assert.Equal(t, int64(23), publicNetworkSpeed.Value.ValueInt64(), "value should be set")
}
