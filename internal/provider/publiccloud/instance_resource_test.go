package publiccloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptContractToContractResource(t *testing.T) {
	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	sdkContract := publicCloud.Contract{
		BillingFrequency: publicCloud.BILLINGFREQUENCY__1,
		Term:             publicCloud.CONTRACTTERM__3,
		Type:             publicCloud.CONTRACTTYPE_HOURLY,
		EndsAt:           *publicCloud.NewNullableTime(&endsAt),
		State:            publicCloud.CONTRACTSTATE_ACTIVE,
	}

	want := contractResourceModel{
		BillingFrequency: basetypes.NewInt64Value(1),
		Term:             basetypes.NewInt64Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := adaptContractToContractResource(sdkContract)

	assert.Equal(t, want, got)
}

func Test_contractResourceModel_IsContractTermValid(t *testing.T) {
	t.Run(
		"false is returned when contract term is monthly and contract term is 0",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__0,
				Type: publicCloud.CONTRACTTYPE_MONTHLY,
			}

			contract := adaptContractToContractResource(sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.False(t, got)
			assert.Equal(t, reasonContractTermCannotBeZero, reason)
		},
	)

	t.Run(
		"false is returned when contract term is hourly and contract term is not 0",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__3,
				Type: publicCloud.CONTRACTTYPE_HOURLY,
			}

			contract := adaptContractToContractResource(sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.False(t, got)
			assert.Equal(t, reasonContractTermMustBeZero, reason)
		},
	)

	t.Run(
		"true is returned when contract term is valid",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__0,
				Type: publicCloud.CONTRACTTYPE_HOURLY,
			}

			contract := adaptContractToContractResource(sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.True(t, got)
			assert.Equal(t, reasonNone, reason)
		},
	)
}

func generateContractObject(
	billingFrequency *int,
	contractTerm *int,
	contractType *string,
) types.Object {
	defaultContractType := "MONTHLY"
	defaultContractTerm := 3
	defaultBillingFrequency := 1

	if contractType == nil {
		contractType = &defaultContractType
	}
	if contractTerm == nil {
		contractTerm = &defaultContractTerm
	}
	if billingFrequency == nil {
		billingFrequency = &defaultBillingFrequency
	}

	contract, _ := types.ObjectValueFrom(
		context.TODO(),
		contractResourceModel{}.AttributeTypes(),
		contractResourceModel{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			State:            basetypes.NewStringUnknown(),
		},
	)

	return contract
}

func generateInstanceResourceModel() instanceResourceModel {
	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

	image, _ := types.ObjectValueFrom(
		context.TODO(),
		imageResourceModel{}.AttributeTypes(),
		imageResourceModel{
			ID:           basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
			MarketApps:   emptyList,
			StorageTypes: emptyList,
		},
	)

	contract := generateContractObject(nil, nil, nil)
	instance := instanceResourceModel{
		ID:                  basetypes.NewStringValue("id"),
		Region:              basetypes.NewStringValue("eu-west-3"),
		Type:                basetypes.NewStringValue("lsw.m5a.4xlarge"),
		RootDiskStorageType: basetypes.NewStringValue("CENTRAL"),
		RootDiskSize:        basetypes.NewInt64Value(int64(55)),
		Image:               image,
		Contract:            contract,
		MarketAppID:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
	}

	return instance
}

func Test_adaptInstanceToInstanceResource(t *testing.T) {
	marketAppId := "marketAppId"
	reference := "reference"

	instance := publicCloud.Instance{
		Id:                  "id",
		Type:                publicCloud.TYPENAME_C3_2XLARGE,
		Region:              "region",
		Reference:           *publicCloud.NewNullableString(&reference),
		MarketAppId:         *publicCloud.NewNullableString(&marketAppId),
		State:               publicCloud.STATE_CREATING,
		RootDiskSize:        50,
		RootDiskStorageType: publicCloud.STORAGETYPE_CENTRAL,
		Contract: publicCloud.Contract{
			Type: publicCloud.CONTRACTTYPE_MONTHLY,
		},
		Image: publicCloud.Image{
			Id: "UBUNTU_20_04_64BIT",
		},
		Ips: []publicCloud.Ip{
			{
				Ip: "127.0.0.1",
			},
		},
	}

	got, err := adaptInstanceToInstanceResource(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := imageResourceModel{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.ID.ValueString())

	contract := contractResourceModel{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []iPResourceModel
	got.IPs.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].IP.ValueString())
}

func Test_adaptInstanceDetailsToInstanceResource(t *testing.T) {
	marketAppId := "marketAppId"
	reference := "reference"

	instance := publicCloud.InstanceDetails{
		Id:                  "id",
		Type:                publicCloud.TYPENAME_C3_2XLARGE,
		Region:              "region",
		Reference:           *publicCloud.NewNullableString(&reference),
		MarketAppId:         *publicCloud.NewNullableString(&marketAppId),
		State:               publicCloud.STATE_CREATING,
		RootDiskSize:        50,
		RootDiskStorageType: publicCloud.STORAGETYPE_CENTRAL,
		Contract: publicCloud.Contract{
			Type: publicCloud.CONTRACTTYPE_MONTHLY,
		},
		Image: publicCloud.Image{
			Id: "UBUNTU_20_04_64BIT",
		},
		Ips: []publicCloud.IpDetails{
			{
				Ip: "127.0.0.1",
			},
		},
	}

	got, err := adaptInstanceDetailsToInstanceResource(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := imageResourceModel{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.ID.ValueString())

	contract := contractResourceModel{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []iPResourceModel
	got.IPs.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].IP.ValueString())
}

func Test_instanceResourceModel_GetLaunchInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceResourceModel()
		instance.MarketAppID = basetypes.NewStringPointerValue(nil)
		instance.Reference = basetypes.NewStringPointerValue(nil)
		instance.RootDiskSize = basetypes.NewInt64PointerValue(nil)

		got, err := instance.GetLaunchInstanceOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.REGIONNAME_EU_WEST_3, got.Region)
		assert.Equal(t, publicCloud.TYPENAME_M5A_4XLARGE, got.Type)
		assert.Equal(t, publicCloud.STORAGETYPE_CENTRAL, got.RootDiskStorageType)
		assert.Equal(t, "UBUNTU_20_04_64BIT", got.ImageId)
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.ContractType)
		assert.Equal(t, publicCloud.CONTRACTTERM__3, got.ContractTerm)
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__1, got.BillingFrequency)

		marketAppId, _ := got.GetMarketAppIdOk()
		assert.Nil(t, marketAppId)

		reference, _ := got.GetReferenceOk()
		assert.Nil(t, reference)

		rootDiskSize, _ := got.GetRootDiskSizeOk()
		assert.Nil(t, rootDiskSize)
	})

	t.Run("optional values are passed", func(t *testing.T) {
		instance := generateInstanceResourceModel()

		got, err := instance.GetLaunchInstanceOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, int32(55), *got.RootDiskSize)
	})

	t.Run(
		"returns error if invalid rootDiskStorageType is passed",
		func(t *testing.T) {
			instance := generateInstanceResourceModel()
			instance.RootDiskStorageType = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instance := generateInstanceResourceModel()
			instance.Type = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractType is passed",
		func(t *testing.T) {
			contractType := "tralala"
			instance := generateInstanceResourceModel()
			contract := generateContractObject(nil, nil, &contractType)
			instance.Contract = contract

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			instance := generateInstanceResourceModel()
			contract := generateContractObject(nil, &contractTerm, nil)
			instance.Contract = contract

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			instance := generateInstanceResourceModel()
			contract := generateContractObject(&billingFrequency, nil, nil)
			instance.Contract = contract

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run("returns error if invalid region is passed", func(t *testing.T) {
		instance := generateInstanceResourceModel()
		instance.Region = basetypes.NewStringValue("tralala")

		_, err := instance.GetLaunchInstanceOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run(
		"returns error if image resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceResourceModel()
			instance.Image = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".imageResourceModel")
		},
	)

	t.Run(
		"returns error if contractResourceModel resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceResourceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".contractResourceModel")
		},
	)
}

func Test_instanceResourceModel_GetUpdateInstanceOpts(t *testing.T) {
	t.Run("optional values are set", func(t *testing.T) {
		instance := generateInstanceResourceModel()

		got, err := instance.GetUpdateInstanceOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.TYPENAME_M5A_4XLARGE, *got.Type)
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, *got.ContractType)
		assert.Equal(t, publicCloud.CONTRACTTERM__3, *got.ContractTerm)
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__1, *got.BillingFrequency)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, int32(55), *got.RootDiskSize)
	})

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instance := generateInstanceResourceModel()
			instance.Type = basetypes.NewStringValue("tralala")

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractType is passed",
		func(t *testing.T) {
			contractType := "tralala"
			instance := generateInstanceResourceModel()
			contract := generateContractObject(nil, nil, &contractType)
			instance.Contract = contract

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			instance := generateInstanceResourceModel()
			contract := generateContractObject(nil, &contractTerm, nil)
			instance.Contract = contract

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			instance := generateInstanceResourceModel()
			contract := generateContractObject(&billingFrequency, nil, nil)
			instance.Contract = contract

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if contractResourceModel resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceResourceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".contractResourceModel")
		},
	)
}

func Test_adaptIpToIPResource(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := iPResourceModel{
		IP: basetypes.NewStringValue("127.0.0.1"),
	}
	got := adaptIpToIPResource(sdkIp)

	assert.Equal(t, want, got)
}

func Test_adaptIpDetailsToIPResource(t *testing.T) {
	sdkIpDetails := publicCloud.IpDetails{
		Ip: "127.0.0.1",
	}

	want := iPResourceModel{
		IP: basetypes.NewStringValue("127.0.0.1"),
	}
	got := adaptIpDetailsToIPResource(sdkIpDetails)

	assert.Equal(t, want, got)
}
