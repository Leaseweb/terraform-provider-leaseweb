package to_sdk_model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/stretchr/testify/assert"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestAdaptToLaunchInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Type = "tralala"

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
		instance.Type = "tralala"

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("required values are set", func(t *testing.T) {
		instance, _ := public_cloud.NewCreateInstance(
			"region",
			"lsw.c3.4xlarge",
			enum.StorageTypeCentral,
			"ALMALINUX_8_64BIT",
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			public_cloud.OptionalCreateInstanceValues{},
			[]string{"lsw.c3.4xlarge"},
		)

		got, err := AdaptToLaunchInstanceOpts(*instance)

		assert.NoError(t, err)
		assert.Equal(t, "region", string(got.Region))
		assert.Equal(t, publicCloud.TYPENAME_C3_4XLARGE, got.Type)
		assert.Equal(
			t,
			publicCloud.STORAGETYPE_CENTRAL,
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

		instance, _ := public_cloud.NewCreateInstance(
			"",
			"lsw.m3.large",
			enum.StorageTypeCentral,
			"ALMALINUX_8_64BIT",
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			public_cloud.OptionalCreateInstanceValues{
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
			public_cloud.Instance{Type: "tralala"},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			public_cloud.Instance{Contract: public_cloud.Contract{Type: "tralala"}},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			public_cloud.Instance{Contract: public_cloud.Contract{Term: 55}},
		)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			public_cloud.Instance{Contract: public_cloud.Contract{BillingFrequency: 55}},
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

		instance, _ := public_cloud.NewUpdateInstance(
			"",
			public_cloud.OptionalUpdateInstanceValues{
				Type:             &instanceType,
				Reference:        &reference,
				ContractType:     &contractType,
				Term:             &contractTerm,
				BillingFrequency: &billingFrequency,
				RootDiskSize:     rootDiskSize,
			},
			[]string{},
			"lsw.c3.large",
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

func generateDomainInstance() public_cloud.Instance {
	rootDiskSize, _ := value_object.NewRootDiskSize(5)

	return public_cloud.NewInstance(
		"",
		"region",
		public_cloud.Image{},
		enum.StateCreating,
		*rootDiskSize,
		"lsw.m3.xlarge",
		enum.StorageTypeCentral,
		public_cloud.Ips{},
		public_cloud.Contract{Type: enum.ContractTypeMonthly},
		public_cloud.OptionalInstanceValues{},
	)
}
