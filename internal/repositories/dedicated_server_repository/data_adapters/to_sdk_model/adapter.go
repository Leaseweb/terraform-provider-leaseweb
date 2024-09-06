// Package to_sdk_model implements adapters to convert public_cloud domain entities to sdk models.
package to_sdk_model

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
)

// AdaptToDataTrafficNotificationSettingOpts adapts a DataTrafficNotificationSetting domain entity to supported DataTrafficNotificationSetting opts.
func AdaptToDataTrafficNotificationSettingOpts(dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) dedicatedServer.DataTrafficNotificationSettingOpts {
	return *dedicatedServer.NewDataTrafficNotificationSettingOpts(
		dataTrafficNotificationSetting.Frequency,
		dataTrafficNotificationSetting.Threshold,
		dataTrafficNotificationSetting.Unit,
	)
}
