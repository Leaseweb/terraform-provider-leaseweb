package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := value_object.NewUuid(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error converting id to uuid", err.Error())
		return
	}

	instance, err := i.client.PublicCloud.GetInstance(*id, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Instance", err.Error())
		return
	}

	state.Populate(*instance, ctx)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
