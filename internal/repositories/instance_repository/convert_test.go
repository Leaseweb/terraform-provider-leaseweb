package instance_repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var instanceId = "5d7f8262-d77f-4476-8da8-6a84f8f2ae8d"

var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_convertImage(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkImage := publicCloud.NewImageDetails(
			publicCloud.IMAGEID_UBUNTU_24_04_64_BIT,
			"name",
			"version",
			"family",
			"flavour",
			"architecture",
			[]string{"marketApp"},
			[]string{"storageType"},
		)

		got, err := convertImage(*sdkImage)

		assert.Nil(t, err)
		assert.Equal(t, enum.Ubuntu240464Bit, got.Id)
		assert.Equal(t, "name", got.Name)
		assert.Equal(t, "version", got.Version)
		assert.Equal(t, "family", got.Family)
		assert.Equal(t, "flavour", got.Flavour)
		assert.Equal(t, "architecture", got.Architecture)
		assert.Equal(t, []string{"marketApp"}, got.MarketApps)
		assert.Equal(t, []string{"storageType"}, got.StorageTypes)
	})

	t.Run("invalid imageId returns error", func(t *testing.T) {
		sdkImage := publicCloud.ImageDetails{Id: "tralala"}

		_, err := convertImage(sdkImage)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertNetworkSpeed(t *testing.T) {
	sdkNetworkSpeed := publicCloud.NewNetworkSpeed(1, "unit")
	got := convertNetworkSpeed(*sdkNetworkSpeed)

	assert.Equal(t, 1, got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_convertMemory(t *testing.T) {
	sdkMemory := publicCloud.NewMemory(1, "unit")
	got := convertMemory(*sdkMemory)

	assert.Equal(t, float64(1), got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_convertCpu(t *testing.T) {
	sdkCpu := publicCloud.NewCpu(1, "unit")
	got := convertCpu(*sdkCpu)

	assert.Equal(t, 1, got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_convertResources(t *testing.T) {
	sdkResources := publicCloud.NewResources(
		publicCloud.Cpu{Unit: "cpu"},
		publicCloud.Memory{Unit: "memory"},
		publicCloud.NetworkSpeed{Unit: "publicNetworkSpeed"},
		publicCloud.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	got := convertResources(*sdkResources)

	assert.Equal(t, "cpu", got.Cpu.Unit)
	assert.Equal(t, "memory", got.Memory.Unit)
	assert.Equal(t, "publicNetworkSpeed", got.PublicNetworkSpeed.Unit)
	assert.Equal(t, "privateNetworkSpeed", got.PrivateNetworkSpeed.Unit)
}

func Test_convertInstance(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		startedAt := time.Now()
		autoScalingGroupId := value_object.NewGeneratedUuid()

		sdkInstance := generateInstanceDetails(t, &startedAt, nil)
		autoScalingGroup := domain.AutoScalingGroup{Id: autoScalingGroupId}

		got, err := convertInstance(sdkInstance, &autoScalingGroup)

		assert.NoError(t, err)
		assert.Equal(
			t,
			"5d7f8262-d77f-4476-8da8-6a84f8f2ae8d",
			got.Id.String(),
		)
		assert.Equal(t, "lsw.m3.large", got.Type)
		assert.Equal(t, "cpu", got.Resources.Cpu.Unit)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, startedAt, *got.StartedAt)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, enum.StateRunning, got.State)
		assert.Equal(t, "productType", got.ProductType)
		assert.True(t, got.HasPublicIpv4)
		assert.False(t, got.HasPrivateNetwork)
		assert.Equal(t, 6, got.RootDiskSize.Value)
		assert.Equal(t, enum.RootDiskStorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "isoId", got.Iso.Id)
		assert.Equal(t, "privateNetworkId", got.PrivateNetwork.Id)
		assert.Equal(t, enum.Centos764Bit, got.Image.Id)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(t, autoScalingGroupId, got.AutoScalingGroup.Id)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{Id: "tralala"},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid Image returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{
				Id:    uuid.New().String(),
				Image: publicCloud.ImageDetails{Id: "tralala"},
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{
				Id:    uuid.New().String(),
				Image: publicCloud.ImageDetails{Id: publicCloud.IMAGEID_ALMALINUX_8_64_BIT},
				State: "tralala",
			},
			nil,
		)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{
				Id:           uuid.New().String(),
				Image:        publicCloud.ImageDetails{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
				State:        publicCloud.STATE_RUNNING,
				RootDiskSize: 5000,
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{
				Id:                  uuid.New().String(),
				Image:               publicCloud.ImageDetails{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
				State:               publicCloud.STATE_RUNNING,
				RootDiskSize:        50,
				RootDiskStorageType: "tralala",
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid ip returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{
				Id:                  uuid.New().String(),
				Image:               publicCloud.ImageDetails{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
				State:               publicCloud.STATE_RUNNING,
				RootDiskSize:        50,
				RootDiskStorageType: publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
				Ips:                 []publicCloud.IpDetails{{NetworkType: "tralala"}},
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		_, err := convertInstance(
			publicCloud.InstanceDetails{
				Id:                  "5d7f8262-d77f-4476-8da8-6a84f8f2ae8d",
				Image:               publicCloud.ImageDetails{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
				State:               publicCloud.STATE_RUNNING,
				Type:                publicCloud.INSTANCETYPENAME_M3_LARGE,
				RootDiskSize:        50,
				RootDiskStorageType: publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
				Contract:            publicCloud.Contract{BillingFrequency: 55},
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})
}

func generateInstanceDetails(
	t *testing.T,
	startedAt *time.Time,
	id *string,
) publicCloud.InstanceDetails {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	if id == nil {
		id = &instanceId
	}

	return *publicCloud.NewInstanceDetails(
		*id,
		publicCloud.INSTANCETYPENAME_M3_LARGE,
		publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "cpu"}},
		"region",
		*publicCloud.NewNullableString(&reference),
		*publicCloud.NewNullableTime(startedAt),
		*publicCloud.NewNullableString(&marketAppId),
		publicCloud.STATE_RUNNING,
		"productType",
		true,
		false,
		6,
		"CENTRAL",
		publicCloud.Contract{
			BillingFrequency: 1,
			Type:             publicCloud.CONTRACTTYPE_HOURLY,
			State:            publicCloud.CONTRACTSTATE_ACTIVE,
		},
		*publicCloud.NewNullableIso(&publicCloud.Iso{Id: "isoId"}),
		*publicCloud.NewNullablePrivateNetwork(&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"}),
		publicCloud.ImageDetails{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
		[]publicCloud.IpDetails{{Ip: "1.2.3.4", NetworkType: publicCloud.NETWORKTYPE_PUBLIC}},
		*publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
			Id: "autoscalingGroupId",
		}),
	)
}

func Test_convertDdos(t *testing.T) {
	got := convertDdos(publicCloud.Ddos{
		DetectionProfile: "detectionProfile",
		ProtectionType:   "protectionType",
	})

	assert.Equal(t, "detectionProfile", got.DetectionProfile)
	assert.Equal(t, "protectionType", got.ProtectionType)
}

func Test_convertIp(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		reverseLookup := "reverseLookup"

		sdkIp := publicCloud.NewIpDetails(
			"1.2.3.4",
			"prefixLength",
			5,
			true,
			false,
			publicCloud.NETWORKTYPE_INTERNAL,
			*publicCloud.NewNullableString(&reverseLookup),
			*publicCloud.NewNullableDdos(&publicCloud.Ddos{DetectionProfile: "detectionProfile"}),
		)

		got, err := convertIp(*sdkIp)

		assert.NoError(t, err)
		assert.Equal(t, "1.2.3.4", got.Ip)
		assert.Equal(t, "prefixLength", got.PrefixLength)
		assert.Equal(t, 5, got.Version)
		assert.True(t, got.NullRouted)
		assert.False(t, got.MainIp)
		assert.Equal(t, enum.NetworkTypeInternal, got.NetworkType)
		assert.Equal(t, "reverseLookup", *got.ReverseLookup)
		assert.Equal(t, "detectionProfile", got.Ddos.DetectionProfile)
	})

	t.Run("error returned for invalid networkType", func(t *testing.T) {
		_, err := convertIp(publicCloud.IpDetails{NetworkType: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertIps(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := convertIps([]publicCloud.IpDetails{{
			Ip:          "1.2.3.4",
			NetworkType: publicCloud.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := convertIps([]publicCloud.IpDetails{{NetworkType: "tralala"}})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertContract(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		endsAt := time.Now()
		renewalsAt := time.Now()
		createdAt := time.Now()

		sdkContract := publicCloud.NewContract(
			0,
			1,
			publicCloud.CONTRACTTYPE_MONTHLY,
			*publicCloud.NewNullableTime(&endsAt),
			renewalsAt,
			createdAt,
			publicCloud.CONTRACTSTATE_ACTIVE,
		)

		got, err := convertContract(*sdkContract)

		assert.NoError(t, err)
		assert.Equal(t, enum.ContractBillingFrequencyZero, got.BillingFrequency)
		assert.Equal(t, enum.ContractTermOne, got.Term)
		assert.Equal(t, enum.ContractTypeMonthly, got.Type)
		assert.Equal(t, endsAt, *got.EndsAt)
		assert.Equal(t, renewalsAt, got.RenewalsAt)
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, enum.ContractStateActive, got.State)
	})

	t.Run("error returned for invalid billingFrequency", func(t *testing.T) {
		sdkContract := publicCloud.Contract{BillingFrequency: 45}

		_, err := convertContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "45")
	})

	t.Run("error returned for invalid term", func(t *testing.T) {
		sdkContract := publicCloud.Contract{BillingFrequency: 0, Term: 55}

		_, err := convertContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("error returned for invalid type", func(t *testing.T) {
		sdkContract := publicCloud.Contract{
			BillingFrequency: 0,
			Term:             0,
			Type:             "tralala",
		}

		_, err := convertContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("error returned for invalid state", func(t *testing.T) {
		sdkContract := publicCloud.Contract{
			BillingFrequency: 0,
			Term:             0,
			Type:             publicCloud.CONTRACTTYPE_HOURLY,
			State:            "tralala",
		}

		_, err := convertContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run(
		"error returned when contract cannot be created",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				BillingFrequency: 0,
				Term:             0,
				Type:             publicCloud.CONTRACTTYPE_MONTHLY,
				State:            publicCloud.CONTRACTSTATE_ACTIVE,
			}

			_, err := convertContract(sdkContract)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "contract.term cannot be 0")
		},
	)

}

func Test_convertIso(t *testing.T) {
	got := convertIso(*publicCloud.NewIso("id", "name"))

	assert.Equal(t, "id", got.Id)
	assert.Equal(t, "name", got.Name)
}

func Test_convertPrivateNetwork(t *testing.T) {
	got := convertPrivateNetwork(*publicCloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	))

	assert.Equal(t, "id", got.Id)
	assert.Equal(t, "status", got.Status)
	assert.Equal(t, "subnet", got.Subnet)
}

func Test_convertAutoScalingGroup(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now()
		startsAt := time.Now()
		endsAt := time.Now()
		minimumAmount := int32(1)
		maximumAmount := int32(2)
		cpuThreshold := int32(3)
		warmupTime := int32(4)
		cooldownTime := int32(5)
		desiredAmount := int32(6)
		loadBalancerId := value_object.NewGeneratedUuid()

		sdkAutoScalingGroup := publicCloud.NewAutoScalingGroupDetails(
			instanceId,
			"MANUAL",
			"SCALING",
			*publicCloud.NewNullableInt32(&desiredAmount),
			"region",
			"reference",
			createdAt,
			updatedAt,
			*publicCloud.NewNullableTime(&startsAt),
			*publicCloud.NewNullableTime(&endsAt),
			*publicCloud.NewNullableInt32(&minimumAmount),
			*publicCloud.NewNullableInt32(&maximumAmount),
			*publicCloud.NewNullableInt32(&cpuThreshold),
			*publicCloud.NewNullableInt32(&warmupTime),
			*publicCloud.NewNullableInt32(&cooldownTime),
			*publicCloud.NewNullableLoadBalancer(nil),
		)

		got, err := convertAutoScalingGroup(
			*sdkAutoScalingGroup,
			&domain.LoadBalancer{Id: loadBalancerId},
		)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id.String())
		assert.Equal(t, enum.AutoScalingCpuTypeManual, got.Type)
		assert.Equal(t, enum.AutoScalingGroupStateScaling, got.State)
		assert.Equal(t, 6, *got.DesiredAmount)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "reference", got.Reference.String())
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, updatedAt, got.UpdatedAt)
		assert.Equal(t, startsAt, *got.StartsAt)
		assert.Equal(t, endsAt, *got.EndsAt)
		assert.Equal(t, 1, *got.MinimumAmount)
		assert.Equal(t, 2, *got.MaximumAmount)
		assert.Equal(t, 3, *got.CpuThreshold)
		assert.Equal(t, 4, *got.WarmupTime)
		assert.Equal(t, 5, *got.CooldownTime)
		assert.Equal(t, loadBalancerId, got.LoadBalancer.Id)
	})

	t.Run(
		"returns error when loadBalancer is set but loadBalancerDetails is not passed",
		func(t *testing.T) {
			loadBalancerId := value_object.NewGeneratedUuid()
			autoScalingGroupId := value_object.NewGeneratedUuid()

			_, err := convertAutoScalingGroup(
				publicCloud.AutoScalingGroupDetails{
					Id:           autoScalingGroupId.String(),
					LoadBalancer: *publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{Id: loadBalancerId.String()}),
				},
				nil,
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, loadBalancerId.String())
			assert.ErrorContains(t, err, autoScalingGroupId.String())
		},
	)

	t.Run("invalid id returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroupDetails{Id: "tralala"},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroupDetails{Id: instanceId, Type: "tralala"},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroupDetails{
				Id:    instanceId,
				Type:  publicCloud.AUTOSCALINGGROUPTYPE_CPU_BASED,
				State: "tralala",
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid reference returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroupDetails{
				Id:        instanceId,
				Type:      "MANUAL",
				State:     "SCALING",
				Reference: "........................................................................................................................................................................................................................................................................",
			},
			nil,
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "characters long")
	})
}

func Test_convertNullableStringToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := "value"

		got := convertNullableStringToValue(*publicCloud.NewNullableString(&val))

		assert.Equal(t, "value", *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := convertNullableStringToValue(*publicCloud.NewNullableString(nil))

		assert.Nil(t, got)
	})
}

func Test_convertNullableTimeToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := time.Now()

		got := convertNullableTimeToValue(*publicCloud.NewNullableTime(&val))

		assert.Equal(t, val, *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := convertNullableTimeToValue(*publicCloud.NewNullableTime(nil))

		assert.Nil(t, got)
	})
}

func Test_convertNullableInt32ToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := int32(2)

		got := convertNullableInt32ToValue(*publicCloud.NewNullableInt32(&val))

		assert.Equal(t, int(val), *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := convertNullableInt32ToValue(*publicCloud.NewNullableInt32(nil))

		assert.Nil(t, got)
	})
}

func Test_convertStickySession(t *testing.T) {
	got := convertStickySession(publicCloud.StickySession{
		Enabled:     false,
		MaxLifeTime: 20,
	})

	assert.False(t, got.Enabled)
	assert.Equal(t, 20, got.MaxLifeTime)

}

func Test_convertHealthCheck(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		host := "host"

		sdkHealthCheck := publicCloud.NewHealthCheck(
			"GET",
			"uri",
			*publicCloud.NewNullableString(&host),
			22,
		)

		got, err := convertHealthCheck(*sdkHealthCheck)

		assert.NoError(t, err)
		assert.Equal(t, enum.MethodGet, got.Method)
		assert.Equal(t, "uri", got.Uri)
		assert.Equal(t, "host", *got.Host)
		assert.Equal(t, 22, got.Port)
	})

	t.Run("invalid method returns error", func(t *testing.T) {
		_, err := convertHealthCheck(publicCloud.HealthCheck{Method: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

}

func Test_convertLoadBalancerConfiguration(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkLoadBalancerConfiguration := publicCloud.NewLoadBalancerConfiguration(
			*publicCloud.NewNullableStickySession(&publicCloud.StickySession{MaxLifeTime: 44}),
			"roundrobin",
			*publicCloud.NewNullableHealthCheck(&publicCloud.HealthCheck{Method: "GET"}),
			true, 1, 2)

		got, err := convertLoadBalancerConfiguration(*sdkLoadBalancerConfiguration)

		assert.NoError(t, err)
		assert.Equal(t, 44, got.StickySession.MaxLifeTime)
		assert.Equal(t, enum.BalanceRoundRobin, got.Balance)
		assert.Equal(t, enum.MethodGet, got.HealthCheck.Method)
		assert.True(t, got.XForwardedFor)
	})

	t.Run("invalid balance returns error", func(t *testing.T) {
		_, err := convertLoadBalancerConfiguration(publicCloud.LoadBalancerConfiguration{Balance: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid HealthCheck returns error", func(t *testing.T) {
		_, err := convertLoadBalancerConfiguration(publicCloud.LoadBalancerConfiguration{
			Balance:     publicCloud.BALANCE_ROUNDROBIN,
			HealthCheck: *publicCloud.NewNullableHealthCheck(&publicCloud.HealthCheck{Method: "tralala"}),
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertLoadBalancer(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		reference := "reference"
		startedAt := time.Now()

		sdkLoadBalancer := publicCloud.NewLoadBalancerDetails(
			instanceId,
			"lsw.m3.large",
			publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "unit"}},
			*publicCloud.NewNullableString(&reference),
			"CREATING",
			*publicCloud.NewNullableTime(&startedAt),
			[]publicCloud.IpDetails{{
				Ip:          "1.2.3.4",
				NetworkType: publicCloud.NETWORKTYPE_PUBLIC,
			}},
			"region",
			*publicCloud.NewNullableLoadBalancerConfiguration(&publicCloud.LoadBalancerConfiguration{
				TargetPort: 22,
				Balance:    "roundrobin",
			}),
			*publicCloud.NewNullableAutoScalingGroup(nil),
			*publicCloud.NewNullablePrivateNetwork(&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"}),
			publicCloud.Contract{
				BillingFrequency: 1,
				Type:             publicCloud.CONTRACTTYPE_MONTHLY,
				State:            publicCloud.CONTRACTSTATE_ACTIVE,
				Term:             publicCloud.CONTRACTTERM__1,
			},
		)

		got, err := convertLoadBalancer(*sdkLoadBalancer)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id.String())
		assert.Equal(t, "lsw.m3.large", got.Type)
		assert.Equal(t, "unit", got.Resources.Cpu.Unit)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, enum.StateCreating, got.State)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(t, 22, got.Configuration.TargetPort)
		assert.Equal(t, "privateNetworkId", got.PrivateNetwork.Id)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		_, err := convertLoadBalancer(publicCloud.LoadBalancerDetails{Id: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := convertLoadBalancer(publicCloud.LoadBalancerDetails{
			Id:    instanceId,
			State: "tralala",
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		_, err := convertLoadBalancer(publicCloud.LoadBalancerDetails{
			Id:       instanceId,
			State:    publicCloud.STATE_RUNNING,
			Contract: publicCloud.Contract{BillingFrequency: 55},
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid ips returns error", func(t *testing.T) {
		_, err := convertLoadBalancer(publicCloud.LoadBalancerDetails{
			Id:    instanceId,
			State: publicCloud.STATE_RUNNING,
			Contract: publicCloud.Contract{
				Type:  publicCloud.CONTRACTTYPE_MONTHLY,
				State: publicCloud.CONTRACTSTATE_ACTIVE,
				Term:  publicCloud.CONTRACTTERM__1,
			},
			Ips: []publicCloud.IpDetails{{NetworkType: "tralala"}},
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid configuration returns error", func(t *testing.T) {
		_, err := convertLoadBalancer(publicCloud.LoadBalancerDetails{
			Id:    instanceId,
			State: publicCloud.STATE_RUNNING,
			Contract: publicCloud.Contract{
				Type:  publicCloud.CONTRACTTYPE_MONTHLY,
				State: publicCloud.CONTRACTSTATE_ACTIVE,
				Term:  publicCloud.CONTRACTTERM__1,
			},
			Configuration: *publicCloud.NewNullableLoadBalancerConfiguration(&publicCloud.LoadBalancerConfiguration{Balance: "tralala"}),
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertEntityToLaunchInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {

		_, err := convertEntityToLaunchInstanceOpts(
			domain.Instance{Type: "tralala"},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		_, err := convertEntityToLaunchInstanceOpts(
			domain.Instance{Type: "lsw.m3.large", RootDiskStorageType: "tralala"},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid imageId returns error", func(t *testing.T) {
		_, err := convertEntityToLaunchInstanceOpts(
			domain.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               domain.Image{Id: "tralala"},
			},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {
		_, err := convertEntityToLaunchInstanceOpts(
			domain.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               domain.Image{Id: enum.Ubuntu200464Bit},
				Contract:            domain.Contract{Type: "tralala"},
			},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {
		_, err := convertEntityToLaunchInstanceOpts(
			domain.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               domain.Image{Id: enum.Ubuntu200464Bit},
				Contract: domain.Contract{
					Type: enum.ContractTypeMonthly,
					Term: 55,
				},
			},
		)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {
		_, err := convertEntityToLaunchInstanceOpts(
			domain.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               domain.Image{Id: enum.Ubuntu200464Bit},
				Contract: domain.Contract{
					Type:             enum.ContractTypeMonthly,
					Term:             enum.ContractTermThree,
					BillingFrequency: 55,
				},
			},
		)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("required values are set", func(t *testing.T) {
		instance := domain.NewCreateInstance(
			"region",
			"lsw.c3.4xlarge",
			enum.RootDiskStorageTypeCentral,
			enum.Almalinux864Bit,
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			domain.OptionalCreateInstanceValues{},
		)

		got, err := convertEntityToLaunchInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, publicCloud.INSTANCETYPENAME_C3_4XLARGE, got.Type)
		assert.Equal(
			t,
			publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
			got.RootDiskStorageType,
		)
		assert.Equal(t, publicCloud.IMAGEID_ALMALINUX_8_64_BIT, got.ImageId)
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
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		instance := domain.NewCreateInstance(
			"",
			"lsw.c3.4xlarge",
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

		got, err := convertEntityToLaunchInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, marketAppId, *got.MarketAppId)
		assert.Equal(t, reference, *got.Reference)
		assert.Equal(t, sshKey, *got.SshKey)
	})
}

func Test_convertEntityToUpdateInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {

		_, err := convertEntityToUpdateInstanceOpts(
			domain.Instance{Type: "tralala"},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {

		_, err := convertEntityToUpdateInstanceOpts(
			domain.Instance{Contract: domain.Contract{Type: "tralala"}},
		)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {

		_, err := convertEntityToUpdateInstanceOpts(
			domain.Instance{Contract: domain.Contract{Term: 55}},
		)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {

		_, err := convertEntityToUpdateInstanceOpts(
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

		got, err := convertEntityToUpdateInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.INSTANCETYPENAME_C3_LARGE, got.GetType())
		assert.Equal(t, "reference", got.GetReference())
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.GetContractType())
		assert.Equal(t, publicCloud.CONTRACTTERM__3, got.GetContractTerm())
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__6, got.GetBillingFrequency())
		assert.Equal(t, int32(23), got.GetRootDiskSize())
	})
}

func Test_convertInstanceType(t *testing.T) {
	got := convertInstanceType(publicCloud.InstanceType{Name: "name"})
	want := domain.InstanceType{Name: "name"}

	assert.Equal(t, want, got)
}

func Test_convertRegion(t *testing.T) {
	got := convertRegion(publicCloud.Region{Name: "name", Location: "location"})
	want := domain.Region{Name: "name", Location: "location"}

	assert.Equal(t, want, got)
}
