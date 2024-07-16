package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

func (i *instanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := value_object.NewUuid(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			fmt.Sprintf(
				"Could not convert id %q to uuid: %q",
				state.Id.ValueString(),
				err.Error(),
			),
		)
		return
	}

	err = i.client.PublicCloud.DeleteInstance(*id, ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			"Could not terminate Public CLoud Instance, unexpected error: "+err.Error(),
		)
		return
	}
}
