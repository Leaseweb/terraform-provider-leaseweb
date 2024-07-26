package public_cloud_repository

import (
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

var instanceId = "5d7f8262-d77f-4476-8da8-6a84f8f2ae8d"
var autoScalingGroupId = "90b9f2cc-c655-40ea-b01a-58c00e175c96"

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func Test_convertImageDetails(t *testing.T) {
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

		got, err := convertImageDetails(*sdkImage)

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

		_, err := convertImageDetails(sdkImage)
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

func Test_convertInstanceDetails(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		startedAt := time.Now()
		sdkInstance := generateInstanceDetails(t, &startedAt)

		got, err := convertInstanceDetails(sdkInstance)

		assert.NoError(t, err)
		assert.Equal(
			t,
			"5d7f8262-d77f-4476-8da8-6a84f8f2ae8d",
			got.Id.String(),
		)
		assert.Equal(t, string(publicCloud.TYPENAME_M3_LARGE), got.Type.String())
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
		assert.Equal(t, autoScalingGroupId, got.AutoScalingGroup.Id.String())
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Id = "tralala"

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid Image returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Image.Id = "tralala"

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.State = "tralala"

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.RootDiskSize = 5000

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.RootDiskStorageType = "tralala"

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid ip returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Ips = []publicCloud.IpDetails{{NetworkType: "tralala"}}

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Contract.BillingFrequency = 55

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid autoScalingGroup returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.AutoScalingGroup.Get().Id = "tralala"

		_, err := convertInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func generateInstanceDetails(
	t *testing.T,
	startedAt *time.Time,
) publicCloud.InstanceDetails {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	return *publicCloud.NewInstanceDetails(
		instanceId,
		publicCloud.TYPENAME_M3_LARGE,
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
		publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
		publicCloud.Contract{
			BillingFrequency: 1,
			Type:             publicCloud.CONTRACTTYPE_HOURLY,
			State:            publicCloud.CONTRACTSTATE_ACTIVE,
		},
		*publicCloud.NewNullableIso(&publicCloud.Iso{Id: "isoId"}),
		*publicCloud.NewNullablePrivateNetwork(
			&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"},
		),
		publicCloud.ImageDetails{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
		[]publicCloud.IpDetails{
			{Ip: "1.2.3.4", NetworkType: publicCloud.NETWORKTYPE_PUBLIC},
		},
		*publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
			Id:    autoScalingGroupId,
			Type:  publicCloud.AUTOSCALINGGROUPTYPE_CPU_BASED,
			State: publicCloud.AUTOSCALINGGROUPSTATE_ACTIVE,
		}),
	)
}

func generateInstance(
	t *testing.T,
	startedAt *time.Time,
) publicCloud.Instance {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	return *publicCloud.NewInstance(
		instanceId,
		publicCloud.TYPENAME_M3_LARGE,
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
		publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
		publicCloud.Contract{
			BillingFrequency: 1,
			Type:             publicCloud.CONTRACTTYPE_HOURLY,
			State:            publicCloud.CONTRACTSTATE_ACTIVE,
		},
		publicCloud.Image{Id: publicCloud.IMAGEID_CENTOS_7_64_BIT},
		[]publicCloud.Ip{
			{Ip: "1.2.3.4", NetworkType: publicCloud.NETWORKTYPE_PUBLIC},
		},
		*publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
			Id:    autoScalingGroupId,
			Type:  publicCloud.AUTOSCALINGGROUPTYPE_CPU_BASED,
			State: publicCloud.AUTOSCALINGGROUPSTATE_ACTIVE,
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

func Test_convertIpDetails(t *testing.T) {
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
			*publicCloud.NewNullableDdos(
				&publicCloud.Ddos{DetectionProfile: "detectionProfile"},
			),
		)

		got, err := convertIpDetails(*sdkIp)

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
		_, err := convertIpDetails(publicCloud.IpDetails{NetworkType: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertIpsDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := convertIpsDetails([]publicCloud.IpDetails{{
			Ip:          "1.2.3.4",
			NetworkType: publicCloud.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := convertIps([]publicCloud.Ip{{NetworkType: "tralala"}})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertIps(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := convertIps([]publicCloud.Ip{{
			Ip:          "1.2.3.4",
			NetworkType: publicCloud.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := convertIps([]publicCloud.Ip{{NetworkType: "tralala"}})

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

func Test_convertAutoScalingGroupDetails(t *testing.T) {
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
			*publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{
				Id:    loadBalancerId.String(),
				State: publicCloud.STATE_CREATING,
				Type:  publicCloud.TYPENAME_M3_LARGE,
			}),
		)

		got, err := convertAutoScalingGroupDetails(*sdkAutoScalingGroup)

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

	t.Run("invalid id returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{Id: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{Id: instanceId, Type: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{
				Id:    instanceId,
				Type:  publicCloud.AUTOSCALINGGROUPTYPE_CPU_BASED,
				State: "tralala",
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid reference returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{
				Id:        instanceId,
				Type:      "MANUAL",
				State:     "SCALING",
				Reference: "........................................................................................................................................................................................................................................................................",
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "characters long")
	})

	t.Run("invalid loadBalancer returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{
				Id:        instanceId,
				Type:      "MANUAL",
				State:     "SCALING",
				Reference: "reference",
				LoadBalancer: *publicCloud.NewNullableLoadBalancer(
					&publicCloud.LoadBalancer{Id: "tralala"},
				),
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
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
		_, err := convertLoadBalancerConfiguration(
			publicCloud.LoadBalancerConfiguration{Balance: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid HealthCheck returns error", func(t *testing.T) {
		_, err := convertLoadBalancerConfiguration(
			publicCloud.LoadBalancerConfiguration{
				Balance: publicCloud.BALANCE_ROUNDROBIN,
				HealthCheck: *publicCloud.NewNullableHealthCheck(
					&publicCloud.HealthCheck{Method: "tralala"},
				),
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertLoadBalancerDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()
		sdkLoadBalancer := generateLoadBalancerDetails(&startedAt)

		got, err := convertLoadBalancerDetails(sdkLoadBalancer)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id.String())
		assert.Equal(t, string(publicCloud.TYPENAME_M3_LARGE), got.Type.String())
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
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Id = "tralala"

		_, err := convertLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.State = "tralala"

		_, err := convertLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Contract.BillingFrequency = 55

		_, err := convertLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid ips returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Ips = []publicCloud.IpDetails{{NetworkType: "tralala"}}

		_, err := convertLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid configuration returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Configuration.Get().Balance = "tralala"

		_, err := convertLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertEntityToLaunchInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Type = value_object.InstanceType{Type: "tralala"}

		_, err := convertEntityToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.RootDiskStorageType = "tralala"

		_, err := convertEntityToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid imageId returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Image.Id = "tralala"

		_, err := convertEntityToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Contract.Type = "tralala"

		_, err := convertEntityToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Contract.Term = 55

		_, err := convertEntityToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Contract.BillingFrequency = 55

		_, err := convertEntityToLaunchInstanceOpts(instance)

		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		instance := generateDomainInstance()
		instance.Type = value_object.InstanceType{Type: "tralala"}

		_, err := convertEntityToLaunchInstanceOpts(instance)

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

		got, err := convertEntityToLaunchInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, publicCloud.TYPENAME_C3_4XLARGE, got.Type)
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

		got, err := convertEntityToLaunchInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, marketAppId, *got.MarketAppId)
		assert.Equal(t, reference, *got.Reference)
		assert.Equal(t, defaultSshKey, *got.SshKey)
	})
}

func Test_convertEntityToUpdateInstanceOpts(t *testing.T) {
	t.Run("invalid instanceType returns error", func(t *testing.T) {

		_, err := convertEntityToUpdateInstanceOpts(
			domain.Instance{Type: value_object.InstanceType{Type: "tralala"}},
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

		got, err := convertEntityToUpdateInstanceOpts(instance)

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.TYPENAME_C3_LARGE, got.GetType())
		assert.Equal(t, "reference", got.GetReference())
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.GetContractType())
		assert.Equal(t, publicCloud.CONTRACTTERM__3, got.GetContractTerm())
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__6, got.GetBillingFrequency())
		assert.Equal(t, int32(23), got.GetRootDiskSize())
	})
}

func Test_convertInstanceType(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		got, err := convertInstanceType(
			publicCloud.InstanceType{
				Name:      "name",
				Resources: publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "cpu"}},
				Prices:    publicCloud.Prices{Currency: "currency"},
			},
		)
		want := domain.InstanceType{
			Name:      "name",
			Resources: domain.Resources{Cpu: domain.Cpu{Unit: "cpu"}},
			Prices:    domain.Prices{Currency: "currency"},
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("optional values are set", func(t *testing.T) {
		got, err := convertInstanceType(
			publicCloud.InstanceType{
				StorageTypes: []publicCloud.RootDiskStorageType{
					publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
				},
			},
		)
		want := domain.InstanceType{
			StorageTypes: &domain.StorageTypes{enum.RootDiskStorageTypeCentral},
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("invalid storageType returns an error", func(t *testing.T) {
		_, err := convertInstanceType(
			publicCloud.InstanceType{
				StorageTypes: []publicCloud.RootDiskStorageType{"tralala"},
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertRegion(t *testing.T) {
	got := convertRegion(publicCloud.Region{Name: "name", Location: "location"})
	want := domain.Region{Name: "name", Location: "location"}

	assert.Equal(t, want, got)
}

func Test_convertInstance(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()

		sdkInstance := generateInstance(t, &startedAt)

		got, err := convertInstance(sdkInstance)

		assert.NoError(t, err)
		assert.Equal(
			t,
			instanceId,
			got.Id.String(),
		)
		assert.Equal(t, "lsw.m3.large", got.Type.String())
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
		assert.Equal(t, enum.Centos764Bit, got.Image.Id)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(t, autoScalingGroupId, got.AutoScalingGroup.Id.String())
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Id = "tralala"

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid Image returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Id = "tralala"

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.State = "tralala"

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.RootDiskSize = 5000

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.RootDiskStorageType = "tralala"

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid ip returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Ips = []publicCloud.Ip{{NetworkType: "tralala"}}

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Contract.BillingFrequency = 55

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid autoScalingGroup returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.AutoScalingGroup.Get().Id = "tralala"

		_, err := convertInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertImage(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkImage := publicCloud.NewImage(
			publicCloud.IMAGEID_UBUNTU_24_04_64_BIT,
			"name",
			"version",
			"family",
			"flavour",
			"architecture",
		)

		got, err := convertImage(*sdkImage)

		assert.Nil(t, err)
		assert.Equal(t, enum.Ubuntu240464Bit, got.Id)
		assert.Equal(t, "name", got.Name)
		assert.Equal(t, "version", got.Version)
		assert.Equal(t, "family", got.Family)
		assert.Equal(t, "flavour", got.Flavour)
		assert.Equal(t, "architecture", got.Architecture)
	})

	t.Run("invalid imageId returns error", func(t *testing.T) {
		sdkImage := publicCloud.Image{Id: "tralala"}

		_, err := convertImage(sdkImage)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_convertIp(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		reverseLookup := "reverseLookup"

		sdkIp := publicCloud.NewIp(
			"1.2.3.4",
			"prefixLength",
			5,
			true,
			false,
			publicCloud.NETWORKTYPE_INTERNAL,
			*publicCloud.NewNullableString(&reverseLookup),
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
	})

	t.Run("error returned for invalid networkType", func(t *testing.T) {
		_, err := convertIpDetails(publicCloud.IpDetails{NetworkType: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
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

		sdkAutoScalingGroup := publicCloud.NewAutoScalingGroup(
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
		)

		got, err := convertAutoScalingGroup(*sdkAutoScalingGroup)

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
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroup{Id: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroup{Id: instanceId, Type: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroup{
				Id:    instanceId,
				Type:  publicCloud.AUTOSCALINGGROUPTYPE_CPU_BASED,
				State: "tralala",
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid reference returns error", func(t *testing.T) {
		_, err := convertAutoScalingGroup(
			publicCloud.AutoScalingGroup{
				Id:        instanceId,
				Type:      "MANUAL",
				State:     "SCALING",
				Reference: "........................................................................................................................................................................................................................................................................",
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "characters long")
	})
}

func Test_convertLoadBalancer(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()
		sdkLoadBalancer := generateLoadBalancer(&startedAt)

		got, err := convertLoadBalancer(sdkLoadBalancer)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id.String())
		assert.Equal(t, string(publicCloud.TYPENAME_M3_LARGE), got.Type.String())
		assert.Equal(t, "unit", got.Resources.Cpu.Unit)
		assert.Equal(t, enum.StateCreating, got.State)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancer(nil)
		sdkLoadBalancer.Id = "tralala"

		_, err := convertLoadBalancer(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancer(nil)
		sdkLoadBalancer.State = "tralala"

		_, err := convertLoadBalancer(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func generateLoadBalancerDetails(startedAt *time.Time) publicCloud.LoadBalancerDetails {
	reference := "reference"

	return *publicCloud.NewLoadBalancerDetails(
		instanceId,
		"lsw.m3.large",
		publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "unit"}},
		*publicCloud.NewNullableString(&reference),
		"CREATING",
		*publicCloud.NewNullableTime(startedAt),
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
		*publicCloud.NewNullablePrivateNetwork(
			&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"},
		),
		publicCloud.Contract{
			BillingFrequency: 1,
			Type:             publicCloud.CONTRACTTYPE_MONTHLY,
			State:            publicCloud.CONTRACTSTATE_ACTIVE,
			Term:             publicCloud.CONTRACTTERM__1,
		},
	)
}

func generateLoadBalancer(startedAt *time.Time) publicCloud.LoadBalancer {
	reference := "reference"

	return *publicCloud.NewLoadBalancer(
		instanceId,
		"lsw.m3.large",
		publicCloud.Resources{Cpu: publicCloud.Cpu{Unit: "unit"}},
		*publicCloud.NewNullableString(&reference),
		publicCloud.STATE_CREATING,
		*publicCloud.NewNullableTime(startedAt),
	)
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

func Test_convertStorageTypes(t *testing.T) {
	t.Run("sdk storageTypes are converted correctly", func(t *testing.T) {
		sdkStorageTypes := []publicCloud.RootDiskStorageType{
			publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
			publicCloud.ROOTDISKSTORAGETYPE_LOCAL,
		}
		got, err := convertStorageTypes(sdkStorageTypes)
		want := domain.StorageTypes{
			enum.RootDiskStorageTypeCentral, enum.RootDiskStorageTypeLocal,
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run(
		"error bubbles up when local storageType cannot be created",
		func(t *testing.T) {
			sdkStorageTypes := []publicCloud.RootDiskStorageType{"tralala"}
			got, err := convertStorageTypes(sdkStorageTypes)

			assert.Nil(t, got)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)
}

func Test_convertPrice(t *testing.T) {
	sdkPrice := publicCloud.NewPrice("1", "2")
	got := convertPrice(*sdkPrice)

	want := domain.Price{
		HourlyPrice:  "1",
		MonthlyPrice: "2",
	}

	assert.Equal(t, want, got)
}

func Test_convertStorage(t *testing.T) {
	sdkStorage := publicCloud.NewStorage(
		publicCloud.Price{HourlyPrice: "1"},
		publicCloud.Price{HourlyPrice: "2"},
	)
	got := convertStorage(*sdkStorage)

	want := domain.Storage{
		Local:   domain.Price{HourlyPrice: "1"},
		Central: domain.Price{HourlyPrice: "2"},
	}

	assert.Equal(t, want, got)
}

func Test_convertPrices(t *testing.T) {
	sdkPrices := publicCloud.NewPrices(
		"currency",
		"symbol",
		publicCloud.Price{HourlyPrice: "1"},
		publicCloud.Storage{Central: publicCloud.Price{HourlyPrice: "2"}},
	)
	got := convertPrices(*sdkPrices)

	want := domain.Prices{
		Currency:       "currency",
		CurrencySymbol: "symbol",
		Compute:        domain.Price{HourlyPrice: "1"},
		Storage:        domain.Storage{Central: domain.Price{HourlyPrice: "2"}},
	}

	assert.Equal(t, want, got)
}
