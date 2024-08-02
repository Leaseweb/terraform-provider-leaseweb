package to_data_source_model

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestAdaptInstances(t *testing.T) {
	id := value_object.NewGeneratedUuid()
	instances := domain.Instances{{Id: id}}

	got := AdaptInstances(instances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, id.String(), got.Instances[0].Id.ValueString())
}

func Test_adaptInstance(t *testing.T) {
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	marketAppId := "marketAppId"
	reference := "reference"
	id := value_object.NewGeneratedUuid()
	sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)
	autoScalingGroupId := value_object.NewGeneratedUuid()
	loadBalancerId := value_object.NewGeneratedUuid()

	instance := generateDomainInstance()
	instance.Id = id
	instance.StartedAt = &startedAt
	instance.MarketAppId = &marketAppId
	instance.Reference = &reference
	instance.SshKey = sshKeyValueObject
	instance.AutoScalingGroup.Id = autoScalingGroupId
	instance.AutoScalingGroup.LoadBalancer.Id = loadBalancerId

	got := adaptInstance(instance)

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
		"lsw.c3.large",
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
		"Image should be set",
	)
	assert.Equal(
		t,
		"MONTHLY",
		got.Contract.Type.ValueString(),
		"Contract should be set",
	)
	assert.Equal(
		t,
		"1.2.3.4",
		got.Ips[0].Ip.ValueString(),
		"Ip should be set",
	)
	assert.Equal(
		t,
		"cpuUnit",
		got.Resources.Cpu.Unit.ValueString(),
		"PrivateNetwork should be set",
	)
	assert.Equal(
		t,
		autoScalingGroupId.String(),
		got.AutoScalingGroup.Id.ValueString(),
		"AutoScalingGroup should be set",
	)
	assert.Equal(
		t,
		"isoId",
		got.Iso.Id.ValueString(),
		"Iso should be set",
	)
	assert.Equal(
		t,
		"id",
		got.PrivateNetwork.Id.ValueString(),
		"PrivateNetwork should be set",
	)
	assert.Equal(
		t,
		loadBalancerId.String(),
		got.AutoScalingGroup.LoadBalancer.Id.ValueString(),
		"loadBalancer should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Volume.Unit.ValueString(),
		"volume should be set",
	)
}

func Test_adaptResources(t *testing.T) {
	resources := domain.NewResources(
		domain.Cpu{Unit: "Cpu"},
		domain.Memory{Unit: "Memory"},
		domain.NetworkSpeed{Unit: "publicNetworkSpeed"},
		domain.NetworkSpeed{Unit: "NetworkSpeed"},
	)

	got := adaptResources(resources)

	assert.Equal(
		t,
		"Cpu",
		got.Cpu.Unit.ValueString(),
		"Cpu should be set",
	)
	assert.Equal(
		t,
		"Memory",
		got.Memory.Unit.ValueString(),
		"Memory should be set",
	)
	assert.Equal(
		t,
		"publicNetworkSpeed",
		got.PublicNetworkSpeed.Unit.ValueString(),
		"publicNetworkSpeed should be set",
	)
	assert.Equal(
		t,
		"NetworkSpeed",
		got.PrivateNetworkSpeed.Unit.ValueString(),
		"NetworkSpeed should be set",
	)
}

func Test_adaptCpu(t *testing.T) {

	cpu := domain.NewCpu(1, "unit")
	got := adaptCpu(cpu)

	assert.Equal(
		t,
		int64(1),
		got.Value.ValueInt64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)

}

func Test_adaptNetworkSpeed(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(23, "unit")

	got := adaptNetworkSpeed(networkSpeed)

	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
	assert.Equal(
		t,
		int64(23),
		got.Value.ValueInt64(),
		"value should be set",
	)
}

func Test_adaptMemory(t *testing.T) {
	memory := domain.NewMemory(1, "unit")

	got := adaptMemory(memory)

	assert.Equal(
		t,
		float64(1),
		got.Value.ValueFloat64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
}

func Test_adaptImage(t *testing.T) {
	state := "state"
	stateReason := "stateReason"
	region := "region"
	createdAt := time.Now()
	updatedAt := time.Now()
	custom := false

	image := domain.NewImage(
		"id",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		&state,
		&stateReason,
		&region,
		&createdAt,
		&updatedAt,
		&custom,
		[]string{"one"},
		[]string{"storageType"},
	)

	got := adaptImage(image)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "name", got.Name.ValueString())
	assert.Equal(t, "version", got.Version.ValueString())
	assert.Equal(t, "family", got.Family.ValueString())
	assert.Equal(t, "flavour", got.Flavour.ValueString())
	assert.Equal(t, "architecture", got.Architecture.ValueString())
	assert.Equal(t, "state", got.State.ValueString())
	assert.Equal(t, "stateReason", got.StateReason.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, createdAt.String(), got.CreatedAt.ValueString())
	assert.Equal(t, updatedAt.String(), got.UpdatedAt.ValueString())
	assert.False(t, got.Custom.ValueBool())
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("one")},
		got.MarketApps,
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("storageType")},
		got.StorageTypes,
	)
}

func Test_adaptContract(t *testing.T) {

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

	got := adaptContract(*contract)

	assert.Equal(
		t,
		int64(6),
		got.BillingFrequency.ValueInt64(),
		"billingFrequency should be set",
	)
	assert.Equal(
		t,
		int64(3),
		got.Term.ValueInt64(),
		"term should be set",
	)
	assert.Equal(
		t,
		"MONTHLY",
		got.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		"2023-12-14 17:09:47 +0000 UTC",
		got.EndsAt.ValueString(),
		"endsAt should be set",
	)
	assert.Equal(
		t,
		"2022-12-14 17:09:47 +0000 UTC",
		got.RenewalsAt.ValueString(),
		"renewalsAt should be set",
	)
	assert.Equal(
		t,
		"2021-12-14 17:09:47 +0000 UTC",
		got.CreatedAt.ValueString(),
		"createdAt should be set",
	)
	assert.Equal(
		t,
		"ACTIVE",
		got.State.ValueString(),
		"state should be set",
	)
}

func Test_adaptLoadBalancer(t *testing.T) {

	reference := "reference"
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	id := value_object.NewGeneratedUuid()

	entityLoadBalancer := domain.NewLoadBalancer(
		id,
		value_object.NewUnvalidatedInstanceType("type"),
		domain.Resources{Cpu: domain.Cpu{Unit: "Resources"}},
		"region",
		enum.StateCreating,
		domain.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
		domain.Ips{{Ip: "1.2.3.4"}},
		domain.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &domain.PrivateNetwork{Id: "privateNetworkId"},
			Configuration: &domain.LoadBalancerConfiguration{
				Balance: enum.BalanceSource,
			},
		},
	)

	got := adaptLoadBalancer(entityLoadBalancer)

	assert.Equal(t, id.String(), got.Id.ValueString(), "id is set")
	assert.Equal(
		t,
		"type",
		got.Type.ValueString(),
		"type is set",
	)
	assert.Equal(
		t,
		"Resources",
		got.Resources.Cpu.Unit.ValueString(),
		"Resources is set",
	)
	assert.Equal(
		t,
		"region",
		got.Region.ValueString(),
		"region is set",
	)
	assert.Equal(
		t,
		"reference",
		got.Reference.ValueString(),
		"reference is set",
	)
	assert.Equal(
		t,
		"CREATING",
		got.State.ValueString(),
		"state is set",
	)
	assert.Equal(
		t,
		int64(6),
		got.Contract.BillingFrequency.ValueInt64(),
		"Contract is set",
	)
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.StartedAt.ValueString(),
		"startedAt is set",
	)
	assert.Equal(
		t,
		"1.2.3.4",
		got.Ips[0].Ip.ValueString(),
		"ips is set",
	)
	assert.Equal(
		t,
		"source",
		got.LoadBalancerConfiguration.Balance.ValueString(),
		"configuration is set",
	)
	assert.Equal(
		t,
		"privateNetworkId",
		got.PrivateNetwork.Id.ValueString(),
		"PrivateNetwork is set",
	)
}

func Test_adaptLoadBalancerConfiguration(t *testing.T) {
	configuration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		1,
		2,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &domain.StickySession{MaxLifeTime: 32},
			HealthCheck:   &domain.HealthCheck{Method: enum.MethodGet},
		},
	)

	got := adaptLoadBalancerConfiguration(configuration)

	assert.Equal(t, int64(32), got.StickySession.MaxLifeTime.ValueInt64())
	assert.Equal(t, "source", got.Balance.ValueString())
	assert.Equal(t, "GET", got.HealthCheck.Method.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(1), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(2), got.TargetPort.ValueInt64())
}

func Test_adaptHealthCheck(t *testing.T) {
	host := "host"
	healthCheck := domain.NewHealthCheck(
		"method",
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got := adaptHealthCheck(healthCheck)

	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func Test_adaptStickySession(t *testing.T) {
	stickySession := domain.NewStickySession(false, 1)

	got := adaptStickySession(stickySession)

	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func Test_adaptPrivateNetwork(t *testing.T) {

	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)
	got := adaptPrivateNetwork(privateNetwork)

	assert.Equal(
		t, "id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"status",
		got.Status.ValueString(),
		"status should be set",
	)
	assert.Equal(
		t,
		"subnet",
		got.Subnet.ValueString(),
		"subnet should be set",
	)
}

func Test_adaptIso(t *testing.T) {
	iso := domain.NewIso("id", "name")
	got := adaptIso(iso)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "name", got.Name.ValueString())
}

func Test_adaptIp(t *testing.T) {

	ip := domain.NewIp(
		"Ip",
		"prefixLength",
		46,
		true,
		false,
		enum.NetworkTypeInternal,
		domain.OptionalIpValues{
			Ddos: &domain.Ddos{ProtectionType: "protection-type"},
		},
	)

	got := adaptIp(ip)

	assert.Equal(
		t,
		"Ip",
		got.Ip.ValueString(),
		"Ip should be set",
	)
	assert.Equal(
		t,
		"prefixLength",
		got.PrefixLength.ValueString(),
		"prefix-length should be set",
	)
	assert.Equal(
		t,
		int64(46),
		got.Version.ValueInt64(),
		"version should be set",
	)
	assert.Equal(
		t,
		true,
		got.NullRouted.ValueBool(),
		"nullRouted should be set",
	)
	assert.Equal(
		t,
		false,
		got.MainIp.ValueBool(),
		"mainIp should be set",
	)
	assert.Equal(
		t,
		"INTERNAL",
		got.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"protection-type",
		got.Ddos.ProtectionType.ValueString(),
		"Ddos should be set",
	)
}

func Test_adaptDdos(t *testing.T) {
	ddos := domain.NewDdos(
		"detectionProfile",
		"protectionType",
	)
	got := adaptDdos(ddos)

	assert.Equal(
		t,
		"detectionProfile",
		got.DetectionProfile.ValueString(),
		"detectionProfile should be set",
	)
	assert.Equal(
		t,
		"protectionType",
		got.ProtectionType.ValueString(),
		"protectionType should be set",
	)

}

func generateDomainInstance() domain.Instance {
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

	state := "state"
	stateReason := "stateReason"
	region := "region"
	createdAt := time.Now()
	updatedAt := time.Now()
	custom := false

	image := domain.NewImage(
		"UBUNTU_20_04_64BIT",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		&state,
		&stateReason,
		&region,
		&createdAt,
		&updatedAt,
		&custom,
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
	contractCreatedAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2021-12-14 17:09:47",
	)
	contract, _ := domain.NewContract(
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
		value_object.NewUnvalidatedInstanceType("type"),
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

	volume := domain.NewVolume(1, "unit")

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
		value_object.NewUnvalidatedInstanceType(
			string(publicCloud.TYPENAME_C3_LARGE),
		),
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
			Volume:           &volume,
		},
	)
}

func Test_adaptVolume(t *testing.T) {
	volume := domain.NewVolume(1, "unit")

	got := adaptVolume(volume)

	assert.Equal(t, float64(1), got.Size.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}
