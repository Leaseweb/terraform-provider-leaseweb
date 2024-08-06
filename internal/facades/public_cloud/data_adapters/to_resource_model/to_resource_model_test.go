package to_resource_model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/stretchr/testify/assert"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_adaptImage(t *testing.T) {
	var marketApps []string
	var storageTypes []string

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
		&domain.StorageSize{Unit: "unit"},
		[]string{"one"},
		[]string{"storageType"},
	)

	got, err := adaptImage(context.TODO(), image)

	assert.NoError(t, err)

	got.MarketApps.ElementsAs(context.TODO(), &marketApps, false)
	got.StorageTypes.ElementsAs(context.TODO(), &storageTypes, false)

	assert.Equal(t, "UBUNTU_20_04_64BIT", got.Id.ValueString())
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

	assert.Len(t, marketApps, 1)
	assert.Equal(t, "one", marketApps[0])

	assert.Len(t, storageTypes, 1)
	assert.Equal(t, "storageType", storageTypes[0])

	storageSize := model.StorageSize{}
	got.StorageSize.As(context.TODO(), &storageSize, basetypes.ObjectAsOptions{})
	assert.Equal(t, "unit", storageSize.Unit.ValueString())
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

	contract, _ := domain.NewContract(
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
	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	got, err := adaptPrivateNetwork(context.TODO(), privateNetwork)

	assert.NoError(t, err)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "status", got.Status.ValueString())
	assert.Equal(t, "subnet", got.Subnet.ValueString())
}

func Test_adaptCpu(t *testing.T) {
	entityCpu := domain.NewCpu(1, "unit")
	got, err := adaptCpu(context.TODO(), entityCpu)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Value.ValueInt64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptMemory(t *testing.T) {
	memory := domain.NewMemory(1, "unit")

	got, err := adaptMemory(context.TODO(), memory)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), got.Value.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptNetworkSpeed(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(1, "unit")

	got, err := adaptNetworkSpeed(context.TODO(), networkSpeed)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Value.ValueInt64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptResources(t *testing.T) {
	resources := domain.NewResources(
		domain.Cpu{Unit: "cpu"},
		domain.Memory{Unit: "memory"},
		domain.NetworkSpeed{Unit: "publicNetworkSpeed"},
		domain.NetworkSpeed{Unit: "privateNetworkSpeed"},
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

func Test_adaptHealthCheck(t *testing.T) {
	host := "host"
	healthCheck := domain.NewHealthCheck(
		enum.MethodGet,
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got, err := adaptHealthCheck(context.TODO(), healthCheck)

	assert.NoError(t, err)
	assert.Equal(t, "GET", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func Test_adaptStickySession(t *testing.T) {
	stickySession := domain.NewStickySession(false, 1)

	got, err := adaptStickySession(context.TODO(), stickySession)

	assert.Nil(t, err)
	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func Test_adaptLoadBalancerConfiguration(t *testing.T) {

	loadBalancerConfiguration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		5,
		6,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &domain.StickySession{MaxLifeTime: 5},
			HealthCheck:   &domain.HealthCheck{Method: enum.MethodHead},
		},
	)

	got, err := adaptLoadBalancerConfiguration(
		context.TODO(),
		loadBalancerConfiguration,
	)

	assert.NoError(t, err)
	assert.Equal(t, "source", got.Balance.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(5), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(6), got.TargetPort.ValueInt64())

	stickySession := model.StickySession{}
	got.StickySession.As(
		context.TODO(),
		&stickySession,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, int64(5), stickySession.MaxLifeTime.ValueInt64())

	healthCheck := model.HealthCheck{}
	got.HealthCheck.As(
		context.TODO(),
		&healthCheck,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, "HEAD", healthCheck.Method.ValueString())
}

func Test_adaptDdos(t *testing.T) {
	ddos := domain.NewDdos(
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

func Test_adaptLoadBalancer(t *testing.T) {
	t.Run("loadBalancer Conversion works", func(t *testing.T) {
		reference := "reference"
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		id := value_object.NewGeneratedUuid()

		loadBalancer := domain.NewLoadBalancer(
			id,
			domain.InstanceType{Name: "instanceType"},
			domain.Resources{Cpu: domain.Cpu{Unit: "cpu"}},
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

		got, err := adaptLoadBalancer(
			context.TODO(),
			loadBalancer,
		)

		assert.NoError(t, err)

		assert.Equal(t, id.String(), got.Id.ValueString())
		assert.Equal(
			t,
			"{\"unit\":\"cpu\",\"value\":0}",
			got.Resources.Attributes()["cpu"].String(),
		)
		assert.Equal(t, "region", got.Region.ValueString())
		assert.Equal(t, "reference", got.Reference.ValueString())
		assert.Equal(t, "CREATING", got.State.ValueString())

		assert.Equal(
			t,
			"6",
			got.Contract.Attributes()["billing_frequency"].String(),
		)

		assert.Equal(
			t,
			"2019-09-08 00:00:00 +0000 UTC",
			got.StartedAt.ValueString(),
		)

		var ips []model.Ip
		got.Ips.ElementsAs(context.TODO(), &ips, false)
		assert.Equal(t, "1.2.3.4", ips[0].Ip.ValueString())

		loadBalancerConfiguration := model.LoadBalancerConfiguration{}
		got.LoadBalancerConfiguration.As(
			context.TODO(),
			&loadBalancerConfiguration,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(t, "source", loadBalancerConfiguration.Balance.ValueString())

		privateNetwork := model.PrivateNetwork{}
		got.PrivateNetwork.As(
			context.TODO(),
			&privateNetwork,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(t, "privateNetworkId", privateNetwork.Id.ValueString())

		instanceType := model.InstanceType{}
		got.Type.As(
			context.TODO(),
			&instanceType,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(t, "instanceType", instanceType.Name.ValueString())
	})
}

func TestAdaptInstance(t *testing.T) {
	var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	marketAppId := "marketAppId"
	reference := "reference"
	id := value_object.NewGeneratedUuid()
	rootDiskSize, _ := value_object.NewRootDiskSize(32)
	autoScalingGroupId := value_object.NewGeneratedUuid()
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
	instance.Type.Name = "instanceType"

	got, err := AdaptInstance(instance, context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, id.String(), got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "productType", got.ProductType.ValueString())
	assert.False(t, got.HasPublicIpv4.ValueBool())
	assert.True(t, got.HasPrivateNetwork.ValueBool())
	assert.Equal(t, int64(32), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.StartedAt.ValueString(),
	)
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())

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
		privateNetwork.Id.ValueString(),
	)

	autoScalingGroup := model.AutoScalingGroup{}
	got.AutoScalingGroup.As(
		context.TODO(),
		&autoScalingGroup,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		autoScalingGroupId.String(),
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

	volume := model.Volume{}
	got.Volume.As(context.TODO(), &volume, basetypes.ObjectAsOptions{})
	assert.Equal(t, "unit", volume.Unit.ValueString())

	instanceType := model.InstanceType{}
	got.Type.As(context.TODO(), &instanceType, basetypes.ObjectAsOptions{})
	assert.Equal(t, "instanceType", instanceType.Name.ValueString())

	assert.Equal(t, sshKey, got.SshKey.ValueString())
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
	id := value_object.NewGeneratedUuid()
	reference, _ := value_object.NewAutoScalingGroupReference("reference")
	loadBalancerId := value_object.NewGeneratedUuid()

	autoScalingGroup := domain.NewAutoScalingGroup(
		id,
		"type",
		"state",
		"region",
		*reference,
		createdAt,
		updatedAt,
		domain.AutoScalingGroupOptions{
			DesiredAmount: &desiredAmount,
			StartsAt:      &startsAt,
			EndsAt:        &endsAt,
			MinimumAmount: &minimumAmount,
			MaximumAmount: &maximumAmount,
			CpuThreshold:  &cpuThreshold,
			WarmupTime:    &warmupTime,
			CoolDownTime:  &cooldownTime,
			LoadBalancer: &domain.LoadBalancer{
				Id:        loadBalancerId,
				StartedAt: &time.Time{},
			},
		},
	)

	got, err := adaptAutoScalingGroup(
		context.TODO(),
		autoScalingGroup,
	)

	assert.NoError(t, err)

	assert.Equal(t, id.String(), got.Id.ValueString())
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

	loadBalancer := model.LoadBalancer{}
	got.LoadBalancer.As(
		context.TODO(),
		&loadBalancer,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, loadBalancerId.String(), loadBalancer.Id.ValueString())
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

	storageSize := domain.NewStorageSize(1, "unit")

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
		&storageSize,
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

	loadBalancerCpu := domain.NewCpu(45, "loadBalancerCpuUnit")
	loadBalancerMemory := domain.NewMemory(2, "loadBalancerMemoryUnit")
	loadBalancerPrivateNetworkSpeed := domain.NewNetworkSpeed(
		55,
		"loadBalancerPrivateNetworkSpeedUnit",
	)
	loadBalancerPublicNetworkSpeed := domain.NewNetworkSpeed(
		56,
		"loadBalancerPublicNetworkSpeedUnit",
	)
	instanceTypeResources := domain.NewResources(
		loadBalancerCpu,
		loadBalancerMemory,
		loadBalancerPrivateNetworkSpeed,
		loadBalancerPublicNetworkSpeed,
	)
	instanceTypePricesCompute := domain.NewPrice("5", "6")
	instanceTypePricesStorageLocal := domain.NewPrice(
		"7",
		"8",
	)
	instanceTypePricesStorageCenral := domain.NewPrice(
		"23",
		"4",
	)
	instanceTypePricesStorage := domain.NewStorage(
		instanceTypePricesStorageLocal,
		instanceTypePricesStorageCenral,
	)
	instanceTypePrices := domain.NewPrices(
		"currency",
		"currencySymbol",
		instanceTypePricesCompute,
		instanceTypePricesStorage,
	)
	instanceTypeStorageTypes := domain.StorageTypes{"storageType"}
	instanceType := domain.NewInstanceType(
		"instanceType",
		instanceTypeResources,
		instanceTypePrices,
		domain.OptionalInstanceTypeValues{StorageTypes: &instanceTypeStorageTypes},
	)

	loadBalancer := domain.NewLoadBalancer(
		value_object.NewGeneratedUuid(),
		instanceType,
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
			Volume:           &volume,
		},
	)
}

func Test_adaptVolume(t *testing.T) {
	got, err := adaptVolume(
		context.TODO(),
		domain.Volume{
			Size: 2,
			Unit: "unit",
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, float64(2), got.Size.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptStorageSize(t *testing.T) {
	got, err := adaptStorageSize(
		context.TODO(),
		domain.StorageSize{
			Size: 2,
			Unit: "unit",
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, float64(2), got.Size.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptPrice(t *testing.T) {
	price := domain.NewPrice("1", "2")
	got, err := adaptPrice(context.TODO(), price)

	assert.NoError(t, err)
	assert.Equal(t, "1", got.HourlyPrice.ValueString())
	assert.Equal(t, "2", got.MonthlyPrice.ValueString())
}

func Test_adaptStorage(t *testing.T) {
	storage := domain.NewStorage(
		domain.Price{HourlyPrice: "1"},
		domain.Price{HourlyPrice: "3"},
	)
	got, err := adaptStorage(context.TODO(), storage)

	assert.NoError(t, err)

	local := model.Price{}
	got.Local.As(context.TODO(), &local, basetypes.ObjectAsOptions{})
	assert.Equal(t, "1", local.HourlyPrice.ValueString())

	central := model.Price{}
	got.Central.As(context.TODO(), &central, basetypes.ObjectAsOptions{})
	assert.Equal(t, "3", central.HourlyPrice.ValueString())
}

func Test_adaptPrices(t *testing.T) {
	prices := domain.NewPrices(
		"currency",
		"currencySymbol",
		domain.Price{HourlyPrice: "1"},
		domain.Storage{Central: domain.Price{HourlyPrice: "3"}},
	)
	got, err := adaptPrices(context.TODO(), prices)

	assert.NoError(t, err)

	assert.Equal(t, "currency", got.Currency.ValueString())
	assert.Equal(t, "currencySymbol", got.CurrencySymbol.ValueString())

	compute := model.Price{}
	got.Compute.As(context.TODO(), &compute, basetypes.ObjectAsOptions{})
	assert.Equal(t, "1", compute.HourlyPrice.ValueString())

	storage := model.Storage{}
	got.Storage.As(context.TODO(), &storage, basetypes.ObjectAsOptions{})
	storageCentral := model.Price{}
	storage.Central.As(
		context.TODO(),
		&storageCentral,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, "3", storageCentral.HourlyPrice.ValueString())
}

func Test_adaptInstanceType(t *testing.T) {
	var storageTypes []string

	instanceType := domain.NewInstanceType(
		"name",
		domain.Resources{Cpu: domain.Cpu{Unit: "unit"}},
		domain.Prices{Currency: "currency"},
		domain.OptionalInstanceTypeValues{
			StorageTypes: &domain.StorageTypes{"storageType"},
		},
	)

	got, err := adaptInstanceType(context.TODO(), instanceType)

	assert.NoError(t, err)

	assert.Equal(t, "name", got.Name.ValueString())

	resources := model.Resources{}
	got.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
	cpu := model.Cpu{}
	resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(t, "unit", cpu.Unit.ValueString())

	prices := model.Prices{}
	got.Prices.As(context.TODO(), &prices, basetypes.ObjectAsOptions{})
	assert.Equal(t, "currency", prices.Currency.ValueString())

	got.StorageTypes.ElementsAs(context.TODO(), &storageTypes, false)
	assert.Equal(t, []string{"storageType"}, storageTypes)
}
