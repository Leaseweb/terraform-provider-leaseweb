package public_cloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_convertImageToResourceModel(t *testing.T) {
	image := domain.NewImage(
		enum.Ubuntu200464Bit,
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	got, err := convertImageToResourceModel(context.TODO(), image)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"UBUNTU_20_04_64BIT",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		got.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		got.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		got.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		got.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		got.Architecture.ValueString(),
		"architecture should be set",
	)

	var marketApps []string
	got.MarketApps.ElementsAs(context.TODO(), &marketApps, false)
	assert.Len(t, marketApps, 1)
	assert.Equal(
		t,
		"one",
		marketApps[0],
		"marketApps should be set",
	)

	var storageTypes []string
	got.StorageTypes.ElementsAs(context.TODO(), &storageTypes, false)
	assert.Len(t, storageTypes, 1)
	assert.Equal(
		t,
		"storageType",
		storageTypes[0],
		"storageTypes should be set",
	)
}

func Test_convertContractToResourceModel(t *testing.T) {
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
	got, err := convertContractToResourceModel(context.TODO(), *contract)

	assert.NoError(t, err)

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

func Test_convertPrivateNetworkToResourceModel(t *testing.T) {
	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	got, err := convertPrivateNetworkToResourceModel(
		context.TODO(),
		privateNetwork,
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"id",
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

func Test_convertCpuToResourceModel(t *testing.T) {
	entityCpu := domain.NewCpu(1, "unit")
	got, err := convertCpuToResourceModel(context.TODO(), entityCpu)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}

func Test_convertMemoryToResourceModel(t *testing.T) {
	memory := domain.NewMemory(1, "unit")

	got, err := convertMemoryToResourceModel(context.TODO(), memory)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), got.Value.ValueFloat64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}

func Test_convertNetworkSpeedToResourceModel(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(1, "unit")

	got, err := convertNetworkSpeedToResourceModel(context.TODO(), networkSpeed)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}

func Test_convertResourcesToResourceModel(t *testing.T) {
	resources := domain.NewResources(
		domain.Cpu{Unit: "cpu"},
		domain.Memory{Unit: "memory"},
		domain.NetworkSpeed{Unit: "publicNetworkSpeed"},
		domain.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	got, err := convertResourcesToResourceModel(context.TODO(), resources)

	assert.NoError(t, err)

	cpu := model.Cpu{}
	got.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"cpu",
		cpu.Unit.ValueString(),
		"cpu should be set",
	)

	memory := model.Memory{}
	got.Memory.As(context.TODO(), &memory, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"memory",
		memory.Unit.ValueString(),
		"memory should be set",
	)

	publicNetworkSpeed := model.NetworkSpeed{}
	got.PublicNetworkSpeed.As(
		context.TODO(),
		&publicNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"publicNetworkSpeed",
		publicNetworkSpeed.Unit.ValueString(),
		"publicNetworkSpeed should be set",
	)

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
		"privateNetworkSpeed should be set",
	)
}

func Test_convertHealthCheckToResourceModel(t *testing.T) {
	host := "host"
	healthCheck := domain.NewHealthCheck(
		enum.MethodGet,
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got, err := convertHealthCheckToResourceModel(context.TODO(), healthCheck)

	assert.NoError(t, err)
	assert.Equal(t, "GET", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func Test_convertStickySessionToResourceModel(t *testing.T) {
	stickySession := domain.NewStickySession(false, 1)

	got, err := convertStickySessionToResourceModel(context.TODO(), stickySession)

	assert.Nil(t, err)
	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func Test_convertLoadBalancerConfigurationToResourceModel(t *testing.T) {

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

	got, err := convertLoadBalancerConfigurationToResourceModel(
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

func Test_convertDdosToResourceModel(t *testing.T) {
	ddos := domain.NewDdos("detectionProfile", "protectionType")

	got, err := convertDdosToResourceModel(context.TODO(), ddos)

	assert.NoError(t, err)

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

func Test_convertIpToResourceModel(t *testing.T) {
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
	got, err := convertIpToResourceModel(context.TODO(), ip)

	assert.NoError(t, err)

	assert.Equal(
		t,
		"1.2.3.4",
		got.Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"prefix-length",
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
		"tralala",
		got.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"reverse-lookup",
		got.ReverseLookup.ValueString(),
		"reverseLookup should be set",
	)

	ddos := model.Ddos{}
	got.Ddos.As(context.TODO(), &ddos, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"protection-type",
		ddos.ProtectionType.ValueString(),
		"ddos should be set",
	)
}

func Test_convertLoadBalancerToResourceModel(t *testing.T) {
	t.Run("loadBalancer Conversion works", func(t *testing.T) {
		reference := "reference"
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		id := value_object.NewGeneratedUuid()

		loadBalancer := domain.NewLoadBalancer(
			id,
			"type",
			domain.Resources{Cpu: domain.Cpu{Unit: "cpu"}},
			"region",
			enum.StateCreating,
			domain.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
			domain.Ips{{Ip: "1.2.3.4"}},
			domain.OptionalLoadBalancerValues{
				Reference:      &reference,
				StartedAt:      &startedAt,
				PrivateNetwork: &domain.PrivateNetwork{Id: "privateNetworkId"},
				Configuration:  &domain.LoadBalancerConfiguration{Balance: enum.BalanceSource},
			},
		)

		got, err := convertLoadBalancerToResourceModel(
			context.TODO(),
			loadBalancer,
		)

		assert.NoError(t, err)

		assert.Equal(t, id.String(), got.Id.ValueString())
		assert.Equal(t, "type", got.Type.ValueString())
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
		got.Ips.ElementsAs(
			context.TODO(),
			&ips,
			false,
		)
		assert.Equal(t, "1.2.3.4", ips[0].Ip.ValueString())

		loadBalancerConfiguration := model.LoadBalancerConfiguration{}
		got.LoadBalancerConfiguration.As(
			context.TODO(),
			&loadBalancerConfiguration,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"source",
			loadBalancerConfiguration.Balance.ValueString(),
		)

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
	})
}

func Test_convertInstanceToResourceModel(t *testing.T) {
	var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

	t.Run("instance is converted correctly", func(t *testing.T) {
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		marketAppId := "marketAppId"
		reference := "reference"
		id := value_object.NewGeneratedUuid()
		rootDiskSize, _ := value_object.NewRootDiskSize(32)
		autoScalingGroupId := value_object.NewGeneratedUuid()
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		instance := generateDomainInstance()
		instance.Id = id
		instance.Type = enum.InstanceTypeM5A4Xlarge
		instance.RootDiskSize = *rootDiskSize
		instance.StartedAt = &startedAt
		instance.MarketAppId = &marketAppId
		instance.Reference = &reference
		instance.SshKey = sshKeyValueObject
		instance.PrivateNetwork.Id = "privateNetworkId"
		instance.AutoScalingGroup.Id = autoScalingGroupId
		instance.Resources.Cpu.Unit = "cpu"

		got, err := convertInstanceToResourceModel(instance, context.TODO())

		assert.NoError(t, err)
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
			"lsw.m5a.4xlarge",
			got.Type.ValueString(),
			"type should be set",
		)
		assert.Equal(
			t,
			int64(32),
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
			"reference",
			got.Reference.ValueString(),
			"reference should be set",
		)

		image := model.Image{}
		got.Image.As(
			context.TODO(),
			&image,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"UBUNTU_20_04_64BIT",
			image.Id.ValueString(),
			"image should be set",
		)

		contract := model.Contract{}
		got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"MONTHLY",
			contract.Type.ValueString(),
			"contract should be set",
		)

		iso := model.Iso{}
		got.Iso.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"isoId",
			iso.Id.ValueString(),
			"iso should be set",
		)

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
			"privateNetwork should be set",
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
			"autoScalingGroup should be set",
		)

		var ips []model.Ip
		got.Ips.ElementsAs(context.TODO(), &ips, false)
		assert.Len(t, ips, 1)
		assert.Equal(
			t,
			"1.2.3.4",
			ips[0].Ip.ValueString(),
			"ip should be set",
		)

		resources := model.Resources{}
		cpu := model.Cpu{}
		got.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
		resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
		assert.Equal(
			t,
			"cpu",
			cpu.Unit.ValueString(),
			"privateNetwork should be set",
		)

		assert.Equal(t, sshKey, got.SshKey.ValueString())
	})
}

func Test_convertAutoScalingGroupToResourceModel(t *testing.T) {
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
			LoadBalancer:  &domain.LoadBalancer{Id: loadBalancerId, StartedAt: &time.Time{}},
		},
	)

	got, err := convertAutoScalingGroupToResourceModel(context.TODO(), autoScalingGroup)

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

func Test_convertInstanceResourceModelToCreateInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		got, err := convertInstanceResourceModelToCreateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, enum.InstanceTypeM5A4Xlarge, got.Type)
		assert.Equal(t, enum.RootDiskStorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, enum.Ubuntu200464Bit, got.Image.Id)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermThree, got.Contract.Term)
		assert.Equal(t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
	})

	t.Run("optional values are passed", func(t *testing.T) {
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		got, err := convertInstanceResourceModelToCreateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, 55, got.RootDiskSize.Value)
		assert.Equal(t, defaultSshKey, got.SshKey.String())
	})

	t.Run(
		"returns error if invalid rootDiskStorageType is passed",
		func(t *testing.T) {
			rootDiskStorageType := "tralala"
			instance := generateInstanceModel(
				&rootDiskStorageType,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			)

			_, err := convertInstanceResourceModelToCreateInstanceOpts(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instanceType := "tralala"
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				&instanceType,
			)

			_, err := convertInstanceResourceModelToCreateInstanceOpts(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run("returns error if invalid imageId is passed", func(t *testing.T) {
		imageId := "tralala"
		instance := generateInstanceModel(
			nil,
			&imageId,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		_, err := convertInstanceResourceModelToCreateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("returns error if invalid contractType is passed", func(t *testing.T) {
		contractType := "tralala"
		instance := generateInstanceModel(
			nil,
			nil,
			&contractType,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		_, err := convertInstanceResourceModelToCreateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				&contractTerm,
				nil,
				nil,
				nil,
				nil,
			)

			_, err := convertInstanceResourceModelToCreateInstanceOpts(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				&billingFrequency,
				nil,
				nil,
				nil,
			)

			_, err := convertInstanceResourceModelToCreateInstanceOpts(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run("returns error if invalid sshKey is passed", func(t *testing.T) {
		sshKey := "tralala"
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			&sshKey,
			nil,
			nil,
		)

		_, err := convertInstanceResourceModelToCreateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "ssh key is invalid")
	})

	t.Run(
		"returns error if invalid rootDiskSize is passed",
		func(t *testing.T) {
			rootDiskSize := 1
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				&rootDiskSize,
				nil,
			)

			_, err := convertInstanceResourceModelToCreateInstanceOpts(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "1")
		},
	)
}

func generateInstanceModel(
	rootDiskStorageType *string,
	imageId *string,
	contractType *string,
	contractTerm *int,
	billingFrequency *int,
	sshKey *string,
	rootDiskSize *int,
	instanceType *string,
) model.Instance {
	defaultRootDiskStorageType := "CENTRAL"
	defaultImageId := "UBUNTU_20_04_64BIT"
	defaultContractType := "MONTHLY"
	defaultContractTerm := 3
	defaultBillingFrequency := 1
	defaultRootDiskSize := 55
	defaultInstanceType := "lsw.m5a.4xlarge"

	if rootDiskStorageType == nil {
		rootDiskStorageType = &defaultRootDiskStorageType
	}
	if imageId == nil {
		imageId = &defaultImageId
	}
	if contractType == nil {
		contractType = &defaultContractType
	}
	if contractTerm == nil {
		contractTerm = &defaultContractTerm
	}
	if billingFrequency == nil {
		billingFrequency = &defaultBillingFrequency
	}
	if rootDiskSize == nil {
		rootDiskSize = &defaultRootDiskSize
	}
	if sshKey == nil {
		sshKey = &defaultSshKey
	}
	if instanceType == nil {
		instanceType = &defaultInstanceType
	}

	image, _ := types.ObjectValueFrom(
		context.TODO(),
		model.Image{}.AttributeTypes(),
		model.Image{
			Id:           basetypes.NewStringValue(*imageId),
			Name:         basetypes.NewStringUnknown(),
			Version:      basetypes.NewStringUnknown(),
			Family:       basetypes.NewStringUnknown(),
			Flavour:      basetypes.NewStringUnknown(),
			Architecture: basetypes.NewStringUnknown(),
			MarketApps:   basetypes.NewListUnknown(types.StringType),
			StorageTypes: basetypes.NewListUnknown(types.StringType),
		},
	)

	contract, _ := types.ObjectValueFrom(
		context.TODO(),
		model.Contract{}.AttributeTypes(),
		model.Contract{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			EndsAt:           basetypes.NewStringUnknown(),
			RenewalsAt:       basetypes.NewStringUnknown(),
			CreatedAt:        basetypes.NewStringUnknown(),
			State:            basetypes.NewStringUnknown(),
		},
	)

	instance := model.Instance{
		Id:                  basetypes.NewStringValue(value_object.NewGeneratedUuid().String()),
		Region:              basetypes.NewStringValue("region"),
		Type:                basetypes.NewStringValue(*instanceType),
		RootDiskStorageType: basetypes.NewStringValue(*rootDiskStorageType),
		RootDiskSize:        basetypes.NewInt64Value(int64(*rootDiskSize)),
		Image:               image,
		Contract:            contract,
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
		SshKey:              basetypes.NewStringValue(*sshKey),
	}

	return instance
}

func Test_convertInstancesToDataSourceModel(t *testing.T) {
	id := value_object.NewGeneratedUuid()
	instances := domain.Instances{{Id: id}}

	got := convertInstancesToDataSourceModel(instances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, id.String(), got.Instances[0].Id.ValueString())
}

func Test_convertInstanceToDataSourceModel(t *testing.T) {
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

	got := convertInstanceToDataSourceModel(instance)

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
}

func Test_convertResourcesToDataSourceModel(t *testing.T) {
	resources := domain.NewResources(
		domain.Cpu{Unit: "Cpu"},
		domain.Memory{Unit: "Memory"},
		domain.NetworkSpeed{Unit: "publicNetworkSpeed"},
		domain.NetworkSpeed{Unit: "NetworkSpeed"},
	)

	got := convertResourcesToDataSourceModel(resources)

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

func Test_convertCpuToDataSourceModel(t *testing.T) {

	cpu := domain.NewCpu(1, "unit")
	got := convertCpuToDataSourceModel(cpu)

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

func Test_convertNetworkSpeedToDataSourceModel(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(23, "unit")

	got := convertNetworkSpeedToDataSourceModel(networkSpeed)

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

func Test_convertMemoryToDataSourceModel(t *testing.T) {
	memory := domain.NewMemory(1, "unit")

	got := convertMemoryToDataSourceModel(memory)

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

func Test_convertImageToDataSourceModel(t *testing.T) {

	image := domain.NewImage(
		"id",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	got := convertImageToDataSourceModel(image)

	assert.Equal(
		t,
		"id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		got.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		got.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		got.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		got.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		got.Architecture.ValueString(),
		"architecture should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("one")},
		got.MarketApps,
		"marketApps should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("storageType")},
		got.StorageTypes,
		"storageTypes should be set",
	)
}

func Test_convertContractToDataSourceModel(t *testing.T) {

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

	got := convertContractToDataSourceModel(*contract)

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

func Test_convertLoadBalancerToDataSourceModel(t *testing.T) {

	reference := "reference"
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	id := value_object.NewGeneratedUuid()

	entityLoadBalancer := domain.NewLoadBalancer(
		id,
		"type",
		domain.Resources{Cpu: domain.Cpu{Unit: "Resources"}},
		"region",
		enum.StateCreating,
		domain.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
		domain.Ips{{Ip: "1.2.3.4"}},
		domain.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &domain.PrivateNetwork{Id: "privateNetworkId"},
			Configuration:  &domain.LoadBalancerConfiguration{Balance: enum.BalanceSource},
		},
	)

	got := convertLoadBalancerToDataSourceModel(entityLoadBalancer)

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

func Test_convertLoadBalancerConfigurationToDataSourceModel(t *testing.T) {
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

	got := convertLoadBalancerConfigurationToDataSourceModel(configuration)

	assert.Equal(t, int64(32), got.StickySession.MaxLifeTime.ValueInt64())
	assert.Equal(t, "source", got.Balance.ValueString())
	assert.Equal(t, "GET", got.HealthCheck.Method.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(1), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(2), got.TargetPort.ValueInt64())
}

func Test_convertHealthCheckToDataSourceModel(t *testing.T) {
	host := "host"
	healthCheck := domain.NewHealthCheck(
		"method",
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got := convertHealthCheckToDataSourceModel(healthCheck)

	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func Test_convertStickySessionToDataSourceModel(t *testing.T) {
	stickySession := domain.NewStickySession(false, 1)

	got := convertStickySessionToDataSourceModel(stickySession)

	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func Test_convertPrivateNetworkToDataSourceModel(t *testing.T) {

	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)
	got := convertPrivateNetworkToDataSourceModel(privateNetwork)

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

func Test_convertIsoToDataSourceModel(t *testing.T) {
	iso := domain.NewIso("id", "name")
	got := convertIsoToDataSourceModel(iso)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "name", got.Name.ValueString())
}

func Test_convertIpToDataSourceModel(t *testing.T) {

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

	got := convertIpToDataSourceModel(ip)

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

func Test_convertDdosToDataSourceModel(t *testing.T) {
	ddos := domain.NewDdos("detectionProfile", "protectionType")
	got := convertDdosToDataSourceModel(ddos)

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

func Test_convertInstanceResourceModelToUpdateInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		contract, _ := types.ObjectValueFrom(
			context.TODO(),
			model.Contract{}.AttributeTypes(),
			model.Contract{
				Type:             basetypes.NewStringValue("MONTHLY"),
				Term:             basetypes.NewInt64Value(3),
				BillingFrequency: basetypes.NewInt64Value(3),
			},
		)

		instance := model.Instance{
			Id:           basetypes.NewStringValue(id.String()),
			Contract:     contract,
			RootDiskSize: basetypes.NewInt64Value(65),
		}

		got, diags := convertInstanceResourceModelToUpdateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.Nil(t, diags)
		assert.Equal(t, id, got.Id)
	})

	t.Run("optional values are set", func(t *testing.T) {
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		got, diags := convertInstanceResourceModelToUpdateInstanceOpts(
			instance,
			context.TODO(),
		)

		assert.Nil(t, diags)
		assert.Equal(t, enum.InstanceTypeM5A4Xlarge, got.Type)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermThree, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, 55, got.RootDiskSize.Value)
	})
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

	image := domain.NewImage(
		enum.Ubuntu200464Bit,
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
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
		"type",
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

	autoScalingGroupReference, _ := value_object.NewAutoScalingGroupReference("reference")
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
		enum.InstanceTypeC3Large,
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

func Test_returnError(t *testing.T) {
	t.Run("diagnostics contain errors", func(t *testing.T) {
		diags := diag.Diagnostics{}
		diags.AddError("summary", "detail")

		got := returnError("functionName", diags)
		want := `functionName: "summary" "detail"`

		assert.Error(t, got)
		assert.Equal(t, want, got.Error())
	})

	t.Run("diagnostics do not contain errors", func(t *testing.T) {
		diags := diag.Diagnostics{}

		got := returnError("functionName", diags)

		assert.NoError(t, got)
	})
}

func Test_convertIntArrayToInt64(t *testing.T) {
	want := []int64{5}
	got := convertIntArrayToInt64([]int{5})

	assert.Equal(t, want, got)
}
