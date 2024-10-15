package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/logging"
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

	tflog.Info(ctx, fmt.Sprintf(
		"Updating public cloud instance %q",
		plan.Id.ValueString(),
	))
	updatedInstance, err := i.client.PublicCloudService.UpdateInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error updating instance", err.Error())

		logging.ServiceError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to update public cloud instance %q",
				plan.Id.ValueString(),
			),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, updatedInstance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
