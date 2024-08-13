package shared_validators

import (
	"context"
	"testing"

	validator2 "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestImmutableStringValidator_ValidateString(t *testing.T) {
	t.Run(
		"does not set an error if current value is unknown",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringUnknown(),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("oldValue")
			stateIdValue := basetypes.NewStringValue("id")

			validator := NewImmutableStringValidator(stateIdValue, stateValue)
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
			stateIdValue := basetypes.NewStringValue("id")

			validator := NewImmutableStringValidator(stateIdValue, stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		})

	t.Run(
		"does not set an error if state id value is unknown",
		func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("tralala"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("")
			stateIdValue := basetypes.NewStringUnknown()

			validator := NewImmutableStringValidator(stateIdValue, stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does not set an error if state id value is null", func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("tralala"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("")
			stateIdValue := basetypes.NewStringNull()

			validator := NewImmutableStringValidator(stateIdValue, stateValue)
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
			stateIdValue := basetypes.NewStringValue("id")

			validator := NewImmutableStringValidator(stateIdValue, stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"does set an error if value changes on update", func(t *testing.T) {
			request := validator2.StringRequest{
				ConfigValue: basetypes.NewStringValue("oldValue"),
			}
			response := validator2.StringResponse{}
			stateValue := basetypes.NewStringValue("newValue")
			stateIdValue := basetypes.NewStringValue("id")

			validator := NewImmutableStringValidator(stateIdValue, stateValue)
			validator.ValidateString(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 1)
		})
}
