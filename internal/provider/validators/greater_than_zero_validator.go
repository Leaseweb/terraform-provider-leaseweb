package validator

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// GreaterThanZeroValidator ensures that the given value is greater than zero.
type GreaterThanZeroValidator struct{}

func (v GreaterThanZeroValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {

	if intValue, err := strconv.Atoi(request.ConfigValue.ValueString()); err == nil {
		if intValue <= 0 {
			response.Diagnostics.AddError(
				"Invalid Value",
				fmt.Sprintf("The value must be greater than 0, but got %s.", request.ConfigValue.ValueString()),
			)
		}
		return
	}

	if floatValue, err := strconv.ParseFloat(request.ConfigValue.ValueString(), 64); err == nil {
		if floatValue <= 0 {
			response.Diagnostics.AddError(
				"Invalid Value",
				fmt.Sprintf("The value must be greater than 0, but got %s.", request.ConfigValue.ValueString()),
			)
		}
		return
	}

	response.Diagnostics.AddError("Invalid Value", "The value is not a valid number")
}

var _ validator.String = GreaterThanZeroValidator{}

func (v GreaterThanZeroValidator) Description(ctx context.Context) string {
	return "Ensures that the value is greater than 0"
}

func (v GreaterThanZeroValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// GreaterThanZero returns a new instance of the validator.
func GreaterThanZero() validator.String {
	return GreaterThanZeroValidator{}
}
