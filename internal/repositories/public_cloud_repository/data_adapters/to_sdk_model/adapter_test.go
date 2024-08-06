package to_sdk_model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/stretchr/testify/assert"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestAdaptToLaunchInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Type = domain.InstanceType{Name: "tralala"}

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.RootDiskStorageType = "tralala"

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Contract.Type = "tralala"

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Contract.Term = 55
		instance.Contract.Type = enum.ContractTypeHourly

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Contract.BillingFrequency = 55

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Type = domain.InstanceType{Name: "tralala"}

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("required values are set", func(t *testing.T) {
		instance, _ := domain.NewCreateInstance(
			"region",
			"lsw.c3.4xlarge",
			enum.RootDiskStorageTypeCentral,
			"ALMALINUX_8_64BIT",
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			domain.OptionalCreateInstanceValues{},
			[]string{"lsw.c3.4xlarge"},
		)

		got, err := AdaptToLaunchInstanceOpts(*instance)

		assert.NoError(t, err)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, publicCloud.TYPENAME_C3_4XLARGE, got.Type)
		assert.Equal(
			t,
			publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
			got.RootDiskStorageType,
		)
		assert.Equal(t, "ALMALINUX_8_64BIT", got.ImageId)
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.ContractType)
		assert.Equal(t, publicCloud.CONTRACTTERM__6, got.ContractTerm)
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__3, got.BillingFrequency)

		assert.Nil(t, got.MarketAppId)
		assert.Nil(t, got.Reference)
		assert.Nil(t, got.SshKey)
	})

	t.Run("optional values are set", func(t *testing.T) {
		marketAppId := "marketAppId"
		reference := "reference"
		sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)

		instance, _ := domain.NewCreateInstance(
			"",
			"lsw.m3.large",
			enum.RootDiskStorageTypeCentral,
			"ALMALINUX_8_64BIT",
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			domain.OptionalCreateInstanceValues{
				MarketAppId: &marketAppId,
				Reference:   &reference,
				SshKey:      sshKeyValueObject,
			},
			[]string{"lsw.m3.large"},
		)

		got, err := AdaptToLaunchInstanceOpts(*instance)

		assert.NoError(t, err)
		assert.Equal(t, marketAppId, *got.MarketAppId)
		assert.Equal(t, reference, *got.Reference)
		assert.Equal(t, defaultSshKey, *got.SshKey)
	})
}

func TestAdaptToUpdateInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			domain.Instance{Type: domain.InstanceType{Name: "tralala"}},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			domain.Instance{Contract: domain.Contract{Type: "tralala"}},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			domain.Instance{Contract: domain.Contract{Term: 55}},
		)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			domain.Instance{Contract: domain.Contract{BillingFrequency: 55}},
		)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("values are set", func(t *testing.T) {
		instanceType := "lsw.c3.large"
		reference := "reference"
		contractType := enum.ContractTypeMonthly
		contractTerm := enum.ContractTermThree
		billingFrequency := enum.ContractBillingFrequencySix
		rootDiskSize, _ := value_object.NewRootDiskSize(23)

		instance, _ := domain.NewUpdateInstance(
			value_object.NewGeneratedUuid(),
			domain.OptionalUpdateInstanceValues{
				Type:             &instanceType,
				Reference:        &reference,
				ContractType:     &contractType,
				Term:             &contractTerm,
				BillingFrequency: &billingFrequency,
				RootDiskSize:     rootDiskSize,
			},
			[]string{"lsw.c3.large"},
		)

		got, err := AdaptToUpdateInstanceOpts(*instance)

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.TYPENAME_C3_LARGE, got.GetType())
		assert.Equal(t, "reference", got.GetReference())
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.GetContractType())
		assert.Equal(t, publicCloud.CONTRACTTERM__3, got.GetContractTerm())
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__6, got.GetBillingFrequency())
		assert.Equal(t, int32(23), got.GetRootDiskSize())
	})
}

func generateDomainInstance() domain.Instance {
	rootDiskSize, _ := value_object.NewRootDiskSize(5)

	return domain.NewInstance(
		value_object.NewGeneratedUuid(),
		"region",
		domain.Resources{},
		domain.Image{},
		enum.StateCreating,
		"productType",
		false,
		true,
		*rootDiskSize,
		domain.InstanceType{Name: "lsw.m3.xlarge"},
		enum.RootDiskStorageTypeCentral,
		domain.Ips{},
		domain.Contract{Type: enum.ContractTypeMonthly},
		domain.OptionalInstanceValues{},
	)
}
