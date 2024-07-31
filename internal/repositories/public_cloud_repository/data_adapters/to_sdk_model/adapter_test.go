package to_sdk_model

import (
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

var autoScalingGroupId = "90b9f2cc-c655-40ea-b01a-58c00e175c96"
var instanceId = "5d7f8262-d77f-4476-8da8-6a84f8f2ae8d"

func Test_adaptImageDetails(t *testing.T) {
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

		got, err := adaptImageDetails(*sdkImage)

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

		_, err := adaptImageDetails(sdkImage)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptNetworkSpeed(t *testing.T) {
	sdkNetworkSpeed := publicCloud.NewNetworkSpeed(1, "unit")
	got := adaptNetworkSpeed(*sdkNetworkSpeed)

	assert.Equal(t, 1, got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_adaptMemory(t *testing.T) {
	sdkMemory := publicCloud.NewMemory(1, "unit")
	got := adaptMemory(*sdkMemory)

	assert.Equal(t, float64(1), got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_adaptCpu(t *testing.T) {
	sdkCpu := publicCloud.NewCpu(1, "unit")
	got := adaptCpu(*sdkCpu)

	assert.Equal(t, 1, got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_adaptResources(t *testing.T) {
	sdkResources := publicCloud.NewResources(
		publicCloud.Cpu{Unit: "cpu"},
		publicCloud.Memory{Unit: "memory"},
		publicCloud.NetworkSpeed{Unit: "publicNetworkSpeed"},
		publicCloud.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	got := adaptResources(*sdkResources)

	assert.Equal(t, "cpu", got.Cpu.Unit)
	assert.Equal(t, "memory", got.Memory.Unit)
	assert.Equal(t, "publicNetworkSpeed", got.PublicNetworkSpeed.Unit)
	assert.Equal(t, "privateNetworkSpeed", got.PrivateNetworkSpeed.Unit)
}

func TestAdaptInstanceDetails(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		startedAt := time.Now()
		sdkInstance := generateInstanceDetails(t, &startedAt)

		got, err := AdaptInstanceDetails(sdkInstance)

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

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid Image returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Image.Id = "tralala"

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.State = "tralala"

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.RootDiskSize = 5000

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.RootDiskStorageType = "tralala"

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid ip returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Ips = []publicCloud.IpDetails{{NetworkType: "tralala"}}

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.Contract.BillingFrequency = 55

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid autoScalingGroup returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t, nil)
		sdkInstance.AutoScalingGroup.Get().Id = "tralala"

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptDdos(t *testing.T) {
	got := adaptDdos(publicCloud.Ddos{
		DetectionProfile: "detectionProfile",
		ProtectionType:   "protectionType",
	})

	assert.Equal(t, "detectionProfile", got.DetectionProfile)
	assert.Equal(t, "protectionType", got.ProtectionType)
}

func Test_adaptIpDetails(t *testing.T) {
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

		got, err := adaptIpDetails(*sdkIp)

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
		_, err := adaptIpDetails(publicCloud.IpDetails{NetworkType: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptIpsDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := adaptIpsDetails([]publicCloud.IpDetails{{
			Ip:          "1.2.3.4",
			NetworkType: publicCloud.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := adaptIps([]publicCloud.Ip{{NetworkType: "tralala"}})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptIps(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := adaptIps([]publicCloud.Ip{{
			Ip:          "1.2.3.4",
			NetworkType: publicCloud.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := adaptIps([]publicCloud.Ip{{NetworkType: "tralala"}})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptContract(t *testing.T) {
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

		got, err := adaptContract(*sdkContract)

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

		_, err := adaptContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "45")
	})

	t.Run("error returned for invalid term", func(t *testing.T) {
		sdkContract := publicCloud.Contract{BillingFrequency: 0, Term: 55}

		_, err := adaptContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("error returned for invalid type", func(t *testing.T) {
		sdkContract := publicCloud.Contract{
			BillingFrequency: 0,
			Term:             0,
			Type:             "tralala",
		}

		_, err := adaptContract(sdkContract)

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

		_, err := adaptContract(sdkContract)

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

			_, err := adaptContract(sdkContract)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "contract.term cannot be 0")
		},
	)

}

func Test_adaptIso(t *testing.T) {
	got := adaptIso(*publicCloud.NewIso("id", "name"))

	assert.Equal(t, "id", got.Id)
	assert.Equal(t, "name", got.Name)
}

func Test_adaptPrivateNetwork(t *testing.T) {
	got := adaptPrivateNetwork(*publicCloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	))

	assert.Equal(t, "id", got.Id)
	assert.Equal(t, "status", got.Status)
	assert.Equal(t, "subnet", got.Subnet)
}

func Test_adaptAutoScalingGroupDetails(t *testing.T) {
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

		got, err := AdaptAutoScalingGroupDetails(*sdkAutoScalingGroup)

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
		_, err := AdaptAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{Id: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := AdaptAutoScalingGroupDetails(
			publicCloud.AutoScalingGroupDetails{Id: instanceId, Type: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := AdaptAutoScalingGroupDetails(
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
		_, err := AdaptAutoScalingGroupDetails(
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
		_, err := AdaptAutoScalingGroupDetails(
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

func Test_adaptStickySession(t *testing.T) {
	got := adaptStickySession(publicCloud.StickySession{
		Enabled:     false,
		MaxLifeTime: 20,
	})

	assert.False(t, got.Enabled)
	assert.Equal(t, 20, got.MaxLifeTime)

}

func Test_adaptHealthCheck(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		host := "host"

		sdkHealthCheck := publicCloud.NewHealthCheck(
			"GET",
			"uri",
			*publicCloud.NewNullableString(&host),
			22,
		)

		got, err := adaptHealthCheck(*sdkHealthCheck)

		assert.NoError(t, err)
		assert.Equal(t, enum.MethodGet, got.Method)
		assert.Equal(t, "uri", got.Uri)
		assert.Equal(t, "host", *got.Host)
		assert.Equal(t, 22, got.Port)
	})

	t.Run("invalid method returns error", func(t *testing.T) {
		_, err := adaptHealthCheck(publicCloud.HealthCheck{Method: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

}

func Test_adaptLoadBalancerConfiguration(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkLoadBalancerConfiguration := publicCloud.NewLoadBalancerConfiguration(
			*publicCloud.NewNullableStickySession(&publicCloud.StickySession{MaxLifeTime: 44}),
			"roundrobin",
			*publicCloud.NewNullableHealthCheck(&publicCloud.HealthCheck{Method: "GET"}),
			true, 1, 2)

		got, err := adaptLoadBalancerConfiguration(*sdkLoadBalancerConfiguration)

		assert.NoError(t, err)
		assert.Equal(t, 44, got.StickySession.MaxLifeTime)
		assert.Equal(t, enum.BalanceRoundRobin, got.Balance)
		assert.Equal(t, enum.MethodGet, got.HealthCheck.Method)
		assert.True(t, got.XForwardedFor)
	})

	t.Run("invalid balance returns error", func(t *testing.T) {
		_, err := adaptLoadBalancerConfiguration(
			publicCloud.LoadBalancerConfiguration{Balance: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid HealthCheck returns error", func(t *testing.T) {
		_, err := adaptLoadBalancerConfiguration(
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

func TestAdaptLoadBalancerDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()
		sdkLoadBalancer := generateLoadBalancerDetails(&startedAt)

		got, err := AdaptLoadBalancerDetails(sdkLoadBalancer)

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

		_, err := AdaptLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.State = "tralala"

		_, err := AdaptLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Contract.BillingFrequency = 55

		_, err := AdaptLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid ips returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Ips = []publicCloud.IpDetails{{NetworkType: "tralala"}}

		_, err := AdaptLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid configuration returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancerDetails(nil)
		sdkLoadBalancer.Configuration.Get().Balance = "tralala"

		_, err := AdaptLoadBalancerDetails(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestAdaptInstanceType(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		got, err := AdaptInstanceType(
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
		got, err := AdaptInstanceType(
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
		_, err := AdaptInstanceType(
			publicCloud.InstanceType{
				StorageTypes: []publicCloud.RootDiskStorageType{"tralala"},
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestAdaptRegion(t *testing.T) {
	got := AdaptRegion(publicCloud.Region{Name: "name", Location: "location"})
	want := domain.Region{Name: "name", Location: "location"}

	assert.Equal(t, want, got)
}

func TestAdaptInstance(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()

		sdkInstance := generateInstance(t, &startedAt)

		got, err := AdaptInstance(sdkInstance)

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

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid Image returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Id = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.State = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.RootDiskSize = 5000

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.RootDiskStorageType = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid ip returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Ips = []publicCloud.Ip{{NetworkType: "tralala"}}

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.Contract.BillingFrequency = 55

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid autoScalingGroup returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t, nil)
		sdkInstance.AutoScalingGroup.Get().Id = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptImage(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkImage := publicCloud.NewImage(
			publicCloud.IMAGEID_UBUNTU_24_04_64_BIT,
			"name",
			"version",
			"family",
			"flavour",
			"architecture",
		)

		got, err := adaptImage(*sdkImage)

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

		_, err := adaptImage(sdkImage)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptIp(t *testing.T) {
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

		got, err := adaptIp(*sdkIp)

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
		_, err := adaptIpDetails(publicCloud.IpDetails{NetworkType: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestAdaptAutoScalingGroup(t *testing.T) {
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

		got, err := adaptAutoScalingGroup(*sdkAutoScalingGroup)

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
		_, err := adaptAutoScalingGroup(
			publicCloud.AutoScalingGroup{Id: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := adaptAutoScalingGroup(
			publicCloud.AutoScalingGroup{Id: instanceId, Type: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := adaptAutoScalingGroup(
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
		_, err := adaptAutoScalingGroup(
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

func Test_adaptLoadBalancer(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()
		sdkLoadBalancer := generateLoadBalancer(&startedAt)

		got, err := adaptLoadBalancer(sdkLoadBalancer)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id.String())
		assert.Equal(t, string(publicCloud.TYPENAME_M3_LARGE), got.Type.String())
		assert.Equal(t, "unit", got.Resources.Cpu.Unit)
		assert.Equal(t, enum.StateCreating, got.State)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancer(nil)
		sdkLoadBalancer.Id = "tralala"

		_, err := adaptLoadBalancer(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkLoadBalancer := generateLoadBalancer(nil)
		sdkLoadBalancer.State = "tralala"

		_, err := adaptLoadBalancer(sdkLoadBalancer)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptStorageTypes(t *testing.T) {
	t.Run("sdk storageTypes are adapted correctly", func(t *testing.T) {
		sdkStorageTypes := []publicCloud.RootDiskStorageType{
			publicCloud.ROOTDISKSTORAGETYPE_CENTRAL,
			publicCloud.ROOTDISKSTORAGETYPE_LOCAL,
		}
		got, err := adaptStorageTypes(sdkStorageTypes)
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
			got, err := adaptStorageTypes(sdkStorageTypes)

			assert.Nil(t, got)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)
}

func Test_adaptPrice(t *testing.T) {
	sdkPrice := publicCloud.NewPrice("1", "2")
	got := adaptPrice(*sdkPrice)

	want := domain.Price{
		HourlyPrice:  "1",
		MonthlyPrice: "2",
	}

	assert.Equal(t, want, got)
}

func Test_adaptStorage(t *testing.T) {
	sdkStorage := publicCloud.NewStorage(
		publicCloud.Price{HourlyPrice: "1"},
		publicCloud.Price{HourlyPrice: "2"},
	)
	got := adaptStorage(*sdkStorage)

	want := domain.Storage{
		Local:   domain.Price{HourlyPrice: "1"},
		Central: domain.Price{HourlyPrice: "2"},
	}

	assert.Equal(t, want, got)
}

func Test_adaptPrices(t *testing.T) {
	sdkPrices := publicCloud.NewPrices(
		"currency",
		"symbol",
		publicCloud.Price{HourlyPrice: "1"},
		publicCloud.Storage{Central: publicCloud.Price{HourlyPrice: "2"}},
	)
	got := adaptPrices(*sdkPrices)

	want := domain.Prices{
		Currency:       "currency",
		CurrencySymbol: "symbol",
		Compute:        domain.Price{HourlyPrice: "1"},
		Storage:        domain.Storage{Central: domain.Price{HourlyPrice: "2"}},
	}

	assert.Equal(t, want, got)
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
