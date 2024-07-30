package validator

import (
	"context"
	"testing"

	validator2 "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/handlers/shared"
)

func TestRegionValidator_ValidateString(t *testing.T) {
	t.Run("does not set errors if the region exists", func(t *testing.T) {
		request := validator2.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := validator2.StringResponse{}

		validator := NewRegionValidator(
			func(
				region string,
				ctx context.Context,
			) (bool, []string, *shared.HandlerError) {
				return true, nil, nil
			},
		)
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run("passes region to handler", func(t *testing.T) {
		request := validator2.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := validator2.StringResponse{}

		validator := NewRegionValidator(
			func(
				region string,
				ctx context.Context,
			) (bool, []string, *shared.HandlerError) {
				assert.Equal(t, "region", region)
				return true, nil, nil
			},
		)
		validator.ValidateString(context.TODO(), request, &response)
	})

	t.Run(
		"does not set errors if the region is unknown",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringUnknown(),
			}

			response := validator2.StringResponse{}

			validator := RegionValidator{}
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set errors if the region is null",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringNull(),
			}

			response := validator2.StringResponse{}

			validator := RegionValidator{}
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run("sets an error if the region does not exist", func(t *testing.T) {
		request := validator2.StringRequest{
			ConfigValue: basetypes.NewStringValue("region"),
		}

		response := validator2.StringResponse{}

		validator := NewRegionValidator(
			func(
				region string,
				ctx context.Context,
			) (bool, []string, *shared.HandlerError) {
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
