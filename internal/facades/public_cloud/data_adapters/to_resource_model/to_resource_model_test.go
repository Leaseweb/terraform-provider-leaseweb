package to_resource_model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/stretchr/testify/assert"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_adaptImage(t *testing.T) {
	image := public_cloud.NewImage(
		"UBUNTU_20_04_64BIT",
		"name",
		"family",
		"flavour",
		false,
	)

	got, err := adaptImage(context.TODO(), image)

	assert.NoError(t, err)

	assert.Equal(t, "UBUNTU_20_04_64BIT", got.Id.ValueString())
	assert.Equal(t, "name", got.Name.ValueString())
	assert.Equal(t, "family", got.Family.ValueString())
	assert.Equal(t, "flavour", got.Flavour.ValueString())
	assert.False(t, got.Custom.ValueBool())
}

func Test_AdaptContract(t *testing.T) {
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

	contract, _ := public_cloud.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		renewalsAt,
		createdAt,
		enum.ContractStateActive,
		&endsAt,
	)
	got, err := adaptContract(context.TODO(), *contract)

	assert.NoError(t, err)

	assert.Equal(t, int64(6), got.BillingFrequency.ValueInt64())
	assert.Equal(t, int64(3), got.Term.ValueInt64())
	assert.Equal(t, "MONTHLY", got.Type.ValueString())
	assert.Equal(t, "2023-12-14 17:09:47 +0000 UTC", got.EndsAt.ValueString())
	assert.Equal(t, "2022-12-14 17:09:47 +0000 UTC", got.RenewalsAt.ValueString())
	assert.Equal(t, "2021-12-14 17:09:47 +0000 UTC", got.CreatedAt.ValueString())
	assert.Equal(t, "ACTIVE", got.State.ValueString())
}

func Test_adaptPrivateNetwork(t *testing.T) {
	privateNetwork := public_cloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	got, err := adaptPrivateNetwork(context.TODO(), privateNetwork)

	assert.NoError(t, err)

	assert.Equal(t, "id", got.PrivateNetworkId.ValueString())
	assert.Equal(t, "status", got.Status.ValueString())
	assert.Equal(t, "subnet", got.Subnet.ValueString())
}

func Test_adaptCpu(t *testing.T) {
	entityCpu := public_cloud.NewCpu(1, "unit")
	got, err := adaptCpu(context.TODO(), entityCpu)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Value.ValueInt64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptMemory(t *testing.T) {
	memory := public_cloud.NewMemory(1, "unit")

	got, err := adaptMemory(context.TODO(), memory)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), got.Value.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptNetworkSpeed(t *testing.T) {
	networkSpeed := public_cloud.NewNetworkSpeed(1, "unit")

	got, err := adaptNetworkSpeed(context.TODO(), networkSpeed)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Value.ValueInt64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptResources(t *testing.T) {
	resources := public_cloud.NewResources(
		public_cloud.Cpu{Unit: "cpu"},
		public_cloud.Memory{Unit: "memory"},
		public_cloud.NetworkSpeed{Unit: "publicNetworkSpeed"},
		public_cloud.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	got, err := adaptResources(context.TODO(), resources)

	assert.NoError(t, err)

	cpu := model.Cpu{}
	got.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(t, "cpu", cpu.Unit.ValueString())

	memory := model.Memory{}
	got.Memory.As(context.TODO(), &memory, basetypes.ObjectAsOptions{})
	assert.Equal(t, "memory", memory.Unit.ValueString())

	publicNetworkSpeed := model.NetworkSpeed{}
	got.PublicNetworkSpeed.As(
		context.TODO(),
		&publicNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, "publicNetworkSpeed", publicNetworkSpeed.Unit.ValueString())

	privateNetworkSpeed := model.NetworkSpeed{}
	got.PrivateNetworkSpeed.As(
		context.TODO(),
		&privateNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"privateNetworkSpeed",
		privateNetworkSpeed.Unit.ValueString(),
	)
}

func Test_adaptDdos(t *testing.T) {
	ddos := public_cloud.NewDdos(
		"detectionProfile",
		"protectionType",
	)

	got, err := adaptDdos(context.TODO(), ddos)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"detectionProfile",
		got.DetectionProfile.ValueString(),
	)
	assert.Equal(t, "protectionType", got.ProtectionType.ValueString())
}

func Test_adaptIp(t *testing.T) {
	reverseLookup := "reverse-lookup"

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
	got, err := adaptIp(context.TODO(), ip)

	assert.NoError(t, err)

	assert.Equal(t, "1.2.3.4", got.Ip.ValueString())
	assert.Equal(t, "prefix-length", got.PrefixLength.ValueString())
	assert.Equal(t, int64(46), got.Version.ValueInt64())
	assert.Equal(t, true, got.NullRouted.ValueBool())
	assert.Equal(t, false, got.MainIp.ValueBool())
	assert.Equal(t, "tralala", got.NetworkType.ValueString())
	assert.Equal(t, "reverse-lookup", got.ReverseLookup.ValueString())

	ddos := model.Ddos{}
	got.Ddos.As(context.TODO(), &ddos, basetypes.ObjectAsOptions{})
	assert.Equal(t, "protection-type", ddos.ProtectionType.ValueString())
}

func TestAdaptInstance(t *testing.T) {
	var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	marketAppId := "marketAppId"
	reference := "reference"
	id := "id"
	rootDiskSize, _ := value_object.NewRootDiskSize(32)
	autoScalingGroupId := "autoScalingGroupId"
	sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

	instance := generateDomainInstance()
	instance.Id = id
	instance.RootDiskSize = *rootDiskSize
	instance.StartedAt = &startedAt
	instance.MarketAppId = &marketAppId
	instance.Reference = &reference
	instance.SshKey = sshKeyValueObject
	instance.PrivateNetwork.Id = "privateNetworkId"
	instance.AutoScalingGroup.Id = autoScalingGroupId
	instance.Resources.Cpu.Unit = "cpu"
	instance.Type = "instanceType"

	got, err := AdaptInstance(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, id, got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "productType", got.ProductType.ValueString())
	assert.False(t, got.HasPublicIpv4.ValueBool())
	assert.True(t, got.HasPrivateNetwork.ValueBool())
	assert.False(t, got.HasUserData.ValueBool())
	assert.Equal(t, int64(32), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.StartedAt.ValueString(),
	)
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "instanceType", got.Type.ValueString())

	image := model.Image{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.Id.ValueString())

	contract := model.Contract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	iso := model.Iso{}
	got.Iso.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
	assert.Equal(t, "isoId", iso.Id.ValueString())

	privateNetwork := model.PrivateNetwork{}
	got.PrivateNetwork.As(
		context.TODO(),
		&privateNetwork,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"privateNetworkId",
		privateNetwork.PrivateNetworkId.ValueString(),
	)

	autoScalingGroup := model.AutoScalingGroup{}
	got.AutoScalingGroup.As(
		context.TODO(),
		&autoScalingGroup,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		autoScalingGroupId,
		autoScalingGroup.Id.ValueString(),
	)

	var ips []model.Ip
	got.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "1.2.3.4", ips[0].Ip.ValueString())

	resources := model.Resources{}
	cpu := model.Cpu{}
	got.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
	resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(t, "cpu", cpu.Unit.ValueString())

	// TODO Enable SSH key support
	//assert.Equal(t, sshKey, got.SshKey.ValueString())
}

func Test_adaptAutoScalingGroup(t *testing.T) {
	desiredAmount := 1
	createdAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-09-08T00:00:00Z")
	startsAt, _ := time.Parse(time.RFC3339, "2010-09-08T00:00:00Z")
	endsAt, _ := time.Parse(time.RFC3339, "2011-09-08T00:00:00Z")
	minimumAmount := 2
	maximumAmount := 3
	cpuThreshold := 4
	warmupTime := 5
	cooldownTime := 6
	id := "id"
	reference, _ := value_object.NewAutoScalingGroupReference("reference")

	autoScalingGroup := public_cloud.NewAutoScalingGroup(
		id,
		"type",
		"state",
		"region",
		*reference,
		createdAt,
		updatedAt,
		public_cloud.AutoScalingGroupOptions{
			DesiredAmount: &desiredAmount,
			StartsAt:      &startsAt,
			EndsAt:        &endsAt,
			MinimumAmount: &minimumAmount,
			MaximumAmount: &maximumAmount,
			CpuThreshold:  &cpuThreshold,
			WarmupTime:    &warmupTime,
			CoolDownTime:  &cooldownTime,
		},
	)

	got, err := adaptAutoScalingGroup(
		context.TODO(),
		autoScalingGroup,
	)

	assert.NoError(t, err)

	assert.Equal(t, id, got.Id.ValueString())
	assert.Equal(t, "type", got.Type.ValueString())
	assert.Equal(t, "state", got.State.ValueString())
	assert.Equal(t, int64(1), got.DesiredAmount.ValueInt64())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.CreatedAt.ValueString(),
	)
	assert.Equal(
		t,
		"2020-09-08 00:00:00 +0000 UTC",
		got.UpdatedAt.ValueString(),
	)
	assert.Equal(
		t,
		"2010-09-08 00:00:00 +0000 UTC",
		got.StartsAt.ValueString(),
	)
	assert.Equal(
		t,
		"2011-09-08 00:00:00 +0000 UTC",
		got.EndsAt.ValueString(),
	)
	assert.Equal(t, int64(2), got.MinimumAmount.ValueInt64())
	assert.Equal(t, int64(3), got.MaximumAmount.ValueInt64())
	assert.Equal(t, int64(4), got.CpuThreshold.ValueInt64())
	assert.Equal(t, int64(5), got.WarmupTime.ValueInt64())
	assert.Equal(t, int64(6), got.CooldownTime.ValueInt64())
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
		"instanceType",
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
