package to_data_source_model

import (
	"testing"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/stretchr/testify/assert"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestAdaptInstances(t *testing.T) {
	id := "id"
	instances := public_cloud.Instances{{Id: id}}

	got := AdaptInstances(instances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, id, got.Instances[0].Id.ValueString())
}

func Test_adaptInstance(t *testing.T) {
	marketAppId := "marketAppId"
	reference := "reference"
	id := "id"
	sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)

	instance := generateDomainInstance()
	instance.Id = id
	instance.MarketAppId = &marketAppId
	instance.Reference = &reference
	instance.SshKey = sshKeyValueObject
	instance.Region = "region"

	got := adaptInstance(instance)

	assert.Equal(t, id, got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "lsw.c3.large", got.Type.ValueString())
	assert.Equal(t, int64(55), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "UBUNTU_20_04_64BIT", got.Image.Id.ValueString())
	assert.Equal(t, "MONTHLY", got.Contract.Type.ValueString())
	assert.Equal(t, "1.2.3.4", got.Ips[0].Ip.ValueString())
}

func Test_adaptImage(t *testing.T) {
	image := public_cloud.NewImage(
		"id",
		"",
		"",
		"",
		false,
	)

	got := adaptImage(image)

	assert.Equal(t, "id", got.Id.ValueString())
}

func Test_adaptContract(t *testing.T) {
	contract, _ := public_cloud.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		time.Now(),
		time.Now(),
		enum.ContractStateActive,
		nil,
	)

	got := adaptContract(*contract)

	assert.Equal(t, int64(6), got.BillingFrequency.ValueInt64())
	assert.Equal(t, int64(3), got.Term.ValueInt64())
	assert.Equal(t, "MONTHLY", got.Type.ValueString())
	assert.Equal(t, "ACTIVE", got.State.ValueString())
}

func Test_adaptIp(t *testing.T) {

	ip := public_cloud.NewIp(
		"Ip",
		"",
		46,
		true,
		false,
		enum.NetworkTypeInternal,
		public_cloud.OptionalIpValues{},
	)

	got := adaptIp(ip)

	assert.Equal(t, "Ip", got.Ip.ValueString())
}

func generateDomainInstance() public_cloud.Instance {
	cpu := public_cloud.NewCpu(1, "cpuUnit")
	memory := public_cloud.NewMemory(2, "memoryUnit")
	publicNetworkSpeed := public_cloud.NewNetworkSpeed(
		3,
		"publicNetworkSpeedUnit",
	)
	privateNetworkSpeed := public_cloud.NewNetworkSpeed(
		4,
		"privateNetworkSpeedUnit",
	)

	resources := public_cloud.NewResources(
		cpu,
		memory,
		publicNetworkSpeed,
		privateNetworkSpeed,
	)

	image := public_cloud.NewImage(
		"UBUNTU_20_04_64BIT",
		"name",
		"family",
		"flavour",
		false,
	)

	rootDiskSize, _ := value_object.NewRootDiskSize(55)

	reverseLookup := "reverseLookup"
	ip := public_cloud.NewIp(
		"1.2.3.4",
		"prefix-length",
		46,
		true,
		false,
		"tralala",
		public_cloud.OptionalIpValues{
			Ddos:          &public_cloud.Ddos{ProtectionType: "protection-type"},
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
	contractCreatedAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2021-12-14 17:09:47",
	)
	contract, _ := public_cloud.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		renewalsAt,
		contractCreatedAt,
		enum.ContractStateActive,
		&endsAt,
	)

	reference := "reference"
	marketAppId := "marketAppId"
	sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)
	startedAt := time.Now()

	privateNetwork := public_cloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
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
	autoScalingGroup := public_cloud.NewAutoScalingGroup(
		"",
		"type",
		"state",
		"autoScalingGroupRegion",
		*autoScalingGroupReference,
		autoScalingGroupCreatedAt,
		autoScalingGroupUpdatedAt,
		public_cloud.AutoScalingGroupOptions{
			DesiredAmount: &autoScalingGroupDesiredAmount,
			StartsAt:      &autoScalingGroupStartsAt,
			EndsAt:        &autoScalingGroupEndsAt,
			MinimumAmount: &autoScalingMinimumAmount,
			MaximumAmount: &autoScalingMaximumAmount,
			CpuThreshold:  &autoScalingCpuThreshold,
			WarmupTime:    &autoScalingWarmupTime,
			CoolDownTime:  &autoScalingCooldownTime,
		})

	return public_cloud.NewInstance(
		"",
		"region",
		resources,
		image,
		enum.StateCreating,
		"productType",
		false,
		true,
		false,
		*rootDiskSize,
		"lsw.c3.large",
		enum.StorageTypeCentral,
		public_cloud.Ips{ip},
		*contract,
		public_cloud.OptionalInstanceValues{
			Reference:        &reference,
			Iso:              &public_cloud.Iso{Id: "isoId"},
			MarketAppId:      &marketAppId,
			SshKey:           sshKeyValueObject,
			StartedAt:        &startedAt,
			PrivateNetwork:   &privateNetwork,
			AutoScalingGroup: &autoScalingGroup,
		},
	)
}
