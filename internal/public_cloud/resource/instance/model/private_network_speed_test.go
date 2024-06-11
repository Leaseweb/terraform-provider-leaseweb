package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newPrivateNetworkSpeed(t *testing.T) {
	sdkPrivateNetworkSpeed := publicCloud.NewPrivateNetworkSpeed()
	sdkPrivateNetworkSpeed.SetUnit("unit")
	sdkPrivateNetworkSpeed.SetValue(23)

	privateNetworkSpeed := newPrivateNetworkSpeed(sdkPrivateNetworkSpeed)

	assert.Equal(t, "unit", privateNetworkSpeed.Unit.ValueString(), "unit should be set")
	assert.Equal(t, int64(23), privateNetworkSpeed.Value.ValueInt64(), "value should be set")
}

func TestPrivateNetworkSpeed_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		PrivateNetworkSpeed{}.attributeTypes(),
		PrivateNetworkSpeed{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
