package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
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

	instance, diags := plan.GenerateUpdateInstanceEntity(ctx)
	if diags.HasError() {
		return
	}

	updatedInstance, err := i.client.PublicCloud.UpdateInstance(*instance, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error updating instance", err.Error())
		return
	}

	plan.Populate(*updatedInstance, ctx)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
