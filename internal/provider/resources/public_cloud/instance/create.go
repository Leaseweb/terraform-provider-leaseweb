package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-leaseweb/internal/provider/logging"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func (i *instanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan model.Instance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating public cloud instance")
	instance, err := i.client.PublicCloudHandler.CreateInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Instance", err.Error())

		logging.HandleError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Error creating public cloud instance",
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
