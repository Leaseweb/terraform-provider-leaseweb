package publiccloud

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptLoadBalancerDetailsToLoadBalancerDataSource(t *testing.T) {
	reference := "reference"

	sdkLoadBalancerDetails := publicCloud.LoadBalancerDetails{
		Id:        "id",
		Region:    "region",
		Reference: *publicCloud.NewNullableString(&reference),
		State:     publicCloud.STATE_CREATING,
		Type:      publicCloud.TYPENAME_C3_2XLARGE,
		Ips: []publicCloud.IpDetails{
			{Ip: "127.0.0.1"},
		},
		Contract: publicCloud.Contract{
			Term: publicCloud.CONTRACTTERM__1,
		},
	}

	got := adaptLoadBalancerDetailsToLoadBalancerDataSource(sdkLoadBalancerDetails)

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
	sdkLoadBalancers := []publicCloud.LoadBalancerDetails{
		{Id: "id"},
	}

	got := adaptLoadBalancersToLoadBalancersDataSource(sdkLoadBalancers)

	assert.Len(t, got.LoadBalancers, 1)
	assert.Equal(t, "id", got.LoadBalancers[0].ID.ValueString())
}
