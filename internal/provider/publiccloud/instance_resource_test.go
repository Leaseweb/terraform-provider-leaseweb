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

func Test_resourceModelContract_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		resourceModelContract{}.AttributeTypes(),
		resourceModelContract{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}

func Test_adaptSdkContractToResourceContract(t *testing.T) {
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

	want := resourceModelContract{
		BillingFrequency: basetypes.NewInt64Value(1),
		Term:             basetypes.NewInt64Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := adaptSdkContractToResourceContract(sdkContract)

	assert.Equal(t, want, got)
}

func Test_contract_IsContractTermValid(t *testing.T) {
	t.Run(
		"false is returned when contract term is monthly and contract term is 0",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__0,
				Type: publicCloud.CONTRACTTYPE_MONTHLY,
			}

			contract := adaptSdkContractToResourceContract(sdkContract)

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

			contract := adaptSdkContractToResourceContract(sdkContract)

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

			contract := adaptSdkContractToResourceContract(sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.True(t, got)
			assert.Equal(t, reasonNone, reason)
		},
	)
}

func Test_adaptSdkImageToResourceImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id: "imageId",
	}

	want := resourceModelImage{
		ID: basetypes.NewStringValue("imageId"),
	}
	got := adaptSdkImageToResourceImage(sdkImage)

	assert.Equal(t, want, got)
}

func GenerateContractObject(
	billingFrequency *int,
	contractTerm *int,
	contractType *string,
	endsAt *string,
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
		resourceModelContract{}.AttributeTypes(),
		resourceModelContract{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			State:            basetypes.NewStringUnknown(),
			EndsAt:           basetypes.NewStringPointerValue(endsAt),
		},
	)

	return contract
}

func generateInstanceModel() resourceModelInstance {
	image, _ := types.ObjectValueFrom(
		context.TODO(),
		resourceModelImage{}.AttributeTypes(),
		resourceModelImage{
			ID: basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
		},
	)

	contract := GenerateContractObject(
		nil,
		nil,
		nil,
		nil,
	)

	instance := resourceModelInstance{
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

func Test_adaptSdkInstanceToResourceInstance(t *testing.T) {
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

	got, err := adaptSdkInstanceToResourceInstance(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := resourceModelImage{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.ID.ValueString())

	contract := resourceModelContract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []resourceModelIP
	got.IPs.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].IP.ValueString())
}

func Test_adaptSdkInstanceDetailsToResourceInstance(t *testing.T) {
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

	got, err := adaptSdkInstanceDetailsToResourceInstance(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := resourceModelImage{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.ID.ValueString())

	contract := resourceModelContract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []resourceModelIP
	got.IPs.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].IP.ValueString())
}

func TestInstance_GetLaunchInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceModel()
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
		instance := generateInstanceModel()

		got, err := instance.GetLaunchInstanceOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, int32(55), *got.RootDiskSize)
	})

	t.Run(
		"returns error if invalid rootDiskStorageType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.RootDiskStorageType = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
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
			instance := generateInstanceModel()
			contract := GenerateContractObject(
				nil,
				nil,
				&contractType,
				nil,
			)
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
			instance := generateInstanceModel()
			contract := GenerateContractObject(
				nil,
				&contractTerm,
				nil,
				nil,
			)
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
			instance := generateInstanceModel()
			contract := GenerateContractObject(
				&billingFrequency,
				nil,
				nil,
				nil,
			)
			instance.Contract = contract

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid region is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Region = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if resourceModelImage resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Image = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".resourceModelImage")
		},
	)

	t.Run(
		"returns error if resourceModelContract resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".resourceModelContract")
		},
	)
}

func TestInstance_GetUpdateInstanceOpts(t *testing.T) {
	t.Run("optional values are set", func(t *testing.T) {
		instance := generateInstanceModel()

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
			instance := generateInstanceModel()
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
			instance := generateInstanceModel()
			contract := GenerateContractObject(
				nil,
				nil,
				&contractType,
				nil,
			)
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
			instance := generateInstanceModel()
			contract := GenerateContractObject(
				nil,
				&contractTerm,
				nil,
				nil,
			)
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
			instance := generateInstanceModel()
			contract := GenerateContractObject(
				&billingFrequency,
				nil,
				nil,
				nil,
			)
			instance.Contract = contract

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if resourceModelContract resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetUpdateInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".resourceModelContract")
		},
	)
}

func Test_instance_CanBeTerminated(t *testing.T) {
	t.Run("instance can be terminated", func(t *testing.T) {
		instance := generateInstanceModel()
		instance.State = basetypes.NewStringValue(string(publicCloud.STATE_UNKNOWN))

		got := instance.CanBeTerminated(context.TODO())

		assert.Nil(t, got)
	})

	t.Run(
		"instance cannot be terminated if state is CREATING/DESTROYING/DESTROYED",
		func(t *testing.T) {
			tests := []struct {
				name           string
				state          publicCloud.State
				reasonContains string
			}{
				{
					name:           "state is CREATING",
					state:          publicCloud.STATE_CREATING,
					reasonContains: "CREATING",
				},
				{
					name:           "state is DESTROYING",
					state:          publicCloud.STATE_DESTROYING,
					reasonContains: "DESTROYING",
				},
				{
					name:           "state is DESTROYED",
					state:          publicCloud.STATE_DESTROYED,
					reasonContains: "DESTROYED",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					instance := generateInstanceModel()
					instance.State = basetypes.NewStringValue(string(tt.state))

					got := instance.CanBeTerminated(context.TODO())

					assert.NotNil(t, got)
					assert.Contains(t, *got, tt.reasonContains)
				})
			}
		},
	)

	t.Run(
		"instance cannot be terminated if contract.endsAt is set",
		func(t *testing.T) {
			endsAt := "2023-12-14 17:09:47 +0000 UTC"

			contract := GenerateContractObject(nil, nil, nil, &endsAt)

			instance := generateInstanceModel()
			instance.State = basetypes.NewStringValue(string(publicCloud.STATE_UNKNOWN))
			instance.Contract = contract

			got := instance.CanBeTerminated(context.TODO())

			assert.NotNil(t, got)
			assert.Contains(t, *got, "2023-12-14 17:09:47 +0000 UTC")
		},
	)
}

func Test_adaptSdkIpToResourceIp(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := resourceModelIP{
		IP: basetypes.NewStringValue("127.0.0.1"),
	}
	got := adaptSdkIpToResourceIP(sdkIp)

	assert.Equal(t, want, got)
}

func Test_adaptSdkIpDetailsToResourceIp(t *testing.T) {
	sdkIpDetails := publicCloud.IpDetails{
		Ip: "127.0.0.1",
	}

	want := resourceModelIP{
		IP: basetypes.NewStringValue("127.0.0.1"),
	}
	got := adaptSdkIpDetailsToResourceIP(sdkIpDetails)

	assert.Equal(t, want, got)
}
