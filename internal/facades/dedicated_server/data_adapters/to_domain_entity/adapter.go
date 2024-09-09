// Package to_domain_entity implements adapters to convert resource models to domain entities.
package to_domain_entity

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

// AdaptToCreateNotificationSettingBandwidthOpts transforms model.NotificationSettingBandwidth to dedicated_server.NotificationSettingBandwidth
// entity with all supported fields for creating notification setting bandwidth.
func AdaptToCreateNotificationSettingBandwidthOpts(
	resourceModel model.NotificationSettingBandwidth,
) *dedicated_server.NotificationSettingBandwidth {
	createNotificationSettingBandwidthOpts := dedicated_server.NewCreateNotificationSettingBandwidth(
		resourceModel.ServerId.ValueString(),
		resourceModel.Frequency.ValueString(),
		resourceModel.Threshold.ValueString(),
		resourceModel.Unit.ValueString(),
	)

	return createNotificationSettingBandwidthOpts
}
