package validator

import (
	"context"
	"testing"

	validator2 "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestNonUpdatableStringValidator_ValidateString(t *testing.T) {
	t.Run(
		"does not set an error if current value is unknown",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringUnknown(),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("oldValue")

			validator := NewNonUpdatableStringValidator(stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set an error if current value is null", func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringNull(),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("oldValue")

			validator := NewNonUpdatableStringValidator(stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		})

	t.Run(
		"does not set an error if state value is unknown",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("tralala"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringUnknown()

			validator := NewNonUpdatableStringValidator(stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set an error if state value is null", func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("tralala"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringNull()

			validator := NewNonUpdatableStringValidator(stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		})

	t.Run(
		"does not set an error if value does not change on update",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("tralala"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("tralala")

			validator := NewNonUpdatableStringValidator(stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set an error if value does change on update",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("oldValue"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("newValue")

			validator := NewNonUpdatableStringValidator(stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 1)
		},
	)
}
