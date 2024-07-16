package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
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

	instance, generateCreateInstanceEntityDiags := plan.GenerateCreateInstanceEntity(ctx)
	if generateCreateInstanceEntityDiags.HasError() {
		for _, diag := range generateCreateInstanceEntityDiags.Errors() {
			resp.Diagnostics.Append(diag)
		}

		return
	}

	instance, err := i.client.PublicCloud.CreateInstance(*instance, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Instance",
			err.Error(),
		)
		return
	}

	diags = plan.Populate(*instance, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
