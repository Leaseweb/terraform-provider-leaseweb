package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

			request := validator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := validator.ObjectResponse{}

			contractTermValidator := contractTermValidator{}
			contractTermValidator.ValidateObject(context.TODO(), request, &response)

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

			request := validator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := validator.ObjectResponse{}

			contractTermValidator := contractTermValidator{}
			contractTermValidator.ValidateObject(context.TODO(), request, &response)

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

			request := validator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := validator.ObjectResponse{}

			contractTermValidator := contractTermValidator{}
			contractTermValidator.ValidateObject(context.TODO(), request, &response)

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
		request := validator.ObjectRequest{}
		response := validator.ObjectResponse{}

		instanceTerminationValidator := instanceTerminationValidator{}
		instanceTerminationValidator.ValidateObject(context.TODO(), request, &response)

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
			request := validator.ObjectRequest{ConfigValue: instanceObject}
			response := validator.ObjectResponse{}

			instanceTerminationValidator := instanceTerminationValidator{}
			instanceTerminationValidator.ValidateObject(context.TODO(), request, &response)

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
			request := validator.ObjectRequest{ConfigValue: instanceObject}
			response := validator.ObjectResponse{}

			instanceTerminationValidator := instanceTerminationValidator{}
			instanceTerminationValidator.ValidateObject(context.TODO(), request, &response)

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

		instanceTypeValidator := instanceTypeValidator{}

		response := validator.StringResponse{}
		instanceTypeValidator.ValidateString(
			context.TODO(),
			validator.StringRequest{ConfigValue: basetypes.NewStringUnknown()},
			&response,
		)

		assert.Equal(t, 0, countIsInstanceTypeAvailableForRegionIsCalled)
		assert.Equal(t, 0, countCanInstanceTypeBeUsedWithInstanceIsCalled)
	})

	t.Run("nothing happens if instanceType does not change", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		instanceTypeValidator := instanceTypeValidator{}

		response := validator.StringResponse{}
		instanceTypeValidator.ValidateString(
			context.TODO(),
			validator.StringRequest{
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
			instanceTypeValidator := instanceTypeValidator{
				availableInstanceTypes: []string{"tralala"},
			}

			response := validator.StringResponse{}
			instanceTypeValidator.ValidateString(
				context.TODO(),
				validator.StringRequest{
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
			instanceTypeValidator := instanceTypeValidator{
				availableInstanceTypes: []string{"tralala"},
			}

			response := validator.StringResponse{}
			instanceTypeValidator.ValidateString(
				context.TODO(),
				validator.StringRequest{
					ConfigValue: basetypes.NewStringValue("tralala"),
				},
				&response,
			)

			assert.Len(t, response.Diagnostics, 0)
		},
	)
}

func Test_newInstanceTypeValidator(t *testing.T) {
	instanceTypeValidator := newInstanceTypeValidator(
		basetypes.NewStringValue("currentInstanceType"),
		[]string{"type1"},
	)

	assert.Equal(
		t,
		[]string{"type1", "currentInstanceType"},
		instanceTypeValidator.availableInstanceTypes,
	)
}

func Test_regionValidator_ValidateString(t *testing.T) {
	t.Run("does not set errors if the region exists", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := validator.StringResponse{}

		regionValidator := regionValidator{
			regions: []string{"region"},
		}
		regionValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run(
		"does not set errors if the region is unknown",
		func(t *testing.T) {
			request := validator.StringRequest{
				ConfigValue: basetypes.NewStringUnknown(),
			}

			response := validator.StringResponse{}

			regionValidator := regionValidator{}
			regionValidator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set errors if the region is null",
		func(t *testing.T) {
			request := validator.StringRequest{
				ConfigValue: basetypes.NewStringNull(),
			}

			response := validator.StringResponse{}

			regionValidator := regionValidator{}
			regionValidator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run("sets an error if the region does not exist", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := validator.StringResponse{}

		regionValidator := regionValidator{
			regions: []string{"tralala"},
		}

		regionValidator.ValidateString(context.TODO(), request, &response)

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
