package publiccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Object = InstanceTerminationValidator{}

// InstanceTerminationValidator validates if the ResourceModelInstance is allowed to be terminated.
type InstanceTerminationValidator struct{}

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
	instance := ResourceModelInstance{}

	diags := request.ConfigValue.As(ctx, &instance, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	reason := instance.CanBeTerminated(ctx)

	if reason != nil {
		response.Diagnostics.AddError(
			"ResourceModelInstance is not allowed to be terminated",
			string(*reason),
		)
	}
}
