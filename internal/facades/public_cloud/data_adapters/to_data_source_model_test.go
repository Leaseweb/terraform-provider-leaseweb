package data_adapters

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
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
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	marketAppId := "marketAppId"
	reference := "reference"
	id := "id"
	sshKeyValueObject, _ := value_object.NewSshKey(defaultSshKey)
	autoScalingGroupId := "autoScalingGroupId"
	loadBalancerId := "loadBalancerId"

	instance := generateDomainInstance()
	instance.Id = id
	instance.StartedAt = &startedAt
	instance.MarketAppId = &marketAppId
	instance.Reference = &reference
	instance.SshKey = sshKeyValueObject
	instance.AutoScalingGroup.Id = autoScalingGroupId
	instance.AutoScalingGroup.LoadBalancer.Id = loadBalancerId
	instance.Region = public_cloud.Region{Name: "region"}

	got := adaptInstance(instance)

	assert.Equal(t, id, got.Id.ValueString())
	assert.Equal(t, "region", got.Region.Name.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "productType", got.ProductType.ValueString())
	assert.False(t, got.HasPublicIpv4.ValueBool())
	assert.True(t, got.HasPrivateNetwork.ValueBool())
	assert.Equal(t, "lsw.c3.large", got.Type.Name.ValueString())
	assert.Equal(t, int64(55), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "2019-09-08 00:00:00 +0000 UTC", got.StartedAt.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "UBUNTU_20_04_64BIT", got.Image.Id.ValueString())
	assert.Equal(t, "MONTHLY", got.Contract.Type.ValueString())
	assert.Equal(t, "1.2.3.4", got.Ips[0].Ip.ValueString())
	assert.Equal(t, "cpuUnit", got.Resources.Cpu.Unit.ValueString())
	assert.Equal(
		t,
		autoScalingGroupId,
		got.AutoScalingGroup.Id.ValueString(),
	)
	assert.Equal(t, "isoId", got.Iso.Id.ValueString())
	assert.Equal(t, "id", got.PrivateNetwork.Id.ValueString())
	assert.Equal(
		t,
		loadBalancerId,
		got.AutoScalingGroup.LoadBalancer.Id.ValueString(),
	)
}

func Test_adaptResources(t *testing.T) {
	resources := public_cloud.NewResources(
		public_cloud.Cpu{Unit: "Cpu"},
		public_cloud.Memory{Unit: "Memory"},
		public_cloud.NetworkSpeed{Unit: "publicNetworkSpeed"},
		public_cloud.NetworkSpeed{Unit: "NetworkSpeed"},
	)

	got := adaptResources(resources)

	assert.Equal(t, "Cpu", got.Cpu.Unit.ValueString())
	assert.Equal(t, "Memory", got.Memory.Unit.ValueString())
	assert.Equal(t, `publicNetworkSpeed`, got.PublicNetworkSpeed.Unit.ValueString())
	assert.Equal(t, "NetworkSpeed", got.PrivateNetworkSpeed.Unit.ValueString())
}

func Test_adaptCpu(t *testing.T) {
	cpu := public_cloud.NewCpu(1, "unit")
	got := adaptCpu(cpu)

	assert.Equal(t, int64(1), got.Value.ValueInt64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptNetworkSpeed(t *testing.T) {
	networkSpeed := public_cloud.NewNetworkSpeed(23, "unit")

	got := adaptNetworkSpeed(networkSpeed)

	assert.Equal(t, "unit", got.Unit.ValueString())
	assert.Equal(t, int64(23), got.Value.ValueInt64())
}

func Test_adaptMemory(t *testing.T) {
	memory := public_cloud.NewMemory(1, "unit")

	got := adaptMemory(memory)

	assert.Equal(t, float64(1), got.Value.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptImage(t *testing.T) {
	state := "state"
	stateReason := "stateReason"
	region := public_cloud.Region{Name: "region"}
	createdAt := time.Now()
	updatedAt := time.Now()
	architecture := "architecture"
	version := "version"

	image := public_cloud.NewImage(
		"id",
		"name",
		&version,
		"family",
		"flavour",
		&architecture,
		&state,
		&stateReason,
		&region,
		&createdAt,
		&updatedAt,
		false,
		&public_cloud.StorageSize{Unit: "unit"},
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
	assert.Equal(t, "region", got.Region.Name.ValueString())
	assert.Equal(t, createdAt.String(), got.CreatedAt.ValueString())
	assert.Equal(t, updatedAt.String(), got.UpdatedAt.ValueString())
	assert.False(t, got.Custom.ValueBool())
	assert.Equal(t, "unit", got.StorageSize.Unit.ValueString())
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

	contract, _ := public_cloud.NewContract(
		enum.ContractBillingFrequencySix,
		enum.ContractTermThree,
		enum.ContractTypeMonthly,
		renewalsAt,
		createdAt,
		enum.ContractStateActive,
		&endsAt,
	)

	got := adaptContract(*contract)

	assert.Equal(t, int64(6), got.BillingFrequency.ValueInt64())
	assert.Equal(t, int64(3), got.Term.ValueInt64())
	assert.Equal(t, "MONTHLY", got.Type.ValueString())
	assert.Equal(t, "2023-12-14 17:09:47 +0000 UTC", got.EndsAt.ValueString())
	assert.Equal(t, "2022-12-14 17:09:47 +0000 UTC", got.RenewalsAt.ValueString())
	assert.Equal(t, "2021-12-14 17:09:47 +0000 UTC", got.CreatedAt.ValueString())
	assert.Equal(t, "ACTIVE", got.State.ValueString())
}

func Test_adaptLoadBalancer(t *testing.T) {

	reference := "reference"
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	id := "id"

	entityLoadBalancer := public_cloud.NewLoadBalancer(
		id,
		public_cloud.InstanceType{Name: "type"},
		public_cloud.Resources{Cpu: public_cloud.Cpu{Unit: "Resources"}},
		public_cloud.Region{Name: "region"},
		enum.StateCreating,
		public_cloud.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
		public_cloud.Ips{{Ip: "1.2.3.4"}},
		public_cloud.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &public_cloud.PrivateNetwork{Id: "privateNetworkId"},
			Configuration: &public_cloud.LoadBalancerConfiguration{
				Balance: enum.BalanceSource,
			},
		},
	)

	got := adaptLoadBalancer(entityLoadBalancer)

	assert.Equal(t, id, got.Id.ValueString(), "id is set")
	assert.Equal(t, "type", got.Type.Name.ValueString())
	assert.Equal(t, "Resources", got.Resources.Cpu.Unit.ValueString())
	assert.Equal(t, "region", got.Region.Name.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(6), got.Contract.BillingFrequency.ValueInt64())
	assert.Equal(t, "2019-09-08 00:00:00 +0000 UTC", got.StartedAt.ValueString())
	assert.Equal(t, "1.2.3.4", got.Ips[0].Ip.ValueString())
	assert.Equal(t, "source", got.LoadBalancerConfiguration.Balance.ValueString())
	assert.Equal(t, "privateNetworkId", got.PrivateNetwork.Id.ValueString())
}

func Test_adaptLoadBalancerConfiguration(t *testing.T) {
	configuration := public_cloud.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		1,
		2,
		public_cloud.OptionalLoadBalancerConfigurationOptions{
			StickySession: &public_cloud.StickySession{MaxLifeTime: 32},
			HealthCheck:   &public_cloud.HealthCheck{Method: enum.MethodGet},
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
	healthCheck := public_cloud.NewHealthCheck(
		"method",
		"uri",
		22,
		public_cloud.OptionalHealthCheckValues{Host: &host},
	)

	got := adaptHealthCheck(healthCheck)

	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func Test_adaptStickySession(t *testing.T) {
	stickySession := public_cloud.NewStickySession(false, 1)

	got := adaptStickySession(stickySession)

	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func Test_adaptPrivateNetwork(t *testing.T) {

	privateNetwork := public_cloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)
	got := adaptPrivateNetwork(privateNetwork)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "status", got.Status.ValueString())
	assert.Equal(t, "subnet", got.Subnet.ValueString())
}

func Test_adaptIso(t *testing.T) {
	iso := public_cloud.NewIso("id", "name")
	got := adaptIso(iso)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "name", got.Name.ValueString())
}

func Test_adaptIp(t *testing.T) {

	ip := public_cloud.NewIp(
		"Ip",
		"prefixLength",
		46,
		true,
		false,
		enum.NetworkTypeInternal,
		public_cloud.OptionalIpValues{
			Ddos: &public_cloud.Ddos{ProtectionType: "protection-type"},
		},
	)

	got := adaptIp(ip)

	assert.Equal(t, "Ip", got.Ip.ValueString())
	assert.Equal(t, "prefixLength", got.PrefixLength.ValueString())
	assert.Equal(t, int64(46), got.Version.ValueInt64())
	assert.Equal(t, true, got.NullRouted.ValueBool())
	assert.Equal(t, false, got.MainIp.ValueBool())
	assert.Equal(t, "INTERNAL", got.NetworkType.ValueString())
	assert.Equal(t, "protection-type", got.Ddos.ProtectionType.ValueString())
}

func Test_adaptDdos(t *testing.T) {
	ddos := public_cloud.NewDdos(
		"detectionProfile",
		"protectionType",
	)
	got := adaptDdos(ddos)

	assert.Equal(t, "detectionProfile", got.DetectionProfile.ValueString())
	assert.Equal(t, "protectionType", got.ProtectionType.ValueString())
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

	state := "state"
	stateReason := "stateReason"
	region := public_cloud.Region{Name: "region"}
	createdAt := time.Now()
	updatedAt := time.Now()
	architecture := "architecture"
	version := "version"

	storageSize := public_cloud.NewStorageSize(5, "storageSizeUnit")

	image := public_cloud.NewImage(
		"UBUNTU_20_04_64BIT",
		"name",
		&version,
		"family",
		"flavour",
		&architecture,
		&state,
		&stateReason,
		&region,
		&createdAt,
		&updatedAt,
		false,
		&storageSize,
		[]string{"one"},
		[]string{"storageType"},
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

	stickySession := public_cloud.NewStickySession(true, 5)

	host := "host"
	healthCheck := public_cloud.NewHealthCheck(
		enum.MethodGet,
		"uri",
		22,
		public_cloud.OptionalHealthCheckValues{Host: &host},
	)

	loadBalancerConfiguration := public_cloud.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		5,
		6,
		public_cloud.OptionalLoadBalancerConfigurationOptions{
			StickySession: &stickySession,
			HealthCheck:   &healthCheck,
		},
	)

	loadBalancer := public_cloud.NewLoadBalancer(
		"",
		public_cloud.InstanceType{Name: "type"},
		resources,
		public_cloud.Region{Name: "loadBalancerRegion"},
		enum.StateCreating,
		*contract,
		public_cloud.Ips{ip},
		public_cloud.OptionalLoadBalancerValues{
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
	autoScalingGroup := public_cloud.NewAutoScalingGroup(
		"",
		"type",
		"state",
		public_cloud.Region{Name: "autoScalingGroupRegion"},
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
			LoadBalancer:  &loadBalancer,
		})

	return public_cloud.NewInstance(
		"",
		public_cloud.Region{Name: "region"},
		resources,
		image,
		enum.StateCreating,
		"productType",
		false,
		true,
		*rootDiskSize,
		public_cloud.InstanceType{Name: "lsw.c3.large"},
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

func Test_adaptStorageSize(t *testing.T) {
	storageSize := public_cloud.NewStorageSize(1, "unit")

	got := adaptStorageSize(storageSize)

	assert.Equal(t, float64(1), got.Size.ValueFloat64())
	assert.Equal(t, "unit", got.Unit.ValueString())
}

func Test_adaptPrice(t *testing.T) {
	got := adaptPrice(public_cloud.NewPrice("hourly", "monthly"))

	assert.Equal(t, "hourly", got.HourlyPrice.ValueString())
	assert.Equal(t, "monthly", got.MonthlyPrice.ValueString())
}

func Test_adaptStorage(t *testing.T) {
	got := adaptStorage(
		public_cloud.NewStorage(
			public_cloud.Price{HourlyPrice: "hourlyLocal"},
			public_cloud.Price{HourlyPrice: "hourlyCentral"},
		),
	)

	assert.Equal(t, "hourlyLocal", got.Local.HourlyPrice.ValueString())
	assert.Equal(t, "hourlyCentral", got.Central.HourlyPrice.ValueString())
}

func Test_adaptPrices(t *testing.T) {
	got := adaptPrices(
		public_cloud.NewPrices(
			"currency",
			"currencySymbol",
			public_cloud.Price{HourlyPrice: "computePrice"},
			public_cloud.Storage{Local: public_cloud.Price{HourlyPrice: "storagePrice"}},
		),
	)

	assert.Equal(t, "currency", got.Currency.ValueString())
	assert.Equal(t, "currencySymbol", got.CurrencySymbol.ValueString())
	assert.Equal(t, "computePrice", got.Compute.HourlyPrice.ValueString())
	assert.Equal(t, "storagePrice", got.Storage.Local.HourlyPrice.ValueString())
}

func Test_adaptInstanceType(t *testing.T) {
	t.Run("required values are adapted", func(t *testing.T) {
		got := adaptInstanceType(
			public_cloud.NewInstanceType(
				"name",
				public_cloud.Resources{Cpu: public_cloud.Cpu{Unit: "cpuUnit"}},
				public_cloud.Prices{Currency: "currency"},
				public_cloud.OptionalInstanceTypeValues{},
			),
		)

		assert.Equal(t, "name", got.Name.ValueString())
		assert.Equal(t, "cpuUnit", got.Resources.Cpu.Unit.ValueString())
		assert.Equal(t, "currency", got.Prices.Currency.ValueString())
		assert.Nil(t, got.StorageTypes)
	})

	t.Run("optional values are adapted", func(t *testing.T) {
		got := adaptInstanceType(
			public_cloud.NewInstanceType(
				"",
				public_cloud.Resources{},
				public_cloud.Prices{},
				public_cloud.OptionalInstanceTypeValues{
					StorageTypes: &public_cloud.StorageTypes{enum.StorageTypeLocal},
				},
			),
		)

		assert.Equal(t, []string{"LOCAL"}, got.StorageTypes)
	})
}

func Test_adaptRegion(t *testing.T) {
	sdkRegion := public_cloud.NewRegion("name", "location")

	want := model.Region{
		Location: basetypes.NewStringValue("location"),
		Name:     basetypes.NewStringValue("name"),
	}
	got := adaptRegion(sdkRegion)

	assert.Equal(t, want, *got)
}
