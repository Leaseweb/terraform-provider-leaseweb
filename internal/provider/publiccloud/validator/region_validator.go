package validator

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = RegionValidator{}

// RegionValidator validates if a region exists.
type RegionValidator struct {
	regions []string
}

func (r RegionValidator) Description(ctx context.Context) string {
	return `Determines whether a region exists`
}

func (r RegionValidator) MarkdownDescription(ctx context.Context) string {
	return r.Description(ctx)
}

func (r RegionValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	// If the region is unknown or null, there is nothing to validate.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	regionExists := slices.Contains(r.regions, request.ConfigValue.ValueString())

	if !regionExists {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Region",
			fmt.Sprintf(
				"Attribute region value must be one of: %q, got: %q",
				r.regions,
				request.ConfigValue.ValueString(),
			),
		)
	}
}

func NewRegionValidator(regions []string) RegionValidator {
	return RegionValidator{regions: regions}
}
