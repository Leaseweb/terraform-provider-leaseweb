package to_domain_entity

import (
	"testing"
	"time"

	sdk "github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/stretchr/testify/assert"
)

var autoScalingGroupId = "90b9f2cc-c655-40ea-b01a-58c00e175c96"
var instanceId = "5d7f8262-d77f-4476-8da8-6a84f8f2ae8d"

func Test_adaptNetworkSpeed(t *testing.T) {
	sdkNetworkSpeed := sdk.NewNetworkSpeed(1, "unit")
	got := adaptNetworkSpeed(*sdkNetworkSpeed)

	assert.Equal(t, 1, got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_adaptMemory(t *testing.T) {
	sdkMemory := sdk.NewMemory(1, "unit")
	got := adaptMemory(*sdkMemory)

	assert.Equal(t, float64(1), got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_adaptCpu(t *testing.T) {
	sdkCpu := sdk.NewCpu(1, "unit")
	got := adaptCpu(*sdkCpu)

	assert.Equal(t, 1, got.Value)
	assert.Equal(t, "unit", got.Unit)
}

func Test_adaptResources(t *testing.T) {
	sdkResources := sdk.NewResources(
		sdk.Cpu{Unit: "cpu"},
		sdk.Memory{Unit: "memory"},
		sdk.NetworkSpeed{Unit: "publicNetworkSpeed"},
		sdk.NetworkSpeed{Unit: "privateNetworkSpeed"},
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
			got.Id,
		)
		assert.Equal(t, string(sdk.TYPENAME_M3_LARGE), got.Type.String())
		assert.Equal(t, "cpu", got.Resources.Cpu.Unit)
		assert.Equal(t, domain.Region{Name: "region"}, got.Region)
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
		assert.Equal(t, "CENTOS_7_64BIT", got.Image.Id)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(t, autoScalingGroupId, got.AutoScalingGroup.Id)
		assert.Equal(t, "unit", got.Volume.Unit)
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
		sdkInstance.Ips = []sdk.IpDetails{{NetworkType: "tralala"}}

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
}

func Test_adaptDdos(t *testing.T) {
	got := adaptDdos(sdk.Ddos{
		DetectionProfile: "detectionProfile",
		ProtectionType:   "protectionType",
	})

	assert.Equal(t, "detectionProfile", got.DetectionProfile)
	assert.Equal(t, "protectionType", got.ProtectionType)
}

func Test_adaptIpDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		reverseLookup := "reverseLookup"

		sdkIp := sdk.NewIpDetails(
			"1.2.3.4",
			"prefixLength",
			5,
			true,
			false,
			sdk.NETWORKTYPE_INTERNAL,
			*sdk.NewNullableString(&reverseLookup),
			*sdk.NewNullableDdos(
				&sdk.Ddos{DetectionProfile: "detectionProfile"},
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
		_, err := adaptIpDetails(sdk.IpDetails{NetworkType: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptIpsDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := adaptIpsDetails([]sdk.IpDetails{{
			Ip:          "1.2.3.4",
			NetworkType: sdk.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := adaptIps([]sdk.Ip{{NetworkType: "tralala"}})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptIps(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := adaptIps([]sdk.Ip{{
			Ip:          "1.2.3.4",
			NetworkType: sdk.NETWORKTYPE_PUBLIC,
		}})

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})

	t.Run("error returned for invalid ip", func(t *testing.T) {
		_, err := adaptIps([]sdk.Ip{{NetworkType: "tralala"}})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptContract(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		endsAt := time.Now()
		renewalsAt := time.Now()
		createdAt := time.Now()

		sdkContract := sdk.NewContract(
			0,
			1,
			sdk.CONTRACTTYPE_MONTHLY,
			*sdk.NewNullableTime(&endsAt),
			renewalsAt,
			createdAt,
			sdk.CONTRACTSTATE_ACTIVE,
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
		sdkContract := sdk.Contract{BillingFrequency: 45}

		_, err := adaptContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "45")
	})

	t.Run("error returned for invalid term", func(t *testing.T) {
		sdkContract := sdk.Contract{BillingFrequency: 0, Term: 55}

		_, err := adaptContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("error returned for invalid type", func(t *testing.T) {
		sdkContract := sdk.Contract{
			BillingFrequency: 0,
			Term:             0,
			Type:             "tralala",
		}

		_, err := adaptContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("error returned for invalid state", func(t *testing.T) {
		sdkContract := sdk.Contract{
			BillingFrequency: 0,
			Term:             0,
			Type:             sdk.CONTRACTTYPE_HOURLY,
			State:            "tralala",
		}

		_, err := adaptContract(sdkContract)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run(
		"error returned when contract cannot be created",
		func(t *testing.T) {
			sdkContract := sdk.Contract{
				BillingFrequency: 0,
				Term:             0,
				Type:             sdk.CONTRACTTYPE_MONTHLY,
				State:            sdk.CONTRACTSTATE_ACTIVE,
			}

			_, err := adaptContract(sdkContract)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "contract.term cannot be 0")
		},
	)

}

func Test_adaptIso(t *testing.T) {
	got := adaptIso(*sdk.NewIso("id", "name"))

	assert.Equal(t, "id", got.Id)
	assert.Equal(t, "name", got.Name)
}

func Test_adaptPrivateNetwork(t *testing.T) {
	got := adaptPrivateNetwork(*sdk.NewPrivateNetwork(
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
		loadBalancerId := "loadBalancerId"

		sdkAutoScalingGroup := sdk.NewAutoScalingGroupDetails(
			instanceId,
			"MANUAL",
			"SCALING",
			*sdk.NewNullableInt32(&desiredAmount),
			"region",
			"reference",
			createdAt,
			updatedAt,
			*sdk.NewNullableTime(&startsAt),
			*sdk.NewNullableTime(&endsAt),
			*sdk.NewNullableInt32(&minimumAmount),
			*sdk.NewNullableInt32(&maximumAmount),
			*sdk.NewNullableInt32(&cpuThreshold),
			*sdk.NewNullableInt32(&warmupTime),
			*sdk.NewNullableInt32(&cooldownTime),
			*sdk.NewNullableLoadBalancer(&sdk.LoadBalancer{
				Id:    loadBalancerId,
				State: sdk.STATE_CREATING,
				Type:  sdk.TYPENAME_M3_LARGE,
			}),
		)

		got, err := AdaptAutoScalingGroupDetails(*sdkAutoScalingGroup)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, enum.AutoScalingCpuTypeManual, got.Type)
		assert.Equal(t, enum.AutoScalingGroupStateScaling, got.State)
		assert.Equal(t, 6, *got.DesiredAmount)
		assert.Equal(t, domain.Region{Name: "region"}, got.Region)
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

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := AdaptAutoScalingGroupDetails(
			sdk.AutoScalingGroupDetails{Id: instanceId, Type: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := AdaptAutoScalingGroupDetails(
			sdk.AutoScalingGroupDetails{
				Id:    instanceId,
				Type:  sdk.AUTOSCALINGGROUPTYPE_CPU_BASED,
				State: "tralala",
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid reference returns error", func(t *testing.T) {
		_, err := AdaptAutoScalingGroupDetails(
			sdk.AutoScalingGroupDetails{
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
			sdk.AutoScalingGroupDetails{
				Id:        instanceId,
				Type:      "MANUAL",
				State:     "SCALING",
				Reference: "reference",
				LoadBalancer: *sdk.NewNullableLoadBalancer(
					&sdk.LoadBalancer{State: "tralala"},
				),
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptStickySession(t *testing.T) {
	got := adaptStickySession(sdk.StickySession{
		Enabled:     false,
		MaxLifeTime: 20,
	})

	assert.False(t, got.Enabled)
	assert.Equal(t, 20, got.MaxLifeTime)

}

func Test_adaptHealthCheck(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		host := "host"

		sdkHealthCheck := sdk.NewHealthCheck(
			"GET",
			"uri",
			*sdk.NewNullableString(&host),
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
		_, err := adaptHealthCheck(sdk.HealthCheck{Method: "tralala"})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

}

func Test_adaptLoadBalancerConfiguration(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkLoadBalancerConfiguration := sdk.NewLoadBalancerConfiguration(
			*sdk.NewNullableStickySession(&sdk.StickySession{MaxLifeTime: 44}),
			"roundrobin",
			*sdk.NewNullableHealthCheck(&sdk.HealthCheck{Method: "GET"}),
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
			sdk.LoadBalancerConfiguration{Balance: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid HealthCheck returns error", func(t *testing.T) {
		_, err := adaptLoadBalancerConfiguration(
			sdk.LoadBalancerConfiguration{
				Balance: sdk.BALANCE_ROUNDROBIN,
				HealthCheck: *sdk.NewNullableHealthCheck(
					&sdk.HealthCheck{Method: "tralala"},
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
		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, string(sdk.TYPENAME_M3_LARGE), got.Type.String())
		assert.Equal(t, "unit", got.Resources.Cpu.Unit)
		assert.Equal(t, domain.Region{Name: "region"}, got.Region)
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
		sdkLoadBalancer.Ips = []sdk.IpDetails{{NetworkType: "tralala"}}

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
			sdk.InstanceType{
				Name:      "name",
				Resources: sdk.Resources{Cpu: sdk.Cpu{Unit: "cpu"}},
				Prices:    sdk.Prices{Currency: "currency"},
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
			sdk.InstanceType{
				StorageTypes: []sdk.RootDiskStorageType{
					sdk.ROOTDISKSTORAGETYPE_CENTRAL,
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
			sdk.InstanceType{
				StorageTypes: []sdk.RootDiskStorageType{"tralala"},
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestAdaptRegion(t *testing.T) {
	got := AdaptRegion(sdk.Region{Name: "name", Location: "location"})
	want := domain.Region{Name: "name", Location: "location"}

	assert.Equal(t, want, got)
}

func TestAdaptInstance(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		startedAt := time.Now()

		sdkInstance := generateInstance(t, &startedAt)

		got, err := AdaptInstance(sdkInstance)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, "lsw.m3.large", got.Type.String())
		assert.Equal(t, "cpu", got.Resources.Cpu.Unit)
		assert.Equal(t, domain.Region{Name: "region"}, got.Region)
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
		assert.Equal(t, "CENTOS_7_64BIT", got.Image.Id)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
		assert.Equal(t, autoScalingGroupId, got.AutoScalingGroup.Id)
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
		sdkInstance.Ips = []sdk.Ip{{NetworkType: "tralala"}}

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
		sdkInstance.AutoScalingGroup.Get().Type = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func Test_adaptImage(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkImage := sdk.NewImage(
			"UBUNTU_24_04_64BIT",
			"name",
			"family",
			"flavour",
			false,
		)

		got := adaptImage(*sdkImage)
		want := domain.Image{
			Id:           "UBUNTU_24_04_64BIT",
			Name:         "name",
			Family:       "family",
			Flavour:      "flavour",
			Custom:       false,
			MarketApps:   []string{},
			StorageTypes: []string{},
		}

		assert.Equal(t, want, got)
	})
}

func Test_adaptIp(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		reverseLookup := "reverseLookup"

		sdkIp := sdk.NewIp(
			"1.2.3.4",
			"prefixLength",
			5,
			true,
			false,
			sdk.NETWORKTYPE_INTERNAL,
			*sdk.NewNullableString(&reverseLookup),
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
		_, err := adaptIpDetails(sdk.IpDetails{NetworkType: "tralala"})

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

		sdkAutoScalingGroup := sdk.NewAutoScalingGroup(
			instanceId,
			"MANUAL",
			"SCALING",
			*sdk.NewNullableInt32(&desiredAmount),
			"region",
			"reference",
			createdAt,
			updatedAt,
			*sdk.NewNullableTime(&startsAt),
			*sdk.NewNullableTime(&endsAt),
			*sdk.NewNullableInt32(&minimumAmount),
			*sdk.NewNullableInt32(&maximumAmount),
			*sdk.NewNullableInt32(&cpuThreshold),
			*sdk.NewNullableInt32(&warmupTime),
			*sdk.NewNullableInt32(&cooldownTime),
		)

		got, err := adaptAutoScalingGroup(*sdkAutoScalingGroup)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, enum.AutoScalingCpuTypeManual, got.Type)
		assert.Equal(t, enum.AutoScalingGroupStateScaling, got.State)
		assert.Equal(t, 6, *got.DesiredAmount)
		assert.Equal(
			t,
			domain.Region{
				Name: "region",
			},
			got.Region,
		)
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

	t.Run("invalid type returns error", func(t *testing.T) {
		_, err := adaptAutoScalingGroup(
			sdk.AutoScalingGroup{Id: instanceId, Type: "tralala"},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		_, err := adaptAutoScalingGroup(
			sdk.AutoScalingGroup{
				Id:    instanceId,
				Type:  sdk.AUTOSCALINGGROUPTYPE_CPU_BASED,
				State: "tralala",
			},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid reference returns error", func(t *testing.T) {
		_, err := adaptAutoScalingGroup(
			sdk.AutoScalingGroup{
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
		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, string(sdk.TYPENAME_M3_LARGE), got.Type.String())
		assert.Equal(t, "unit", got.Resources.Cpu.Unit)
		assert.Equal(t, enum.StateCreating, got.State)
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
		sdkStorageTypes := []sdk.RootDiskStorageType{
			sdk.ROOTDISKSTORAGETYPE_CENTRAL,
			sdk.ROOTDISKSTORAGETYPE_LOCAL,
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
			sdkStorageTypes := []sdk.RootDiskStorageType{"tralala"}
			got, err := adaptStorageTypes(sdkStorageTypes)

			assert.Nil(t, got)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)
}

func Test_adaptPrice(t *testing.T) {
	sdkPrice := sdk.NewPrice("1", "2")
	got := adaptPrice(*sdkPrice)

	want := domain.Price{
		HourlyPrice:  "1",
		MonthlyPrice: "2",
	}

	assert.Equal(t, want, got)
}

func Test_adaptStorage(t *testing.T) {
	sdkStorage := sdk.NewStorage(
		sdk.Price{HourlyPrice: "1"},
		sdk.Price{HourlyPrice: "2"},
	)
	got := adaptStorage(*sdkStorage)

	want := domain.Storage{
		Local:   domain.Price{HourlyPrice: "1"},
		Central: domain.Price{HourlyPrice: "2"},
	}

	assert.Equal(t, want, got)
}

func Test_adaptPrices(t *testing.T) {
	sdkPrices := sdk.NewPrices(
		"currency",
		"symbol",
		sdk.Price{HourlyPrice: "1"},
		sdk.Storage{Central: sdk.Price{HourlyPrice: "2"}},
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
) sdk.InstanceDetails {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	return *sdk.NewInstanceDetails(
		instanceId,
		sdk.TYPENAME_M3_LARGE,
		sdk.Resources{Cpu: sdk.Cpu{Unit: "cpu"}},
		"region",
		*sdk.NewNullableString(&reference),
		*sdk.NewNullableTime(startedAt),
		*sdk.NewNullableString(&marketAppId),
		sdk.STATE_RUNNING,
		"productType",
		true,
		false,
		6,
		sdk.ROOTDISKSTORAGETYPE_CENTRAL,
		sdk.Contract{
			BillingFrequency: 1,
			Type:             sdk.CONTRACTTYPE_HOURLY,
			State:            sdk.CONTRACTSTATE_ACTIVE,
		},
		*sdk.NewNullableAutoScalingGroup(&sdk.AutoScalingGroup{
			Id:    autoScalingGroupId,
			Type:  sdk.AUTOSCALINGGROUPTYPE_CPU_BASED,
			State: sdk.AUTOSCALINGGROUPSTATE_ACTIVE,
		}),
		sdk.Image{Id: "CENTOS_7_64BIT"},
		*sdk.NewNullableIso(&sdk.Iso{Id: "isoId"}),
		*sdk.NewNullablePrivateNetwork(
			&sdk.PrivateNetwork{PrivateNetworkId: "privateNetworkId"},
		),
		[]sdk.IpDetails{
			{Ip: "1.2.3.4", NetworkType: sdk.NETWORKTYPE_PUBLIC},
		},
		*sdk.NewNullableVolume(&sdk.Volume{Size: 3, Unit: "unit"}),
	)
}

func generateInstance(
	t *testing.T,
	startedAt *time.Time,
) sdk.Instance {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	return *sdk.NewInstance(
		instanceId,
		sdk.TYPENAME_M3_LARGE,
		sdk.Resources{Cpu: sdk.Cpu{Unit: "cpu"}},
		"region",
		*sdk.NewNullableString(&reference),
		*sdk.NewNullableTime(startedAt),
		*sdk.NewNullableString(&marketAppId),
		sdk.STATE_RUNNING,
		"productType",
		true,
		false,
		6,
		sdk.ROOTDISKSTORAGETYPE_CENTRAL,
		sdk.Contract{
			BillingFrequency: 1,
			Type:             sdk.CONTRACTTYPE_HOURLY,
			State:            sdk.CONTRACTSTATE_ACTIVE,
		},
		*sdk.NewNullableAutoScalingGroup(&sdk.AutoScalingGroup{
			Id:    autoScalingGroupId,
			Type:  sdk.AUTOSCALINGGROUPTYPE_CPU_BASED,
			State: sdk.AUTOSCALINGGROUPSTATE_ACTIVE,
		}),
		sdk.Image{Id: "CENTOS_7_64BIT"},
		[]sdk.Ip{
			{Ip: "1.2.3.4", NetworkType: sdk.NETWORKTYPE_PUBLIC},
		},
	)
}

func generateLoadBalancerDetails(startedAt *time.Time) sdk.LoadBalancerDetails {
	reference := "reference"

	return *sdk.NewLoadBalancerDetails(
		instanceId,
		"lsw.m3.large",
		sdk.Resources{Cpu: sdk.Cpu{Unit: "unit"}},
		*sdk.NewNullableString(&reference),
		"CREATING",
		*sdk.NewNullableTime(startedAt),
		[]sdk.IpDetails{{
			Ip:          "1.2.3.4",
			NetworkType: sdk.NETWORKTYPE_PUBLIC,
		}},
		"region",
		*sdk.NewNullableLoadBalancerConfiguration(&sdk.LoadBalancerConfiguration{
			TargetPort: 22,
			Balance:    "roundrobin",
		}),
		*sdk.NewNullableAutoScalingGroup(nil),
		*sdk.NewNullablePrivateNetwork(
			&sdk.PrivateNetwork{PrivateNetworkId: "privateNetworkId"},
		),
		sdk.Contract{
			BillingFrequency: 1,
			Type:             sdk.CONTRACTTYPE_MONTHLY,
			State:            sdk.CONTRACTSTATE_ACTIVE,
			Term:             sdk.CONTRACTTERM__1,
		},
	)
}

func generateLoadBalancer(startedAt *time.Time) sdk.LoadBalancer {
	reference := "reference"

	return *sdk.NewLoadBalancer(
		instanceId,
		"lsw.m3.large",
		sdk.Resources{Cpu: sdk.Cpu{Unit: "unit"}},
		*sdk.NewNullableString(&reference),
		sdk.STATE_CREATING,
		*sdk.NewNullableTime(startedAt),
	)
}

func Test_adaptVolume(t *testing.T) {
	sdkVolume := sdk.NewVolume(1, "unit")
	got := adaptVolume(*sdkVolume)
	want := domain.Volume{Size: 1, Unit: "unit"}

	assert.Equal(t, want, got)
}

func TestAdaptImageDetails(t *testing.T) {
	state := "state"
	stateReason := "stateReason"
	createdAt := time.Now()
	updatedAt := time.Now()
	version := "version"
	architecture := "architecture"
	region := "region"

	sdkImageDetails := sdk.NewImageDetails(
		"id",
		"name",
		"family",
		"flavour",
		false,
		*sdk.NewNullableStorageSize(
			sdk.NewStorageSize(float32(1), "unit"),
		),
		*sdk.NewNullableString(&state),
		*sdk.NewNullableString(&stateReason),
		*sdk.NewNullableString(&region),
		*sdk.NewNullableTime(&createdAt),
		*sdk.NewNullableTime(&updatedAt),
		"version",
		"architecture",
		[]string{"marketApp"},
		[]string{"storageType"},
	)

	got := AdaptImageDetails(*sdkImageDetails)

	want := domain.Image{
		Id:           "id",
		Name:         "name",
		Version:      &version,
		Family:       "family",
		Flavour:      "flavour",
		Architecture: &architecture,
		State:        &state,
		StateReason:  &stateReason,
		Region: &domain.Region{
			Name: "region",
		},
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
		Custom:       false,
		StorageSize:  &domain.StorageSize{Size: 1, Unit: "unit"},
		MarketApps:   []string{"marketApp"},
		StorageTypes: []string{"storageType"},
	}

	assert.Equal(t, want, got)
}

func Test_adaptStorageSize(t *testing.T) {
	sdkStorageSize := sdk.NewStorageSize(1, "unit")

	got := adaptStorageSize(*sdkStorageSize)
	want := domain.StorageSize{Size: 1, Unit: "unit"}

	assert.Equal(t, want, got)
}
