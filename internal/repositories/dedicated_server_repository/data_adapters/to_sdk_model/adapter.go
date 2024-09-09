// Package to_sdk_model implements adapters to convert dedicated_server domain entities to sdk models.
package to_sdk_model

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
)

// AdaptToCreateNotificationSettingBandwidthOpts adapts dedicated_server.NotificationSettingBandwidth to dedicatedServer.BandwidthNotificationSettingOpts.
func AdaptToCreateNotificationSettingBandwidthOpts(notificationSettingBandwidth dedicated_server.NotificationSettingBandwidth) *dedicatedServer.BandwidthNotificationSettingOpts {

	return dedicatedServer.NewBandwidthNotificationSettingOpts(
		notificationSettingBandwidth.Frequency,
		notificationSettingBandwidth.Threshold,
		notificationSettingBandwidth.Unit,
	)
}
