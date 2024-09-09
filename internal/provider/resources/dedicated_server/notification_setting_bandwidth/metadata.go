package notification_setting_bandwidth

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (n *notificationSettingBandwidthResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_notification_setting_bandwidth"
}
