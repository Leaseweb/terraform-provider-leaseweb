package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
)

var _ validator.Object = InstanceTerminationValidator{}

// InstanceTerminationValidator validates if the Instance is allowed to be terminated.
type InstanceTerminationValidator struct {
	canInstanceBeTerminated func(
		id string,
		ctx context.Context,
	) (bool, *string, *errors.ServiceError)
}

func (i InstanceTerminationValidator) Description(ctx context.Context) string {
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

func (i InstanceTerminationValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i InstanceTerminationValidator) ValidateObject(
	ctx context.Context,
	request validator.ObjectRequest,
	response *validator.ObjectResponse,
) {
	instance := model.Instance{}

	diags := request.ConfigValue.As(ctx, &instance, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	_, reason, err := i.canInstanceBeTerminated(
		instance.Id.ValueString(),
		ctx,
	)
	if err != nil {
		response.Diagnostics.AddError("ValidateObject", err.Error())
		return
	}

	if reason != nil {
		response.Diagnostics.AddError(
			"Instance is not allowed to be terminated",
			*reason,
		)
	}
}

func ValidateInstanceTermination(
	canInstanceBeTerminated func(
		id string,
		ctx context.Context,
	) (bool, *string, *errors.ServiceError),
) InstanceTerminationValidator {
	return InstanceTerminationValidator{
		canInstanceBeTerminated: canInstanceBeTerminated,
	}
}
