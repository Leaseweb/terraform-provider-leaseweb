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

func TestAdaptInstanceDetails(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t)

		got, err := AdaptInstanceDetails(sdkInstance)

		assert.NoError(t, err)
		assert.Equal(
			t,
			"5d7f8262-d77f-4476-8da8-6a84f8f2ae8d",
			got.Id,
		)
		assert.Equal(t, string(sdk.TYPENAME_M3_LARGE), got.Type)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, enum.StateRunning, got.State)
		assert.Equal(t, 6, got.RootDiskSize.Value)
		assert.Equal(t, enum.StorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "CENTOS_7_64BIT", got.Image.Id)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t)
		sdkInstance.State = "tralala"

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t)
		sdkInstance.RootDiskSize = 5000

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t)
		sdkInstance.RootDiskStorageType = "tralala"

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkInstance := generateInstanceDetails(t)
		sdkInstance.Contract.BillingFrequency = 55

		_, err := AdaptInstanceDetails(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})
}

func Test_adaptIpDetails(t *testing.T) {
	sdkIp := sdk.NewIpDetails(
		"1.2.3.4",
		"",
		5,
		true,
		false,
		sdk.NETWORKTYPE_INTERNAL,
		*sdk.NewNullableString(nil),
		*sdk.NewNullableDdos(nil),
	)

	got := adaptIpDetails(*sdkIp)

	assert.Equal(t, "1.2.3.4", got.Ip)
}

func Test_adaptIpsDetails(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got := adaptIpsDetails([]sdk.IpDetails{{
			Ip: "1.2.3.4",
		}})

		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})
}

func Test_adaptIps(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got := adaptIps([]sdk.Ip{{
			Ip: "1.2.3.4",
		}})

		assert.Len(t, got, 1)
		assert.Equal(t, "1.2.3.4", got[0].Ip)
	})
}

func Test_adaptContract(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		endsAt := time.Now()

		sdkContract := sdk.NewContract(
			0,
			1,
			sdk.CONTRACTTYPE_MONTHLY,
			*sdk.NewNullableTime(&endsAt),
			time.Now(),
			time.Now(),
			sdk.CONTRACTSTATE_ACTIVE,
		)

		got, err := adaptContract(*sdkContract)

		assert.NoError(t, err)
		assert.Equal(t, enum.ContractBillingFrequencyZero, got.BillingFrequency)
		assert.Equal(t, enum.ContractTermOne, got.Term)
		assert.Equal(t, enum.ContractTypeMonthly, got.Type)
		assert.Equal(t, endsAt, *got.EndsAt)
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

func TestAdaptInstance(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkInstance := generateInstance(t)

		got, err := AdaptInstance(sdkInstance)

		assert.NoError(t, err)
		assert.Equal(t, instanceId, got.Id)
		assert.Equal(t, "lsw.m3.large", got.Type)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, enum.StateRunning, got.State)
		assert.Equal(t, 6, got.RootDiskSize.Value)
		assert.Equal(t, enum.StorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "CENTOS_7_64BIT", got.Image.Id)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)
	})

	t.Run("invalid state returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t)
		sdkInstance.State = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
	})

	t.Run("invalid rootDiskSize returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t)
		sdkInstance.RootDiskSize = 5000

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "5000")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t)
		sdkInstance.RootDiskStorageType = "tralala"

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contract returns error", func(t *testing.T) {
		sdkInstance := generateInstance(t)
		sdkInstance.Contract.BillingFrequency = 55

		_, err := AdaptInstance(sdkInstance)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})
}

func Test_adaptImage(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		sdkImage := sdk.NewImage(
			"UBUNTU_24_04_64BIT",
			"",
			"",
			"",
			false,
		)

		got := adaptImage(*sdkImage)
		want := domain.Image{
			Id: "UBUNTU_24_04_64BIT",
		}

		assert.Equal(t, want, got)
	})
}

func Test_adaptIp(t *testing.T) {
	sdkIp := sdk.NewIp(
		"1.2.3.4",
		"",
		5,
		true,
		false,
		sdk.NETWORKTYPE_INTERNAL,
		*sdk.NewNullableString(nil),
	)

	got := adaptIp(*sdkIp)

	assert.Equal(t, "1.2.3.4", got.Ip)
}

func generateInstanceDetails(t *testing.T) sdk.InstanceDetails {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	return *sdk.NewInstanceDetails(
		instanceId,
		sdk.TYPENAME_M3_LARGE,
		sdk.Resources{Cpu: sdk.Cpu{Unit: "cpu"}},
		"region",
		*sdk.NewNullableString(&reference),
		*sdk.NewNullableTime(nil),
		*sdk.NewNullableString(&marketAppId),
		sdk.STATE_RUNNING,
		"productType",
		true,
		false,
		false,
		6,
		sdk.STORAGETYPE_CENTRAL,
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
	)
}

func generateInstance(t *testing.T) sdk.Instance {
	t.Helper()

	reference := "reference"
	marketAppId := "marketAppId"

	return *sdk.NewInstance(
		instanceId,
		sdk.TYPENAME_M3_LARGE,
		sdk.Resources{Cpu: sdk.Cpu{Unit: "cpu"}},
		"region",
		*sdk.NewNullableString(&reference),
		*sdk.NewNullableTime(nil),
		*sdk.NewNullableString(&marketAppId),
		sdk.STATE_RUNNING,
		"productType",
		true,
		false,
		false,
		6,
		sdk.STORAGETYPE_CENTRAL,
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
