package data_traffic_notification_setting

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (d *dataTrafficNotificationSettingResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_data_traffic_notification_setting"
}
