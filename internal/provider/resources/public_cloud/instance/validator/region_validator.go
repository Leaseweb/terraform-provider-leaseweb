package validator

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-leaseweb/internal/handlers/shared"
)

var _ validator.String = RegionValidator{}

// RegionValidator validates if a region exists.
type RegionValidator struct {
	doesRegionExist func(
		region string,
		ctx context.Context,
	) (bool, []string, *shared.HandlerError)
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

	regionExists, currentRegions, err := r.doesRegionExist(
		request.ConfigValue.ValueString(),
		ctx,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	if !regionExists {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Region",
			fmt.Sprintf(
				"Attribute region value must be one of: %q, got: %q",
				currentRegions,
				request.ConfigValue.ValueString(),
			),
		)
	}
}

func NewRegionValidator(doesRegionExist func(
	region string,
	ctx context.Context,
) (bool, []string, *shared.HandlerError)) RegionValidator {
	return RegionValidator{doesRegionExist: doesRegionExist}
}
