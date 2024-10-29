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
			contract := contractResourceModel{}
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
			contract := contractResourceModel{
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
			contract := contractResourceModel{
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
