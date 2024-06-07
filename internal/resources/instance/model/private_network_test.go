package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newPrivateNetwork(t *testing.T) {
	sdkPrivateNetwork := publicCloud.NewPrivateNetwork()
	sdkPrivateNetwork.SetPrivateNetworkId("id")
	sdkPrivateNetwork.SetStatus("status")
	sdkPrivateNetwork.SetSubnet("subnet")

	privateNetwork := newPrivateNetwork(*sdkPrivateNetwork)

	assert.Equal(t, "id", privateNetwork.Id.ValueString(), "id should be set")
	assert.Equal(t, "status", privateNetwork.Status.ValueString(), "status should be set")
	assert.Equal(t, "subnet", privateNetwork.Subnet.ValueString(), "subnet should be set")
}

func TestPrivateNetwork_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		PrivateNetwork{}.attributeTypes(),
		PrivateNetwork{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
