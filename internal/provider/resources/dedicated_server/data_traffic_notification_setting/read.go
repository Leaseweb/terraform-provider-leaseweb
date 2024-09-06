package data_traffic_notification_setting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

func (d *dataTrafficNotificationSettingResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.DataTrafficNotificationSetting
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Read dedicated server data traffic notification setting %q - %q",
		state.ServerId.ValueString(),
		state.Id.ValueString(),
	))
	dataTrafficNotificationSetting, err := d.client.DedicatedServerFacade.GetDataTrafficNotificationSetting(
		ctx,
		state.ServerId.ValueString(),
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading dedicated server data traffic notification setting", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to read dedicated server data traffic notification setting %q - %q", state.ServerId.ValueString(), state.Id.ValueString()),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
