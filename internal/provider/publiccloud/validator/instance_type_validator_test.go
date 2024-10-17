package validator

import (
	"context"
	"testing"

	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestInstanceTypeValidator_ValidateString(t *testing.T) {
	t.Run("nothing happens if instanceType is unknown", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		validator := InstanceTypeValidator{}

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

		validator := InstanceTypeValidator{}

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
			validator := InstanceTypeValidator{
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
			validator := InstanceTypeValidator{
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
	validator := NewInstanceTypeValidator(
		basetypes.NewStringValue("currentInstanceType"),
		[]string{"type1"},
	)

	assert.Equal(
		t,
		[]string{"type1", "currentInstanceType"},
		validator.availableInstanceTypes,
	)
}
