package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
)

var _ validator.Object = contractTermValidator{}

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
	contract := resource.Contract{}
	request.ConfigValue.As(ctx, &contract, basetypes.ObjectAsOptions{})
	valid, reason := contract.IsContractTermValid()

	if !valid {
		switch reason {
		case resource.ReasonContractTermCannotBeZero:
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				request.Path.AtName("term"),
				"cannot be 0 when contract.type is \"MONTHLY\"",
				contract.Term.String(),
			))
			return
		case resource.ReasonContractTermMustBeZero:
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

func ContractTermIsValid() validator.Object {
	return contractTermValidator{}
}
