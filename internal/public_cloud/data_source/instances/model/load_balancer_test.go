package model

import (
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newLoadBalancer(t *testing.T) {
	reference := "reference"
	state, _ := publicCloud.NewStateFromValue("RUNNING")
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")

	sdkLoadBalancerDetails := publicCloud.NewLoadBalancerDetails(
		"id",
		"type",
		publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "resources"}},
		"region",
		*publicCloud.NewNullableString(&reference),
		*state,
		publicCloud.Contract{BillingFrequency: 5},
		*publicCloud.NewNullableTime(&startedAt),
		[]publicCloud.IpDetails{{Ip: "1.2.3.4"}},
		*publicCloud.NewNullableLoadBalancerConfiguration(&publicCloud.LoadBalancerConfiguration{Balance: "balance"}),
		*publicCloud.NewNullableAutoScalingGroup(nil),
		*publicCloud.NewNullablePrivateNetwork(&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"}),
	)

	got := newLoadBalancer(*sdkLoadBalancerDetails)

	assert.Equal(t, "id", got.Id.ValueString(), "id is set")
	assert.Equal(
		t,
		"type",
		got.Type.ValueString(),
		"type is set",
	)
	assert.Equal(
		t,
		"resources",
		got.Resources.Cpu.Unit.ValueString(),
		"resources is set",
	)
	assert.Equal(
		t,
		"region",
		got.Region.ValueString(),
		"region is set",
	)
	assert.Equal(
		t,
		"reference",
		got.Reference.ValueString(),
		"reference is set",
	)
	assert.Equal(
		t,
		"RUNNING",
		got.State.ValueString(),
		"state is set",
	)
	assert.Equal(
		t,
		int64(5),
		got.Contract.BillingFrequency.ValueInt64(),
		"contract is set",
	)
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.StartedAt.ValueString(),
		"startedAt is set",
	)
	assert.Equal(
		t,
		"1.2.3.4",
		got.Ips[0].Ip.ValueString(),
		"ips is set",
	)
	assert.Equal(
		t,
		"balance",
		got.LoadBalancerConfiguration.Balance.ValueString(),
		"configuration is set",
	)
	assert.Equal(
		t,
		"privateNetworkId",
		got.PrivateNetwork.Id.ValueString(),
		"privateNetwork is set",
	)
}
