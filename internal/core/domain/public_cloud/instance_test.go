package public_cloud

import (
	"testing"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/stretchr/testify/assert"
)

var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestNewInstance(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		rootDiskSize, _ := value_object.NewRootDiskSize(5)

		got := NewInstance(
			"id",
			"region",
			Resources{Cpu: Cpu{Unit: "cpu"}},
			Image{Name: "image"},
			enum.StateRunning,
			"productType",
			false,
			true,
			false,
			*rootDiskSize,
			"instanceType",
			enum.StorageTypeCentral,
			Ips{{Ip: "1.2.3.4"}},
			Contract{BillingFrequency: enum.ContractBillingFrequencyOne},
			OptionalInstanceValues{},
		)

		assert.Equal(t, "id", got.Id)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "cpu", got.Resources.Cpu.Unit)
		assert.Equal(t, "image", got.Image.Name)
		assert.Equal(t, enum.StateRunning, got.State)
		assert.Equal(t, "productType", got.ProductType)
		assert.False(t, got.HasPublicIpv4)
		assert.True(t, got.HasPrivateNetwork)
		assert.False(t, got.HasUserData)
		assert.Equal(t, "instanceType", got.Type)
		assert.Equal(t, enum.StorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, 5, got.RootDiskSize.Value)

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
			"",
			"",
			Resources{},
			Image{},
			enum.StateRunning,
			"",
			false,
			true,
			false,
			value_object.RootDiskSize{},
			"",
			enum.StorageTypeCentral,
			Ips{},
			Contract{},
			OptionalInstanceValues{
				Reference:      &reference,
				MarketAppId:    &marketAppId,
				SshKey:         sshKeyValueObject,
				Iso:            &Iso{Id: "isoId"},
				StartedAt:      &startedAt,
				PrivateNetwork: &PrivateNetwork{Id: "privateNetworkId"},
				AutoScalingGroup: &AutoScalingGroup{
					Region: "autoScalingGroupRegion",
				},
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
		got, err := NewCreateInstance(
			"region",
			"instanceType",
			enum.StorageTypeCentral,
			"ALMALINUX_8_64BIT",
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			OptionalCreateInstanceValues{},
			[]string{"instanceType"},
		)

		assert.NoError(t, err)

		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "instanceType", got.Type)
		assert.Equal(t, enum.StorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, "ALMALINUX_8_64BIT", got.Image.Id)
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
		assert.Equal(t, value_object.RootDiskSize{}, got.RootDiskSize)
	})

	t.Run("optional values are set", func(t *testing.T) {
		marketAppId := "marketAppId"
		reference := "reference"
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)
		rootDiskSize, _ := value_object.NewRootDiskSize(6)

		got, err := NewCreateInstance(
			"",
			"instanceType",
			enum.StorageTypeCentral,
			"ALMALINUX_8_64BIT",
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			OptionalCreateInstanceValues{
				MarketAppId:  &marketAppId,
				Reference:    &reference,
				SshKey:       sshKeyValueObject,
				RootDiskSize: rootDiskSize,
			},
			[]string{"instanceType"},
		)

		assert.NoError(t, err)

		assert.Equal(t, marketAppId, *got.MarketAppId)
		assert.Equal(t, reference, *got.Reference)
		assert.Equal(t, sshKeyValueObject, got.SshKey)
		assert.Equal(t, *rootDiskSize, got.RootDiskSize)
	})

	t.Run(
		"passing an invalid instance type returns an error",
		func(t *testing.T) {
			marketAppId := "marketAppId"
			reference := "reference"
			sshKeyValueObject, _ := value_object.NewSshKey(sshKey)
			rootDiskSize, _ := value_object.NewRootDiskSize(6)

			_, err := NewCreateInstance(
				"",
				"instanceType",
				enum.StorageTypeCentral,
				"ALMALINUX_8_64BIT",
				enum.ContractTypeMonthly,
				enum.ContractTermSix,
				enum.ContractBillingFrequencyThree,
				OptionalCreateInstanceValues{
					MarketAppId:  &marketAppId,
					Reference:    &reference,
					SshKey:       sshKeyValueObject,
					RootDiskSize: rootDiskSize,
				},
				[]string{},
			)

			assert.Error(t, err)
			assert.Error(t, err, ErrInvalidInstanceTypePassed{})
		},
	)
}

func TestNewUpdateInstance(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		got, err := NewUpdateInstance(
			"id",
			OptionalUpdateInstanceValues{},
			[]string{},
			"",
		)

		assert.NoError(t, err)

		assert.Equal(t, "id", got.Id)
		assert.Empty(t, got.Type)
		assert.Empty(t, got.Reference)
		assert.Empty(t, got.Contract.Type)
		assert.Empty(t, got.Contract.Term)
		assert.Empty(t, got.Contract.BillingFrequency)
		assert.Empty(t, got.RootDiskSize)
	})

	t.Run("optional values are set", func(t *testing.T) {
		instanceType := "instanceType"
		reference := "reference"
		contractType := enum.ContractTypeMonthly
		contractTerm := enum.ContractTermSix
		billingFrequency := enum.ContractBillingFrequencyThree
		rootDiskSize, _ := value_object.NewRootDiskSize(50)

		got, err := NewUpdateInstance(
			"",
			OptionalUpdateInstanceValues{
				Type:             &instanceType,
				Reference:        &reference,
				ContractType:     &contractType,
				Term:             &contractTerm,
				BillingFrequency: &billingFrequency,
				RootDiskSize:     rootDiskSize,
			},
			[]string{},
			"instanceType",
		)

		assert.NoError(t, err)

		assert.Equal(t, instanceType, got.Type)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermSix, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyThree,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, 50, got.RootDiskSize.Value)
	})

	t.Run(
		"passing an invalid instance type returns an error",
		func(t *testing.T) {
			instanceType := "instanceType"

			_, err := NewUpdateInstance(
				"",
				OptionalUpdateInstanceValues{Type: &instanceType},
				[]string{},
				"",
			)

			assert.Error(t, err)
			assert.Error(t, err, ErrInvalidInstanceTypePassed{})
		},
	)

	t.Run("currentInstanceType is respected", func(t *testing.T) {
		instanceType := "instanceType"

		_, err := NewUpdateInstance(
			"",
			OptionalUpdateInstanceValues{Type: &instanceType},
			[]string{},
			"instanceType",
		)

		assert.NoError(t, err)
	})

	t.Run("allowedInstanceTypes is respected", func(t *testing.T) {
		instanceType := "instanceType"

		_, err := NewUpdateInstance(
			"",
			OptionalUpdateInstanceValues{Type: &instanceType},
			[]string{"instanceType"},
			"",
		)

		assert.NoError(t, err)
	})
}

func TestInstance_CanBeTerminated(t *testing.T) {
	t.Run(
		"Instance cannot be terminated when state is creating",
		func(t *testing.T) {
			instance := Instance{State: enum.StateCreating}
			permission, reason := instance.CanBeTerminated()

			assert.False(t, permission)
			assert.Contains(t, *reason, enum.StateCreating.String())
		},
	)

	t.Run(
		"Instance cannot be terminated when state is destroying",
		func(t *testing.T) {
			instance := Instance{State: enum.StateDestroying}
			permission, reason := instance.CanBeTerminated()

			assert.False(t, permission)
			assert.Contains(t, *reason, enum.StateDestroying.String())
		},
	)

	t.Run(
		"Instance cannot be terminated when state is destroyed",
		func(t *testing.T) {
			instance := Instance{State: enum.StateDestroyed}
			permission, reason := instance.CanBeTerminated()

			assert.False(t, permission)
			assert.Contains(t, *reason, enum.StateDestroyed.String())
		},
	)

	t.Run(
		"Instance cannot be terminated when contract.EndsAt is not nil",
		func(t *testing.T) {
			endsAt := time.Now()
			instance := Instance{Contract: Contract{EndsAt: &endsAt}}
			permission, reason := instance.CanBeTerminated()

			assert.False(t, permission)
			assert.Contains(t, *reason, endsAt.String())
		},
	)

	t.Run(
		"Instance can be terminated in other scenarios",
		func(t *testing.T) {
			instance := Instance{}
			permission, reason := instance.CanBeTerminated()

			assert.True(t, permission)
			assert.Nil(t, reason)
		},
	)
}
