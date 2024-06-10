package model

import (
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
