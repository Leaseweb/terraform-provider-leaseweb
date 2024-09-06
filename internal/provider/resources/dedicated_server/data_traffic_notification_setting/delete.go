package data_traffic_notification_setting

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

func (d *dataTrafficNotificationSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.DataTrafficNotificationSetting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Deleting data traffic notification setting %q - %q",
		state.ServerId.ValueString(),
		state.Id.ValueString(),
	))
	err := d.client.DedicatedServerFacade.DeleteDataTrafficNotificationSetting(ctx, state.ServerId.ValueString(), state.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting data traffic notification setting",
			fmt.Sprintf(
				"Could not delete data traffic notification setting, unexpected error: %q",
				err.Error(),
			),
		)

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Error deleting data traffic notification setting %q - %q",
				state.ServerId.ValueString(),
				state.Id.ValueString(),
			),
			err.Error(),
		)

		return
	}
}
