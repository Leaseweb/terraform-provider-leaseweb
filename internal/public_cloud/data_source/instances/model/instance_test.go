package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_newInstance(t *testing.T) {
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	marketAppId := "marketAppId"
	reference := "reference"
	id := value_object.NewGeneratedUuid()
	sshKeyValueObject, _ := value_object.NewSshKey(sshKey)
	autoScalingGroupId := value_object.NewGeneratedUuid()
	loadBalancerId := value_object.NewGeneratedUuid()

	instance := entity.NewInstance(
		id,
		"region",
		entity.Resources{Cpu: entity.Cpu{Unit: "cpu"}},
		entity.Image{Id: enum.Ubuntu200464Bit},
		"state",
		"productType",
		true,
		false,
		value_object.RootDiskSize{Value: 55},
		"lsw.m3.large",
		enum.RootDiskStorageTypeCentral,
		entity.Ips{{Ip: "1.2.3.4"}},
		entity.Contract{Type: enum.ContractTypeMonthly},
		entity.OptionalInstanceValues{
			Reference:      &reference,
			Iso:            &entity.Iso{Id: "isoId"},
			MarketAppId:    &marketAppId,
			SshKey:         sshKeyValueObject,
			StartedAt:      &startedAt,
			PrivateNetwork: &entity.PrivateNetwork{Id: "privateNetworkId"},
			AutoScalingGroup: &entity.AutoScalingGroup{
				Id:           autoScalingGroupId,
				LoadBalancer: &entity.LoadBalancer{Id: loadBalancerId},
			},
		},
	)

	got := newInstance(instance)

	assert.Equal(
		t,
		id.String(),
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
		int64(55),
		got.RootDiskSize.ValueInt64(),
		"rootDiskSize should be set",
	)
	assert.Equal(
		t,
		"CENTRAL",
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
		"UBUNTU_20_04_64BIT",
		got.Image.Id.ValueString(),
		"image should be set",
	)
	assert.Equal(
		t,
		"MONTHLY",
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
		autoScalingGroupId.String(),
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
		loadBalancerId.String(),
		got.AutoScalingGroup.LoadBalancer.Id.ValueString(),
		"loadBalancer should be set",
	)
}
