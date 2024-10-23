package publiccloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_contractTermValidator_ValidateObject(t *testing.T) {
	t.Run(
		"does not set error if contract term is correct",
		func(t *testing.T) {
			contract := resourceModelContract{}
			configValue, _ := types.ObjectValueFrom(
				context.TODO(),
				contract.AttributeTypes(),
				contract,
			)

			request := terraformValidator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := terraformValidator.ObjectResponse{}

			validator := contractTermValidator{}
			validator.ValidateObject(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"returns expected error if contract term cannot be 0",
		func(t *testing.T) {
			contract := resourceModelContract{
				Type: basetypes.NewStringValue("MONTHLY"),
				Term: basetypes.NewInt64Value(0),
			}
			configValue, _ := types.ObjectValueFrom(
				context.TODO(),
				contract.AttributeTypes(),
				contract,
			)

			request := terraformValidator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := terraformValidator.ObjectResponse{}

			validator := contractTermValidator{}
			validator.ValidateObject(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 1)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"MONTHLY",
			)
		},
	)

	t.Run(
		"returns expected error if contract term must be 0",
		func(t *testing.T) {
			contract := resourceModelContract{
				Type: basetypes.NewStringValue("HOURLY"),
				Term: basetypes.NewInt64Value(3),
			}
			configValue, _ := types.ObjectValueFrom(
				context.TODO(),
				contract.AttributeTypes(),
				contract,
			)

			request := terraformValidator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := terraformValidator.ObjectResponse{}

			validator := contractTermValidator{}
			validator.ValidateObject(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 1)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"HOURLY",
			)
		},
	)
}

func Test_instanceTerminationValidator_ValidateObject(t *testing.T) {
	t.Run("ConfigValue populate errors bubble up", func(t *testing.T) {
		request := terraformValidator.ObjectRequest{}
		response := terraformValidator.ObjectResponse{}

		validator := instanceTerminationValidator{}
		validator.ValidateObject(context.TODO(), request, &response)

		assert.True(t, response.Diagnostics.HasError())
		assert.Contains(
			t,
			response.Diagnostics[0].Summary(),
			"Value Conversion Error",
		)
	})

	t.Run(
		"does not set a diagnostics error if instance is allowed to be terminated",
		func(t *testing.T) {
			instance := generateInstanceModelForValidator()
			instanceObject, _ := basetypes.NewObjectValueFrom(
				context.TODO(),
				instance.AttributeTypes(),
				instance,
			)
			request := terraformValidator.ObjectRequest{ConfigValue: instanceObject}
			response := terraformValidator.ObjectResponse{}

			validator := instanceTerminationValidator{}
			validator.ValidateObject(context.TODO(), request, &response)

			assert.False(t, response.Diagnostics.HasError())
		},
	)

	t.Run(
		"sets a diagnostics error if instance is not allowed to be terminated",
		func(t *testing.T) {
			instance := generateInstanceModelForValidator()
			instance.State = basetypes.NewStringValue("DESTROYED")
			instanceObject, _ := basetypes.NewObjectValueFrom(
				context.TODO(),
				instance.AttributeTypes(),
				instance,
			)
			request := terraformValidator.ObjectRequest{ConfigValue: instanceObject}
			response := terraformValidator.ObjectResponse{}

			validator := instanceTerminationValidator{}
			validator.ValidateObject(context.TODO(), request, &response)

			assert.True(t, response.Diagnostics.HasError())
			assert.Contains(t, response.Diagnostics[0].Detail(), "DESTROYED")
		},
	)
}

func generateInstanceModelForValidator() resourceModelInstance {
	contract := resourceModelContract{}
	contractObject, _ := types.ObjectValueFrom(
		context.TODO(),
		contract.AttributeTypes(),
		contract,
	)

	return resourceModelInstance{
		ID:        basetypes.NewStringUnknown(),
		Region:    basetypes.NewStringUnknown(),
		Reference: basetypes.NewStringUnknown(),
		Image: basetypes.NewObjectUnknown(
			resourceModelImage{}.AttributeTypes(),
		),
		State:               basetypes.NewStringUnknown(),
		Type:                basetypes.NewStringUnknown(),
		RootDiskSize:        basetypes.NewInt64Unknown(),
		RootDiskStorageType: basetypes.NewStringUnknown(),
		Ips: basetypes.NewListUnknown(
			types.ObjectType{
				AttrTypes: resourceModelIp{}.AttributeTypes(),
			},
		),
		Contract:    contractObject,
		MarketAppId: basetypes.NewStringUnknown(),
	}
}

func Test_instanceTypeValidator_ValidateString(t *testing.T) {
	t.Run("nothing happens if instanceType is unknown", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		validator := instanceTypeValidator{}

		response := terraformValidator.StringResponse{}
		validator.ValidateString(
			context.TODO(),
			terraformValidator.StringRequest{ConfigValue: basetypes.NewStringUnknown()},
			&response,
		)

		assert.Equal(t, 0, countIsInstanceTypeAvailableForRegionIsCalled)
		assert.Equal(t, 0, countCanInstanceTypeBeUsedWithInstanceIsCalled)
	})

	t.Run("nothing happens if instanceType does not change", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		validator := instanceTypeValidator{}

		response := terraformValidator.StringResponse{}
		validator.ValidateString(
			context.TODO(),
			terraformValidator.StringRequest{
				ConfigValue: basetypes.NewStringNull(),
			},
			&response,
		)

		assert.Equal(t, 0, countIsInstanceTypeAvailableForRegionIsCalled)
		assert.Equal(t, 0, countCanInstanceTypeBeUsedWithInstanceIsCalled)
	})

	t.Run(
		"attributeError added to response if instanceType cannot be found",
		func(t *testing.T) {
			validator := instanceTypeValidator{
				availableInstanceTypes: []string{"tralala"},
			}

			response := terraformValidator.StringResponse{}
			validator.ValidateString(
				context.TODO(),
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("doesNotExist"),
				},
				&response,
			)

			assert.Contains(
				t,
				response.Diagnostics[0].Detail(),
				"tralala",
			)
			assert.Contains(
				t,
				response.Diagnostics[0].Detail(),
				"doesNotExist",
			)
		},
	)

	t.Run(
		"attributeError not added to response if instanceType can be found",
		func(t *testing.T) {
			validator := instanceTypeValidator{
				availableInstanceTypes: []string{"tralala"},
			}

			response := terraformValidator.StringResponse{}
			validator.ValidateString(
				context.TODO(),
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("tralala"),
				},
				&response,
			)

			assert.Len(t, response.Diagnostics, 0)
		},
	)
}

func Test_newInstanceTypeValidator(t *testing.T) {
	validator := newInstanceTypeValidator(
		basetypes.NewStringValue("currentInstanceType"),
		[]string{"type1"},
	)

	assert.Equal(
		t,
		[]string{"type1", "currentInstanceType"},
		validator.availableInstanceTypes,
	)
}

func Test_regionValidator_ValidateString(t *testing.T) {
	t.Run("does not set errors if the region exists", func(t *testing.T) {
		request := terraformValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := terraformValidator.StringResponse{}

		validator := regionValidator{
			regions: []string{"region"},
		}
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run(
		"does not set errors if the region is unknown",
		func(t *testing.T) {
			request := terraformValidator.StringRequest{
				ConfigValue: basetypes.NewStringUnknown(),
			}

			response := terraformValidator.StringResponse{}

			validator := regionValidator{}
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set errors if the region is null",
		func(t *testing.T) {
			request := terraformValidator.StringRequest{
				ConfigValue: basetypes.NewStringNull(),
			}

			response := terraformValidator.StringResponse{}

			validator := regionValidator{}
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run("sets an error if the region does not exist", func(t *testing.T) {
		request := terraformValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := terraformValidator.StringResponse{}

		validator := regionValidator{
			regions: []string{"tralala"},
		}

		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
		assert.Contains(
			t,
			response.Diagnostics.Errors()[0].Detail(),
			"region",
		)
		assert.Contains(
			t,
			response.Diagnostics.Errors()[0].Detail(),
			"tralala",
		)
	})
}

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
	got, err := adaptSdkContractToResourceContract(context.TODO(), sdkContract)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_contract_IsContractTermValid(t *testing.T) {
	t.Run(
		"false is returned when contract term is monthly and contract term is 0",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__0,
				Type: publicCloud.CONTRACTTYPE_MONTHLY,
			}

			contract, _ := adaptSdkContractToResourceContract(context.TODO(), sdkContract)

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

			contract, _ := adaptSdkContractToResourceContract(context.TODO(), sdkContract)

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

			contract, _ := adaptSdkContractToResourceContract(context.TODO(), sdkContract)

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
		Id: basetypes.NewStringValue("imageId"),
	}
	got, err := adaptSdkImageToResourceImage(context.TODO(), sdkImage)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func generateContractObject(
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
			Id: basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
		},
	)

	contract := generateContractObject(
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
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
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
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := resourceModelImage{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.Id.ValueString())

	contract := resourceModelContract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []resourceModelIp
	got.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].Ip.ValueString())
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
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := resourceModelImage{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.Id.ValueString())

	contract := resourceModelContract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []resourceModelIp
	got.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].Ip.ValueString())
}

func TestInstance_GetLaunchInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceModel()

		got, err := instance.GetLaunchInstanceOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.REGIONNAME_EU_WEST_3, got.Region)
		assert.Equal(t, publicCloud.TYPENAME_M5A_4XLARGE, got.Type)
		assert.Equal(t, publicCloud.STORAGETYPE_CENTRAL, got.RootDiskStorageType)
		assert.Equal(t, "UBUNTU_20_04_64BIT", got.ImageId)
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.ContractType)
		assert.Equal(t, publicCloud.CONTRACTTERM__3, got.ContractTerm)
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__1, got.BillingFrequency)
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
			contract := generateContractObject(
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
			contract := generateContractObject(
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
			contract := generateContractObject(
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

		got, diags := instance.GetUpdateInstanceOpts(context.TODO())

		assert.Nil(t, diags)
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
			contract := generateContractObject(
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
			contract := generateContractObject(
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
			contract := generateContractObject(
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

			contract := generateContractObject(nil, nil, nil, &endsAt)

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

	want := resourceModelIp{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got, err := adaptSdkIpToResourceIp(context.TODO(), sdkIp)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_adaptSdkIpDetailsToResourceIp(t *testing.T) {
	sdkIpDetails := publicCloud.IpDetails{
		Ip: "127.0.0.1",
	}

	want := resourceModelIp{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got, err := adaptSdkIpDetailsToResourceIp(context.TODO(), sdkIpDetails)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_instanceResource_Metadata(t *testing.T) {
	resp := resource.MetadataResponse{}
	instanceResource := NewInstanceResource()

	instanceResource.Metadata(
		context.TODO(),
		resource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(t,
		"tralala_public_cloud_instance",
		resp.TypeName,
		"Type name should be tralala_public_cloud_instance",
	)
}
