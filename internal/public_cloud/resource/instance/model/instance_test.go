package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestInstance_Populate(t *testing.T) {
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	sdkInstanceTypeName, _ := publicCloud.NewInstanceTypeNameFromValue("lsw.m5a.4xlarge")
	marketAppId := "marketAppId"
	reference := "reference"
	state, _ := publicCloud.NewStateFromValue("CREATING")

	sdkInstanceDetails := publicCloud.NewInstanceDetails(
		"id",
		*sdkInstanceTypeName,
		publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "cpu"}},
		"region",
		*publicCloud.NewNullableString(&reference),
		*publicCloud.NewNullableTime(&startedAt),
		*publicCloud.NewNullableString(&marketAppId),
		*state,
		"productType",
		true,
		false,
		32,
		"rootDiskStorageType",
		publicCloud.Contract{Type: "contract"},
		*publicCloud.NewNullableIso(&publicCloud.Iso{Id: "isoId"}),
		*publicCloud.NewNullablePrivateNetwork(&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"}),
		publicCloud.ImageDetails{Id: "imageId"},
		[]publicCloud.IpDetails{{Ip: "1.2.3.4"}},
		*publicCloud.NewNullableAutoScalingGroup(nil),
	)

	sdkAutoScalingGroupDetails := publicCloud.AutoScalingGroupDetails{Id: "autoScalingGroupId"}
	sdkLoadBalancerDetails := publicCloud.LoadBalancerDetails{Id: "loadBalancerId"}

	instance := Instance{}
	instance.Populate(
		sdkInstanceDetails,
		&sdkAutoScalingGroupDetails,
		&sdkLoadBalancerDetails,
		context.TODO(),
	)

	assert.Equal(
		t,
		"id",
		instance.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"region",
		instance.Region.ValueString(),
		"region should be set",
	)
	assert.Equal(
		t,
		"CREATING",
		instance.State.ValueString(),
		"state should be set",
	)
	assert.Equal(
		t,
		"productType",
		instance.ProductType.ValueString(),
		"productType should be set",
	)
	assert.True(
		t,
		instance.HasPublicIpv4.ValueBool(),
		"hasPublicIpv should be set",
	)
	assert.False(
		t,
		instance.HasPrivateNetwork.ValueBool(),
		"hasPrivateNetwork should be set",
	)
	assert.Equal(
		t,
		"lsw.m5a.4xlarge",
		instance.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		int64(32),
		instance.RootDiskSize.ValueInt64(),
		"rootDiskSize should be set",
	)
	assert.Equal(
		t,
		"rootDiskStorageType",
		instance.RootDiskStorageType.ValueString(),
		"rootDiskStorageType should be set",
	)
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		instance.StartedAt.ValueString(),
		"startedAt should be set",
	)
	assert.Equal(
		t,
		"marketAppId",
		instance.MarketAppId.ValueString(),
		"marketAppId should be set",
	)
	assert.Equal(
		t,
		"reference",
		instance.Reference.ValueString(),
		"reference should be set",
	)

	image := Image{}
	instance.Image.As(
		context.TODO(),
		&image,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"imageId",
		image.Id.ValueString(),
		"image should be set",
	)

	contract := Contract{}
	instance.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"contract",
		contract.Type.ValueString(),
		"contract should be set",
	)

	iso := Iso{}
	instance.Iso.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"isoId",
		iso.Id.ValueString(),
		"iso should be set",
	)

	privateNetwork := PrivateNetwork{}
	instance.PrivateNetwork.As(
		context.TODO(),
		&privateNetwork,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"privateNetworkId",
		privateNetwork.Id.ValueString(),
		"privateNetwork should be set",
	)

	autoScalingGroup := AutoScalingGroup{}
	instance.AutoScalingGroup.As(
		context.TODO(),
		&autoScalingGroup,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"autoScalingGroupId",
		autoScalingGroup.Id.ValueString(),
		"autoScalingGroup should be set",
	)

	loadBalancer := LoadBalancer{}
	autoScalingGroup.LoadBalancer.As(
		context.TODO(),
		&loadBalancer,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"loadBalancerId",
		loadBalancer.Id.ValueString(),
		"loadBalancer should be set",
	)

	var ips []Ip
	instance.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Equal(
		t,
		"1.2.3.4",
		ips[0].Ip.ValueString(),
		"ip should be set",
	)

	resources := Resources{}
	cpu := Cpu{}
	instance.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
	resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"cpu",
		cpu.Unit.ValueString(),
		"privateNetwork should be set",
	)
}

func Test_generateAutoScalingGroup(t *testing.T) {
	t.Run("autoScalingGroup is not set", func(t *testing.T) {
		got, diags := generateAutoScalingGroup(
			context.TODO(),
			nil,
			nil,
		)

		assert.Nil(t, diags)
		assert.Equal(t, types.ObjectNull(AutoScalingGroup{}.AttributeTypes()), got)
	})

	t.Run("autoScalingGroup is set", func(t *testing.T) {
		got, diags := generateAutoScalingGroup(
			context.TODO(),
			&publicCloud.AutoScalingGroupDetails{Id: "autoScalingGroupId"},
			nil,
		)

		assert.Nil(t, diags)

		autoScalingGroup := AutoScalingGroup{}
		got.As(
			context.TODO(),
			&autoScalingGroup,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"autoScalingGroupId",
			autoScalingGroup.Id.ValueString(),
			"autoScalingGroup should be set",
		)
	})
}
