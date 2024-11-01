package dedicatedserver

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func Test_greaterThanZeroValidator_ValidateString(t *testing.T) {
	t.Run("does not set errors if the int value is greater than 0", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("2"),
		}

		response := validator.StringResponse{}

		greaterThanZeroValidator := greaterThanZero()
		greaterThanZeroValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run("does not set errors if the float value is greater than 0", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("2.3"),
		}

		response := validator.StringResponse{}

		greaterThanZeroValidator := greaterThanZero()
		greaterThanZeroValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 0)
	})

	t.Run("set errors if the value is 0", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("0"),
		}

		response := validator.StringResponse{}

		greaterThanZeroValidator := greaterThanZero()
		greaterThanZeroValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})

	t.Run("set errors if the int value is less than 0", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("-1"),
		}

		response := validator.StringResponse{}

		greaterThanZeroValidator := greaterThanZero()
		greaterThanZeroValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})

	t.Run("set errors if the float value is less than 0", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("-1.1"),
		}

		response := validator.StringResponse{}

		greaterThanZeroValidator := greaterThanZero()
		greaterThanZeroValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})

	t.Run("set errors if the value is any string", func(t *testing.T) {
		request := validator.StringRequest{
			ConfigValue: basetypes.NewStringValue("test"),
		}

		response := validator.StringResponse{}

		greaterThanZeroValidator := greaterThanZero()
		greaterThanZeroValidator.ValidateString(context.TODO(), request, &response)

		assert.Len(t, response.Diagnostics.Errors(), 1)
	})
}
