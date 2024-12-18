package publiccloud

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptLoadBalancerListItemToLoadBalancerDataSource(t *testing.T) {
	reference := "reference"

	sdkLoadBalancerDetails := publiccloud.LoadBalancerListItem{
		Id:        "id",
		Region:    "region",
		Reference: *publiccloud.NewNullableString(&reference),
		State:     publiccloud.STATE_CREATING,
		Type:      publiccloud.TYPENAME_C3_2XLARGE,
		Ips: []publiccloud.Ip{
			{Ip: "127.0.0.1"},
		},
		Contract: publiccloud.Contract{
			Term: publiccloud.CONTRACTTERM__1,
		},
	}

	got := adaptLoadBalancerListItemToLoadBalancerDataSource(sdkLoadBalancerDetails)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
	assert.Len(t, got.IPs, 1)
	assert.Equal(t, "127.0.0.1", got.IPs[0].IP.ValueString())
	assert.Equal(t, int32(1), got.Contract.Term.ValueInt32())
}

func Test_adaptLoadBalancersToLoadBalancersDatasource(t *testing.T) {
	sdkLoadBalancers := []publiccloud.LoadBalancerListItem{
		{Id: "id"},
	}

	got := adaptLoadBalancersToLoadBalancersDataSource(sdkLoadBalancers)

	assert.Len(t, got.LoadBalancers, 1)
	assert.Equal(t, "id", got.LoadBalancers[0].ID.ValueString())
}
