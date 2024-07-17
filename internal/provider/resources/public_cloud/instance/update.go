package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func (i *instanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan model.Instance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedInstance, err := i.client.PublicCloudHandler.UpdateInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error updating instance", err.Error())
		return
	}

	diags = resp.State.Set(ctx, updatedInstance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
