package to_instance

import (
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestAdaptToLaunchInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Type = value_object.InstanceType{Type: "tralala"}

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
		instance.Type = value_object.InstanceType{Type: "tralala"}

		_, err := AdaptToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("required values are set", func(t *testing.T) {
		instance := domain.NewCreateInstance(
			"region",
			value_object.InstanceType{Type: string(publicCloud.TYPENAME_C3_4XLARGE)},
			enum.RootDiskStorageTypeCentral,
			enum.Almalinux864Bit,
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			domain.OptionalCreateInstanceValues{},
		)

		got, err := AdaptToLaunchInstanceOpts(instance)

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

		instance := domain.NewCreateInstance(
			"",
			value_object.NewUnvalidatedInstanceType(
				string(publicCloud.TYPENAME_M3_LARGE),
			),
			enum.RootDiskStorageTypeCentral,
			enum.Almalinux864Bit,
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			domain.OptionalCreateInstanceValues{
				MarketAppId: &marketAppId,
				Reference:   &reference,
				SshKey:      sshKeyValueObject,
			},
		)

		got, err := AdaptToLaunchInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, marketAppId, *got.MarketAppId)
		assert.Equal(t, reference, *got.Reference)
		assert.Equal(t, defaultSshKey, *got.SshKey)
	})
}

func TestAdaptToUpdateInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {

		_, err := AdaptToUpdateInstanceOpts(
			domain.Instance{Type: value_object.InstanceType{Type: "tralala"}},
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
		instanceType := value_object.InstanceType{
			Type: string(publicCloud.TYPENAME_C3_LARGE),
		}
		reference := "reference"
		contractType := enum.ContractTypeMonthly
		contractTerm := enum.ContractTermThree
		billingFrequency := enum.ContractBillingFrequencySix
		rootDiskSize, _ := value_object.NewRootDiskSize(23)

		instance := domain.NewUpdateInstance(
			value_object.NewGeneratedUuid(),
			domain.OptionalUpdateInstanceValues{
				Type:             &instanceType,
				Reference:        &reference,
				ContractType:     &contractType,
				Term:             &contractTerm,
				BillingFrequency: &billingFrequency,
				RootDiskSize:     rootDiskSize,
			})

		got, err := AdaptToUpdateInstanceOpts(instance)

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
	instanceType := value_object.NewUnvalidatedInstanceType(
		string(publicCloud.TYPENAME_C3_LARGE),
	)
	loadBalancerType, _ := value_object.NewInstanceType(
		"loadBalancerType",
		[]string{"loadBalancerType"},
	)

	cpu := domain.NewCpu(1, "cpuUnit")
	memory := domain.NewMemory(2, "memoryUnit")
	publicNetworkSpeed := domain.NewNetworkSpeed(
		3,
		"publicNetworkSpeedUnit",
	)
	privateNetworkSpeed := domain.NewNetworkSpeed(
		4,
		"privateNetworkSpeedUnit",
	)

	resources := domain.NewResources(
		cpu,
		memory,
		publicNetworkSpeed,
		privateNetworkSpeed,
	)

	image := domain.NewImage(
		enum.Ubuntu200464Bit,
		"name",
		"version",
		"family",
		"flavour",
		[]string{"one"},
		[]string{"storageType"},
	)

	rootDiskSize, _ := value_object.NewRootDiskSize(55)

	reverseLookup := "reverseLookup"
	ip := domain.NewIp(
		"1.2.3.4",
		"prefix-length",
		46,
		true,
		false,
		"tralala",
		domain.OptionalIpValues{
			Ddos:          &domain.Ddos{ProtectionType: "protection-type"},
			ReverseLookup: &reverseLookup,
		},
	)

	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	renewalsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2022-12-14 17:09:47",
	)
	createdAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2021-12-14 17:09:47",
	)
	contract, _ := domain.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		renewalsAt,
		createdAt,
		enum.ContractStateActive,
		&endsAt,
	)

	reference := "reference"
	marketAppId := "marketAppId"
	sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)
	startedAt := time.Now()

	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	stickySession := domain.NewStickySession(true, 5)

	host := "host"
	healthCheck := domain.NewHealthCheck(
		enum.MethodGet,
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	loadBalancerConfiguration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		5,
		6,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &stickySession,
			HealthCheck:   &healthCheck,
		},
	)

	loadBalancer := domain.NewLoadBalancer(
		value_object.NewGeneratedUuid(),
		*loadBalancerType,
		resources,
		"region",
		enum.StateCreating,
		*contract,
		domain.Ips{ip},
		domain.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &privateNetwork,
			Configuration:  &loadBalancerConfiguration,
		},
	)

	autoScalingGroupReference, _ := value_object.NewAutoScalingGroupReference(
		"reference",
	)
	autoScalingGroupCreatedAt := time.Now()
	autoScalingGroupUpdatedAt := time.Now()
	autoScalingGroupDesiredAmount := 1
	autoScalingGroupStartsAt := time.Now()
	autoScalingGroupEndsAt := time.Now()
	autoScalingMinimumAmount := 2
	autoScalingMaximumAmount := 3
	autoScalingCpuThreshold := 4
	autoScalingWarmupTime := 5
	autoScalingCooldownTime := 6
	autoScalingGroup := domain.NewAutoScalingGroup(
		value_object.NewGeneratedUuid(),
		"type",
		"state",
		"region",
		*autoScalingGroupReference,
		autoScalingGroupCreatedAt,
		autoScalingGroupUpdatedAt,
		domain.AutoScalingGroupOptions{
			DesiredAmount: &autoScalingGroupDesiredAmount,
			StartsAt:      &autoScalingGroupStartsAt,
			EndsAt:        &autoScalingGroupEndsAt,
			MinimumAmount: &autoScalingMinimumAmount,
			MaximumAmount: &autoScalingMaximumAmount,
			CpuThreshold:  &autoScalingCpuThreshold,
			WarmupTime:    &autoScalingWarmupTime,
			CoolDownTime:  &autoScalingCooldownTime,
			LoadBalancer:  &loadBalancer,
		})

	return domain.NewInstance(
		value_object.NewGeneratedUuid(),
		"region",
		resources,
		image,
		enum.StateCreating,
		"productType",
		false,
		true,
		*rootDiskSize,
		instanceType,
		enum.RootDiskStorageTypeCentral,
		domain.Ips{ip},
		*contract,
		domain.OptionalInstanceValues{
			Reference:        &reference,
			Iso:              &domain.Iso{Id: "isoId"},
			MarketAppId:      &marketAppId,
			SshKey:           sshKeyValueObject,
			StartedAt:        &startedAt,
			PrivateNetwork:   &privateNetwork,
			AutoScalingGroup: &autoScalingGroup,
		},
	)
}
