package data_traffic_notification_setting

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

func (d *dataTrafficNotificationSettingResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan model.DataTrafficNotificationSetting

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Updating data traffic notification setting %q - %q",
		plan.ServerId.ValueString(),
		plan.Id.ValueString(),
	))
	updatedDataTrafficNotificationSetting, err := d.client.DedicatedServerFacade.UpdateDataTrafficNotificationSetting(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error data traffic notification setting", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to update data traffic notification setting %q - %q",
				plan.ServerId.ValueString(),
				plan.Id.ValueString(),
			),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, updatedDataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
