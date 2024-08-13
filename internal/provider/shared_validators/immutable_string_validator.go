package shared_validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = ImmutableStringValidator{}

// ImmutableStringValidator makes sure that a stringValue cannot be changed on update.
type ImmutableStringValidator struct {
	stateIdValue types.String
	stateValue   types.String
}

func (n ImmutableStringValidator) Description(ctx context.Context) string {
	return "Makes sure that a stringValue cannot be changed on update."
}

func (n ImmutableStringValidator) MarkdownDescription(ctx context.Context) string {
	return n.Description(ctx)
}

func (n ImmutableStringValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	// If the value is unknown or null, there is nothing to validate.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	// On create the value can be anything
	if n.stateIdValue.IsUnknown() || n.stateIdValue.IsNull() {
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

func NewImmutableStringValidator(
	stateIdValue types.String,
	stateValue types.String,
) ImmutableStringValidator {
	return ImmutableStringValidator{
		stateIdValue: stateIdValue,
		stateValue:   stateValue,
	}
}
