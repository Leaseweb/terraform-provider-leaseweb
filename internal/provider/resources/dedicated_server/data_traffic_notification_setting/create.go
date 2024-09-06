package data_traffic_notification_setting

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

func (d *dataTrafficNotificationSettingResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan model.DataTrafficNotificationSetting

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating data traffic notification setting")
	dataTrafficNotificationSetting, err := d.client.DedicatedServerFacade.CreateDataTrafficNotificationSetting(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating data traffic notification setting", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Error creating data traffic notification setting",
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
