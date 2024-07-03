package model

import (
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newInstance(t *testing.T) {
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	sdkInstanceTypeName, _ := publicCloud.NewInstanceTypeNameFromValue("lsw.m3.large")
	marketAppId := "marketAppId"
	reference := "reference"
	iso := publicCloud.Iso{Id: "isoId"}
	privateNetwork := publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"}

	sdkInstanceDetails := publicCloud.NewInstanceDetails(
		"id",
		*sdkInstanceTypeName,
		publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "cpu"}},
		"region",
		*publicCloud.NewNullableString(&reference),
		*publicCloud.NewNullableTime(&startedAt),
		*publicCloud.NewNullableString(&marketAppId),
		"state",
		"productType",
		true,
		false,
		32,
		"rootDiskStorageType",
		publicCloud.Contract{Type: "contract"},
		*publicCloud.NewNullableIso(&iso),
		*publicCloud.NewNullablePrivateNetwork(&privateNetwork),
		publicCloud.OperatingSystemDetails{Id: "operatingSystemId"},
		[]publicCloud.IpDetails{{Ip: "1.2.3.4"}},
		*publicCloud.NewNullableAutoScalingGroup(nil),
	)

	sdkAutoScalingGroupDetails := publicCloud.AutoScalingGroupDetails{Id: "autoScalingGroup"}
	sdkLoadBalancerDetails := publicCloud.LoadBalancerDetails{Id: "loadBalancerId"}

	got := newInstance(
		*sdkInstanceDetails,
		&sdkAutoScalingGroupDetails,
		&sdkLoadBalancerDetails,
	)

	assert.Equal(
		t,
		"id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"region",
		got.Region.ValueString(),
		"region should be set",
	)
	assert.Equal(
		t,
		"state",
		got.State.ValueString(),
		"state should be set",
	)
	assert.Equal(
		t,
		"productType",
		got.ProductType.ValueString(),
		"productType should be set",
	)
	assert.Equal(
		t,
		true,
		got.HasPublicIpv4.ValueBool(),
		"hasPublicIpv should be set",
	)
	assert.Equal(
		t,
		false,
		got.HasPrivateNetwork.ValueBool(),
		"hasPrivateNetwork should be set",
	)
	assert.Equal(
		t,
		"lsw.m3.large",
		got.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		int64(32),
		got.RootDiskSize.ValueInt64(),
		"rootDiskSize should be set",
	)
	assert.Equal(
		t,
		"rootDiskStorageType",
		got.RootDiskStorageType.ValueString(),
		"rootDiskStorageType should be set",
	)
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.StartedAt.ValueString(),
		"startedAt should be set",
	)
	assert.Equal(
		t,
		"marketAppId",
		got.MarketAppId.ValueString(),
		"marketAppId should be set",
	)
	assert.Equal(
		t,
		"operatingSystemId",
		got.OperatingSystem.Id.ValueString(),
		"operating_system should be set",
	)
	assert.Equal(
		t,
		"contract",
		got.Contract.Type.ValueString(),
		"contract should be set",
	)
	assert.Equal(
		t,
		"1.2.3.4",
		got.Ips[0].Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"cpu",
		got.Resources.Cpu.Unit.ValueString(),
		"privateNetwork should be set",
	)
	assert.Equal(
		t,
		"autoScalingGroup",
		got.AutoScalingGroup.Id.ValueString(),
		"autoScalingGroup should be set",
	)
	assert.Equal(
		t,
		"isoId",
		got.Iso.Id.ValueString(),
		"iso should be set",
	)
	assert.Equal(
		t,
		"privateNetworkId",
		got.PrivateNetwork.Id.ValueString(),
		"privateNetwork should be set",
	)
	assert.Equal(
		t,
		"loadBalancerId",
		got.AutoScalingGroup.LoadBalancer.Id.ValueString(),
		"loadBalancer should be set",
	)
}

func Test_generateAutoScalingGroup(t *testing.T) {
	t.Run("sdkAutoScalingGroupDetails is empty", func(t *testing.T) {
		got := generateAutoScalingGroup(nil, nil)
		assert.Nil(t, got)
	})

	t.Run("sdkAutoScalingGroupDetails is set", func(t *testing.T) {
		got := generateAutoScalingGroup(&publicCloud.AutoScalingGroupDetails{}, nil)
		assert.NotNil(t, got)
	})
}
