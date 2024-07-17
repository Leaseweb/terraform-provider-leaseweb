package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
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

	instance, err := i.client.PublicCloudHandler.GetInstance(
		state.Id.ValueString(),
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Instance", err.Error())
		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
