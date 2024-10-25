package publiccloud

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ validator.Object = contractTermValidator{}
	_ validator.Object = instanceTerminationValidator{}
	_ validator.String = regionValidator{}
	_ validator.String = instanceTypeValidator{}
)

// Checks that contractType/contractTerm combination is valid.
type contractTermValidator struct {
}

func (v contractTermValidator) Description(_ context.Context) string {
	return `When contract.type is "MONTHLY", contract.term cannot be 0. When contract.type is "HOURLY", contract.term may only be 0.`
}

func (v contractTermValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v contractTermValidator) ValidateObject(
	ctx context.Context,
	request validator.ObjectRequest,
	response *validator.ObjectResponse,
) {
	contract := contractResourceModel{}
	request.ConfigValue.As(ctx, &contract, basetypes.ObjectAsOptions{})
	valid, reason := contract.IsContractTermValid()

	if !valid {
		switch reason {
		case reasonContractTermCannotBeZero:
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				request.Path.AtName("term"),
				"cannot be 0 when contract.type is \"MONTHLY\"",
				contract.Term.String(),
			))
			return
		case reasonContractTermMustBeZero:
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				request.Path.AtName("term"),
				"must be 0 when contract.type is \"HOURLY\"",
				contract.Term.String(),
			))
			return
		default:
			return
		}
	}
}

// instanceTerminationValidator validates if the instanceResourceModel is allowed to be terminated.
type instanceTerminationValidator struct{}

func (i instanceTerminationValidator) Description(_ context.Context) string {
	return `
Determines whether an instance can be terminated or not. An instance cannot be
terminated if:

- state is equal to Creating
- state is equal to Destroying
- state is equal to Destroyed
- contract.endsAt is set

In all other scenarios an instance can be terminated.
`
}

func (i instanceTerminationValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i instanceTerminationValidator) ValidateObject(
	ctx context.Context,
	request validator.ObjectRequest,
	response *validator.ObjectResponse,
) {
	instance := instanceResourceModel{}

	diags := request.ConfigValue.As(ctx, &instance, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	reason := instance.CanBeTerminated(ctx)

	if reason != nil {
		response.Diagnostics.AddError(
			"instance is not allowed to be terminated",
			string(*reason),
		)
	}
}

// regionValidator validates if a region exists.
type regionValidator struct {
	regions []string
}

func (r regionValidator) Description(_ context.Context) string {
	return `Determines whether a region exists`
}

func (r regionValidator) MarkdownDescription(ctx context.Context) string {
	return r.Description(ctx)
}

func (r regionValidator) ValidateString(
	_ context.Context,
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

type instanceTypeValidator struct {
	availableInstanceTypes []string
}

func (i instanceTypeValidator) Description(_ context.Context) string {
	return "Determines if an instanceType can be used with an instance."
}

func (i instanceTypeValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i instanceTypeValidator) ValidateString(
	_ context.Context,
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

func newInstanceTypeValidator(
	currentInstanceType types.String,
	availableInstanceTypes []string,
) instanceTypeValidator {
	// Include the current instance type as it isn't returned the by api.
	availableInstanceTypes = append(
		availableInstanceTypes,
		currentInstanceType.ValueString(),
	)

	return instanceTypeValidator{
		availableInstanceTypes: availableInstanceTypes,
	}
}
