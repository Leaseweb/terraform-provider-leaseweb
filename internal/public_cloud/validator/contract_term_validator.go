package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

var _ validator.Object = contractTermValidator{}

// contractTermValidator validates that the contract term is correct.
type contractTermValidator struct {
}

func (v contractTermValidator) Description(_ context.Context) string {
	return "Contract Term Validator"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v contractTermValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v contractTermValidator) ValidateObject(
	ctx context.Context,
	request validator.ObjectRequest,
	response *validator.ObjectResponse,
) {

	contract := model.Contract{}
	request.ConfigValue.As(ctx, &contract, basetypes.ObjectAsOptions{})

	if contract.Type.ValueString() == "MONTHLY" && contract.Term.ValueInt64() == 0 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path.AtName("term"),
			"cannot be 0 when contract.type is \"MONTHLY\"",
			contract.Term.String(),
		))
	}

	if contract.Type.ValueString() == "HOURLY" && contract.Term.ValueInt64() != 0 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path.AtName("term"),
			"must be 0 when contract.type is \"HOURLY\"",
			contract.Term.String(),
		))
	}
}

func ContractTermIsValid() validator.Object {
	return contractTermValidator{}
}
