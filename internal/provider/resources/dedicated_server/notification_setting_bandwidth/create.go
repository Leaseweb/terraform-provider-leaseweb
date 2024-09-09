package notification_setting_bandwidth

import (
	"context"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (n *notificationSettingBandwidthResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan model.NotificationSettingBandwidth

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating notification setting bandwidth")
	notificationSettingBandwidth, err := n.client.DedicatedServerFacade.CreateNotificationSettingBandwidth(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating notification setting bandwidth", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Error while creating notification setting bandwidth",
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, notificationSettingBandwidth)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
