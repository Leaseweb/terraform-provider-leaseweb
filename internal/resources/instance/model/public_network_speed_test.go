package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newPublicNetworkSpeed(t *testing.T) {
	sdkPublicNetworkSpeed := publicCloud.NewPublicNetworkSpeed()
	sdkPublicNetworkSpeed.SetUnit("unit")
	sdkPublicNetworkSpeed.SetValue(23)

	publicNetworkSpeed := newPublicNetworkSpeed(sdkPublicNetworkSpeed)

	assert.Equal(t, "unit", publicNetworkSpeed.Unit.ValueString(), "unit should be set")
	assert.Equal(t, int64(23), publicNetworkSpeed.Value.ValueInt64(), "value should be set")
}

func TestPublicNetworkSpeed_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		PublicNetworkSpeed{}.attributeTypes(),
		PublicNetworkSpeed{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
