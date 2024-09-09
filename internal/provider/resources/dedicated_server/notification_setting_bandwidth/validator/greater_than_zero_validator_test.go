package validator

import (
	"context"
	"testing"

	schemaValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestGreaterThanZeroValidator_ValidateString(t *testing.T) {
	t.Run("does not set errors if the value is greater than 0", func(t *testing.T) {
		request := schemaValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("2"),
		}

		response := schemaValidator.StringResponse{}

		validator := GreaterThanZero()
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run("set errors if the value is 0", func(t *testing.T) {
		request := schemaValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("0"),
		}

		response := schemaValidator.StringResponse{}

		validator := GreaterThanZero()
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})

	t.Run("set errors if the value is less than 0", func(t *testing.T) {
		request := schemaValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("-1"),
		}

		response := schemaValidator.StringResponse{}

		validator := GreaterThanZero()
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})

	t.Run("set errors if the value is any string", func(t *testing.T) {
		request := schemaValidator.StringRequest{
			ConfigValue: basetypes.NewStringValue("test"),
		}

		response := schemaValidator.StringResponse{}

		validator := GreaterThanZero()
		validator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})
}
