package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newNetworkSpeed(t *testing.T) {
	sdkNetworkSpeed := publicCloud.NewNetworkSpeed(23, "unit")

	networkSpeed := newNetworkSpeed(sdkNetworkSpeed)

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

func TestNetworkSpeed_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		NetworkSpeed{}.AttributeTypes(),
		NetworkSpeed{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
