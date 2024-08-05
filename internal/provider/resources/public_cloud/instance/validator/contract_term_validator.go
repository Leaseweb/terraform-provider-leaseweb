package validator

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
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
	facade := public_cloud.PublicCloudFacade{}

	contract := model.Contract{}
	request.ConfigValue.As(ctx, &contract, basetypes.ObjectAsOptions{})
	err := facade.ValidateContractTerm(
		contract.Term.ValueInt64(),
		contract.Type.ValueString(),
	)

	if err != nil {
		switch {
		case errors.Is(err, public_cloud.ErrContractTermCannotBeZero):
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				request.Path.AtName("term"),
				"cannot be 0 when contract.type is \"MONTHLY\"",
				contract.Term.String(),
			))
			return
		case errors.Is(err, public_cloud.ErrContractTermMustBeZero):
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
