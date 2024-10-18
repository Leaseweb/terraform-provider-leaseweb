package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func Test_contractTermValidator_ValidateObject(t *testing.T) {
	t.Run(
		"does not set error if contract term is correct",
		func(t *testing.T) {
			contract := ResourceModelContract{}
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
			contract := ResourceModelContract{
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
			contract := ResourceModelContract{
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

func TestInstanceTerminationValidator_ValidateObject(t *testing.T) {
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

func generateInstanceModelForValidator() ResourceModelInstance {
	contract := ResourceModelContract{}
	contractObject, _ := types.ObjectValueFrom(
		context.TODO(),
		contract.AttributeTypes(),
		contract,
	)

	return ResourceModelInstance{
		Id:        basetypes.NewStringUnknown(),
		Region:    basetypes.NewStringUnknown(),
		Reference: basetypes.NewStringUnknown(),
		Image: basetypes.NewObjectUnknown(
			ResourceModelImage{}.AttributeTypes(),
		),
		State:               basetypes.NewStringUnknown(),
		Type:                basetypes.NewStringUnknown(),
		RootDiskSize:        basetypes.NewInt64Unknown(),
		RootDiskStorageType: basetypes.NewStringUnknown(),
		Ips: basetypes.NewListUnknown(
			types.ObjectType{
				AttrTypes: ResourceModelIp{}.AttributeTypes(),
			},
		),
		Contract:    contractObject,
		MarketAppId: basetypes.NewStringUnknown(),
	}
}

func TestInstanceTypeValidator_ValidateString(t *testing.T) {
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

func TestNewInstanceTypeValidator(t *testing.T) {
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

func TestRegionValidator_ValidateString(t *testing.T) {
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
