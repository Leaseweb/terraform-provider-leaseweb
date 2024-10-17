package validator

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = InstanceTypeValidator{}

type InstanceTypeValidator struct {
	availableInstanceTypes []string
}

func (i InstanceTypeValidator) Description(ctx context.Context) string {
	return "Determines if an instanceType can be used with an instance."
}

func (i InstanceTypeValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i InstanceTypeValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	// Nothing to validate here.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	if !slices.Contains(
		i.availableInstanceTypes,
		request.ConfigValue.ValueString(),
	) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Instance Type",
			fmt.Sprintf(
				"Attribute type value must be one of: %q, got: %q",
				i.availableInstanceTypes,
				request.ConfigValue.ValueString(),
			),
		)
	}
}

func NewInstanceTypeValidator(
	currentInstanceType types.String,
	availableInstanceTypes []string,
) InstanceTypeValidator {
	// Include the current instance type as it isn't returned the by api.
	availableInstanceTypes = append(availableInstanceTypes, currentInstanceType.ValueString())

	return InstanceTypeValidator{
		availableInstanceTypes: availableInstanceTypes,
	}
}
