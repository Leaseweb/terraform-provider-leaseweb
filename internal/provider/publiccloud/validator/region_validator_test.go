package validator

import (
	"context"
	"testing"

	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
	"github.com/stretchr/testify/assert"
)

func TestRegionValidator_ValidateString(t *testing.T) {
	t.Run("does not set errors if the region exists", func(t *testing.T) {
		request := terraformValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := terraformValidator.StringResponse{}

		validator := NewRegionValidator(
			func(
				region string,
				ctx context.Context,
			) (bool, []string, *errors.ServiceError) {
				return true, nil, nil
			},
		)
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run("passes region to handler", func(t *testing.T) {
		request := terraformValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := terraformValidator.StringResponse{}

		validator := NewRegionValidator(
			func(
				region string,
				ctx context.Context,
			) (bool, []string, *errors.ServiceError) {
				assert.Equal(t, "region", region)
				return true, nil, nil
			},
		)
		validator.ValidateString(context.TODO(), request, &response)
	})

	t.Run(
		"does not set errors if the region is unknown",
		func(t *testing.T) {
			request := terraformValidator.StringRequest{
				ConfigValue: basetypes.NewStringUnknown(),
			}

			response := terraformValidator.StringResponse{}

			validator := RegionValidator{}
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

			validator := RegionValidator{}
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run("sets an error if the region does not exist", func(t *testing.T) {
		request := terraformValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := terraformValidator.StringResponse{}

		validator := NewRegionValidator(
			func(
				region string,
				ctx context.Context,
			) (bool, []string, *errors.ServiceError) {
				return false, []string{"tralala"}, nil
			},
		)

		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
		assert.Contains(
			t,
			response.Diagnostics.Errors()[0].Detail(),
			"tralala",
		)
	})
}
