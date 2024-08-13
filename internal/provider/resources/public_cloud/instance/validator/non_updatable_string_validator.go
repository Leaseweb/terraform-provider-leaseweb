package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = NonUpdatableStringValidator{}

// NonUpdatableStringValidator makes sure that a stringValue cannot be changed on update.
type NonUpdatableStringValidator struct {
	stateValue types.String
}

func (n NonUpdatableStringValidator) Description(ctx context.Context) string {
	return "Makes sure that a stringValue cannot be changed on update."
}

func (n NonUpdatableStringValidator) MarkdownDescription(ctx context.Context) string {
	return n.Description(ctx)
}

func (n NonUpdatableStringValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	// If the value is unknown or null, there is nothing to validate.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	// On create the value can be anything
	if n.stateValue.IsUnknown() || n.stateValue.IsNull() {
		return
	}

	if n.stateValue != request.ConfigValue {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid value",
			fmt.Sprintf(
				"Attribute value is not allowed to change, was %q got: %q",
				n.stateValue.ValueString(),
				request.ConfigValue.ValueString(),
			),
		)
	}
}

func NewNonUpdatableStringValidator(stateValue types.String) NonUpdatableStringValidator {
	return NonUpdatableStringValidator{stateValue: stateValue}
}
