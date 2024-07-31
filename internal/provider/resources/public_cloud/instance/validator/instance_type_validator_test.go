package validator

import (
	"context"
	"testing"

	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestInstanceTypeValidator_setError(t *testing.T) {
	validator := InstanceTypeValidator{}

	response := terraformValidator.StringResponse{}
	request := terraformValidator.StringRequest{
		ConfigValue: basetypes.NewStringValue("tralala"),
	}

	validator.setError(&response, request, []string{"piet"})

	assert.Len(t, response.Diagnostics.Errors(), 1)
	assert.Contains(
		t,
		response.Diagnostics.Errors()[0].Detail(),
		"tralala",
	)
	assert.Contains(
		t,
		response.Diagnostics.Errors()[0].Detail(),
		"piet",
	)
}

func TestInstanceTypeValidator_ValidateString(t *testing.T) {
	t.Run("nothing happens if instanceType is unknown", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		validator := InstanceTypeValidator{
			isInstanceTypeAvailableForRegion: func(
				instanceType string,
				region string,
				ctx context.Context,
			) (bool, []string, error) {
				countIsInstanceTypeAvailableForRegionIsCalled++

				return false, nil, nil
			},
			canInstanceTypeBeUsedWithInstance: func(
				id string,
				instanceType string,
				ctx context.Context,
			) (bool, []string, error) {
				countCanInstanceTypeBeUsedWithInstanceIsCalled++

				return false, nil, nil
			},
		}

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

		validator := InstanceTypeValidator{
			isInstanceTypeAvailableForRegion: func(
				instanceType string,
				region string,
				ctx context.Context,
			) (bool, []string, error) {
				countIsInstanceTypeAvailableForRegionIsCalled++

				return false, nil, nil
			},
			canInstanceTypeBeUsedWithInstance: func(
				id string,
				instanceType string,
				ctx context.Context,
			) (bool, []string, error) {
				countCanInstanceTypeBeUsedWithInstanceIsCalled++

				return false, nil, nil
			},
		}

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
		"logic for created instance is followed if instanceId is null",
		func(t *testing.T) {
			countIsInstanceTypeAvailableForRegionIsCalled := 0

			validator := InstanceTypeValidator{
				isInstanceTypeAvailableForRegion: func(
					instanceType string,
					region string,
					ctx context.Context,
				) (bool, []string, error) {
					countIsInstanceTypeAvailableForRegionIsCalled++

					return false, nil, nil
				},
				instanceId: basetypes.NewStringNull(),
			}

			response := terraformValidator.StringResponse{}
			validator.ValidateString(
				context.TODO(),
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("oldInstanceType"),
				},
				&response,
			)

			assert.Equal(t, 1, countIsInstanceTypeAvailableForRegionIsCalled)
		},
	)

	t.Run(
		"logic for updated instance is followed if instanceId is known",
		func(t *testing.T) {
			countCanInstanceTypeBeUsedWithInstance := 0

			validator := InstanceTypeValidator{
				canInstanceTypeBeUsedWithInstance: func(
					id string,
					instanceType string,
					ctx context.Context,
				) (bool, []string, error) {
					countCanInstanceTypeBeUsedWithInstance++

					return false, nil, nil
				},
				instanceId: basetypes.NewStringValue(""),
			}

			response := terraformValidator.StringResponse{}
			validator.ValidateString(
				context.TODO(),
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("oldInstanceType"),
				},
				&response,
			)

			assert.Equal(t, 1, countCanInstanceTypeBeUsedWithInstance)
		},
	)
}

func TestInstanceTypeValidator_validateCreatedInstance(t *testing.T) {
	t.Run(
		"instanceType & region are passed to isInstanceTypeAvailableForRegion",
		func(t *testing.T) {
			validator := InstanceTypeValidator{
				isInstanceTypeAvailableForRegion: func(
					instanceType string,
					region string,
					ctx context.Context,
				) (bool, []string, error) {
					assert.Equal(t, "region", region)
					assert.Equal(t, "instanceType", instanceType)

					return false, nil, nil
				},
				region: basetypes.NewStringValue("region"),
			}

			response := terraformValidator.StringResponse{}
			validator.validateCreatedInstance(
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("instanceType"),
				},
				&response,
				context.TODO(),
			)
		},
	)

	t.Run(
		"no errors are set if instanceType is valid",
		func(t *testing.T) {
			validator := InstanceTypeValidator{
				isInstanceTypeAvailableForRegion: func(
					instanceType string,
					region string,
					ctx context.Context,
				) (bool, []string, error) {
					return true, nil, nil
				},
			}

			response := terraformValidator.StringResponse{}
			validator.validateCreatedInstance(
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue(""),
				},
				&response,
				context.TODO(),
			)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"errors are set if instanceType is valid for creation",
		func(t *testing.T) {
			validator := InstanceTypeValidator{
				isInstanceTypeAvailableForRegion: func(
					instanceType string,
					region string,
					ctx context.Context,
				) (bool, []string, error) {
					return false, []string{"tralala"}, nil
				},
			}

			response := terraformValidator.StringResponse{}
			validator.validateCreatedInstance(
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("piet"),
				},
				&response,
				context.TODO(),
			)

			assert.Len(t, response.Diagnostics.Errors(), 1)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"tralala",
			)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"piet",
			)
		},
	)
}

func TestInstanceTypeValidator_validateUpdatedInstance(t *testing.T) {
	t.Run(
		"no errors are set if instanceType is valid ",
		func(t *testing.T) {
			validator := InstanceTypeValidator{
				canInstanceTypeBeUsedWithInstance: func(
					id string,
					instanceType string,
					ctx context.Context,
				) (bool, []string, error) {
					return true, nil, nil
				},
			}

			response := terraformValidator.StringResponse{}
			validator.validateUpdatedInstance(
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue(""),
				},
				&response,
				context.TODO(),
			)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"instanceType & id are passed to canInstanceTypeBeUsedWithInstance",
		func(t *testing.T) {
			validator := InstanceTypeValidator{
				canInstanceTypeBeUsedWithInstance: func(
					id string,
					instanceType string,
					ctx context.Context,
				) (bool, []string, error) {
					assert.Equal(t, "instanceType", instanceType)
					assert.Equal(t, "instanceId", id)

					return true, nil, nil
				},
				instanceId: basetypes.NewStringValue("instanceId"),
			}

			response := terraformValidator.StringResponse{}
			validator.validateUpdatedInstance(
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("instanceType"),
				},
				&response,
				context.TODO(),
			)
		},
	)

	t.Run(
		"errors are set if instanceType is not valid",
		func(t *testing.T) {
			validator := InstanceTypeValidator{
				canInstanceTypeBeUsedWithInstance: func(
					id string,
					instanceType string,
					ctx context.Context,
				) (bool, []string, error) {
					return false, []string{"tralala"}, nil
				},
			}

			response := terraformValidator.StringResponse{}
			validator.validateUpdatedInstance(
				terraformValidator.StringRequest{
					ConfigValue: basetypes.NewStringValue("piet"),
				},
				&response,
				context.TODO(),
			)

			assert.Len(t, response.Diagnostics.Errors(), 1)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"tralala",
			)
			assert.Contains(
				t, response.Diagnostics.Errors()[0].Detail(),
				"piet",
			)
		},
	)
}

func TestNewInstanceTypeValidator(t *testing.T) {
	validator := NewInstanceTypeValidator(
		func(
			instanceType string,
			region string,
			ctx context.Context,
		) (bool, []string, error) {
			return false, []string{"tralala"}, nil
		},
		func(
			id string,
			instanceType string,
			ctx context.Context,
		) (bool, []string, error) {
			return false, []string{"blah"}, nil
		},
		basetypes.NewStringValue("instanceId"),
		basetypes.NewStringValue("region"),
	)

	assert.Equal(t, "instanceId", validator.instanceId.ValueString())
	assert.Equal(t, "region", validator.region.ValueString())

	_, got, _ := validator.canInstanceTypeBeUsedWithInstance(
		"",
		"",
		context.TODO(),
	)
	assert.Equal(t, []string{"blah"}, got)

	_, got, _ = validator.isInstanceTypeAvailableForRegion(
		"",
		"",
		context.TODO(),
	)
	assert.Equal(t, []string{"tralala"}, got)
}
