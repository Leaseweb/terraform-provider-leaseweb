package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestNewInstance(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()
		rootDiskSize, _ := value_object.NewRootDiskSize(5)

		got := NewInstance(
			instanceId,
			"region",
			Resources{Cpu: Cpu{Unit: "cpu"}},
			Image{Name: "image"},
			enum.StateRunning,
			"productType",
			false,
			true,
			*rootDiskSize,
			"instanceType",
			enum.RootDiskStorageTypeCentral,
			Ips{{Ip: "1.2.3.4"}},
			Contract{BillingFrequency: enum.ContractBillingFrequencyOne},
			OptionalInstanceValues{},
		)

		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "cpu", got.Resources.Cpu.Unit)
		assert.Equal(t, "image", got.Image.Name)
		assert.Equal(t, enum.StateRunning, got.State)
		assert.Equal(t, "productType", got.ProductType)
		assert.False(t, got.HasPublicIpv4)
		assert.True(t, got.HasPrivateNetwork)
		assert.Equal(t, "instanceType", got.Type)
		assert.Equal(t, enum.RootDiskStorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, int64(5), got.RootDiskSize.Value)

		assert.Nil(t, got.Reference)
		assert.Nil(t, got.Iso)
		assert.Nil(t, got.MarketAppId)
		assert.Nil(t, got.SshKey)
		assert.Nil(t, got.StartedAt)
		assert.Nil(t, got.PrivateNetwork)
		assert.Nil(t, got.AutoScalingGroup)
	})

	t.Run("optional values are set", func(t *testing.T) {
		reference := "Reference"
		marketAppId := "marketAppId"
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)
		startedAt := time.Now()

		got := NewInstance(
			value_object.NewGeneratedUuid(),
			"",
			Resources{},
			Image{},
			enum.StateRunning,
			"",
			false,
			true,
			value_object.RootDiskSize{},
			"",
			enum.RootDiskStorageTypeCentral,
			Ips{},
			Contract{},
			OptionalInstanceValues{
				Reference:        &reference,
				MarketAppId:      &marketAppId,
				SshKey:           sshKeyValueObject,
				Iso:              &Iso{Id: "isoId"},
				StartedAt:        &startedAt,
				PrivateNetwork:   &PrivateNetwork{Id: "privateNetworkId"},
				AutoScalingGroup: &AutoScalingGroup{Region: "autoScalingGroupRegion"},
			},
		)

		assert.Equal(t, "Reference", *got.Reference)
		assert.Equal(t, Iso{Id: "isoId"}, *got.Iso)
		assert.Equal(
			t,
			PrivateNetwork{Id: "privateNetworkId"},
			*got.PrivateNetwork,
		)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(
			t,
			sshKey,
			got.SshKey.String(),
		)
		assert.Equal(t, startedAt, *got.StartedAt)
		assert.Equal(t, "autoScalingGroupRegion", got.AutoScalingGroup.Region)
	})
}

func TestNewCreateInstance(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		got := NewCreateInstance(
			"region",
			"instanceType",
			enum.RootDiskStorageTypeCentral,
			enum.Almalinux864Bit,
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			OptionalCreateInstanceValues{},
		)

		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "instanceType", got.Type)
		assert.Equal(t, enum.RootDiskStorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, enum.Almalinux864Bit, got.Image.Id)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermSix, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyThree,
			got.Contract.BillingFrequency,
		)

		assert.Nil(t, got.MarketAppId)
		assert.Nil(t, got.Reference)
		assert.Nil(t, got.SshKey)
	})

	t.Run("optional values are set", func(t *testing.T) {
		marketAppId := "marketAppId"
		reference := "reference"
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		got := NewCreateInstance(
			"",
			"",
			enum.RootDiskStorageTypeCentral,
			enum.Almalinux864Bit,
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			OptionalCreateInstanceValues{
				MarketAppId: &marketAppId,
				Reference:   &reference,
				SshKey:      sshKeyValueObject,
			},
		)

		assert.Equal(t, marketAppId, *got.MarketAppId)
		assert.Equal(t, reference, *got.Reference)
		assert.Equal(t, sshKeyValueObject, got.SshKey)
	})

}

func TestNewUpdateInstance(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()

		got := NewUpdateInstance(id, OptionalUpdateInstanceValues{})

		assert.Equal(t, id, got.Id)

		assert.Empty(t, got.Type)
		assert.Empty(t, got.Reference)
		assert.Empty(t, got.Contract.Type)
		assert.Empty(t, got.Contract.Term)
		assert.Empty(t, got.Contract.BillingFrequency)
		assert.Empty(t, got.RootDiskSize)
	})

	t.Run("optional values are set", func(t *testing.T) {
		instanceType := "lsw.m4.large"
		reference := "reference"
		contractType := enum.ContractTypeMonthly
		contractTerm := enum.ContractTermSix
		billingFrequency := enum.ContractBillingFrequencyThree
		rootDiskSize, _ := value_object.NewRootDiskSize(50)

		got := NewUpdateInstance(
			value_object.NewGeneratedUuid(),
			OptionalUpdateInstanceValues{
				Type:             &instanceType,
				Reference:        &reference,
				ContractType:     &contractType,
				Term:             &contractTerm,
				BillingFrequency: &billingFrequency,
				RootDiskSize:     rootDiskSize,
			},
		)

		assert.Equal(t, "lsw.m4.large", got.Type)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermSix, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyThree,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, int64(50), got.RootDiskSize.Value)
	})

}
