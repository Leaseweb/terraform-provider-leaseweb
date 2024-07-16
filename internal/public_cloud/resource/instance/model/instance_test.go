package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestInstance_Populate(t *testing.T) {
	t.Run("instance is populated", func(t *testing.T) {
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		marketAppId := "marketAppId"
		reference := "reference"
		id := value_object.NewGeneratedUuid()
		rootDiskSize, _ := value_object.NewRootDiskSize(32)
		autoScalingGroupId := value_object.NewGeneratedUuid()
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		instance := domain.NewInstance(
			id,
			"region",
			domain.Resources{Cpu: domain.Cpu{Unit: "cpu"}},
			domain.Image{Id: enum.Ubuntu200464Bit},
			enum.StateCreating,
			"productType",
			false,
			true,
			*rootDiskSize,
			"lsw.m5a.4xlarge",
			enum.RootDiskStorageTypeCentral,
			domain.Ips{{Ip: "1.2.3.4"}},
			domain.Contract{Type: enum.ContractTypeMonthly},
			domain.OptionalInstanceValues{
				Reference:        &reference,
				Iso:              &domain.Iso{Id: "isoId"},
				MarketAppId:      &marketAppId,
				SshKey:           sshKeyValueObject,
				StartedAt:        &startedAt,
				PrivateNetwork:   &domain.PrivateNetwork{Id: "privateNetworkId"},
				AutoScalingGroup: &domain.AutoScalingGroup{Id: autoScalingGroupId},
			},
		)

		got := Instance{}
		got.Populate(instance, context.TODO())

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
			"CREATING",
			got.State.ValueString(),
			"state should be set",
		)
		assert.Equal(
			t,
			"productType",
			got.ProductType.ValueString(),
			"productType should be set",
		)
		assert.False(
			t,
			got.HasPublicIpv4.ValueBool(),
			"hasPublicIpv should be set",
		)
		assert.True(
			t,
			got.HasPrivateNetwork.ValueBool(),
			"hasPrivateNetwork should be set",
		)
		assert.Equal(
			t,
			"lsw.m5a.4xlarge",
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
			"reference",
			got.Reference.ValueString(),
			"reference should be set",
		)

		image := Image{}
		got.Image.As(
			context.TODO(),
			&image,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"UBUNTU_20_04_64BIT",
			image.Id.ValueString(),
			"image should be set",
		)

		contract := Contract{}
		got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"MONTHLY",
			contract.Type.ValueString(),
			"contract should be set",
		)

		iso := Iso{}
		got.Iso.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"isoId",
			iso.Id.ValueString(),
			"iso should be set",
		)

		privateNetwork := PrivateNetwork{}
		got.PrivateNetwork.As(
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
		got.AutoScalingGroup.As(
			context.TODO(),
			&autoScalingGroup,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			autoScalingGroupId.String(),
			autoScalingGroup.Id.ValueString(),
			"autoScalingGroup should be set",
		)

		var ips []Ip
		got.Ips.ElementsAs(context.TODO(), &ips, false)
		assert.Len(t, ips, 1)
		assert.Equal(
			t,
			"1.2.3.4",
			ips[0].Ip.ValueString(),
			"ip should be set",
		)

		resources := Resources{}
		cpu := Cpu{}
		got.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
		resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"cpu",
			cpu.Unit.ValueString(),
			"privateNetwork should be set",
		)

		assert.Equal(t, sshKey, got.SshKey.ValueString())
	})
}

func TestInstance_GenerateCreateInstanceEntity(t *testing.T) {
	t.Run("required values are passed", func(t *testing.T) {
		instanceEntity := domain.Instance{
			Region:              "region",
			Type:                "lsw.m5a.4xlarge",
			RootDiskStorageType: enum.RootDiskStorageTypeCentral,
			Image:               domain.Image{Id: enum.Ubuntu200464Bit},
			Contract: domain.Contract{
				Type:             enum.ContractTypeMonthly,
				Term:             enum.ContractTermThree,
				BillingFrequency: enum.ContractBillingFrequencyThree,
			},
		}

		instance := Instance{}
		instance.Populate(instanceEntity, context.TODO())

		got, diags := instance.GenerateCreateInstanceEntity(context.TODO())

		assert.False(t, diags.HasError())
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "lsw.m5a.4xlarge", got.Type)
		assert.Equal(t, enum.RootDiskStorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, enum.Ubuntu200464Bit, got.Image.Id)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermThree, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyThree,
			got.Contract.BillingFrequency,
		)
		assert.Nil(t, got.MarketAppId)
		assert.Nil(t, got.Reference)
		assert.Equal(t, value_object.RootDiskSize{}, got.RootDiskSize)
	})

	t.Run("optional values are passed", func(t *testing.T) {
		marketAppId := "marketAppId"
		reference := "reference"
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		instanceEntity := domain.Instance{
			Region:              "region",
			Type:                "lsw.m5a.4xlarge",
			RootDiskStorageType: enum.RootDiskStorageTypeCentral,
			RootDiskSize:        value_object.RootDiskSize{Value: 55},
			Image:               domain.Image{Id: enum.Ubuntu200464Bit},
			Contract: domain.Contract{
				Type:             enum.ContractTypeMonthly,
				Term:             enum.ContractTermThree,
				BillingFrequency: enum.ContractBillingFrequencyThree,
			},
			MarketAppId: &marketAppId,
			Reference:   &reference,
			SshKey:      sshKeyValueObject,
		}

		instance := Instance{}
		instance.Populate(instanceEntity, context.TODO())

		got, diags := instance.GenerateCreateInstanceEntity(context.TODO())

		assert.False(t, diags.HasError())
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, 55, got.RootDiskSize.Value)
		assert.Equal(t, sshKey, got.SshKey.String())
	})
}

func TestInstance_GenerateUpdateInstanceEntity(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()

		instanceEntity := domain.Instance{Id: id}

		instance := Instance{}
		instance.Populate(instanceEntity, context.TODO())

		got, diags := instance.GenerateUpdateInstanceEntity(context.TODO())

		assert.False(t, diags.HasError())
		assert.Equal(t, id, got.Id)
	})

	t.Run("optional values are set", func(t *testing.T) {
		reference := "reference"
		rootDiskSize, _ := value_object.NewRootDiskSize(65)

		instanceEntity := domain.Instance{
			Type: "lsw.m5a.4xlarge",
			Contract: domain.Contract{
				Type:             enum.ContractTypeMonthly,
				Term:             enum.ContractTermThree,
				BillingFrequency: enum.ContractBillingFrequencyThree,
			},
			Reference:    &reference,
			RootDiskSize: *rootDiskSize,
		}

		instance := Instance{}
		instance.Populate(instanceEntity, context.TODO())

		got, diags := instance.GenerateUpdateInstanceEntity(context.TODO())

		assert.False(t, diags.HasError())
		assert.Equal(t, "lsw.m5a.4xlarge", got.Type)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermThree, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyThree,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, 65, got.RootDiskSize.Value)
	})
}
