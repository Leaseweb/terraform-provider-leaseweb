package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type NotificationSettingBandwidthUnit string

func (u NotificationSettingBandwidthUnit) String() string {
	return string(u)
}

func (u NotificationSettingBandwidthUnit) Values() []string {
	return enum_utils.ConvertStringEnumToValues(notificationSettingBandwidthUnitValues)
}

const (
	NotificationSettingBandwidthUnitMbps NotificationSettingBandwidthUnit = "Mbps"
	NotificationSettingBandwidthUnitGbps NotificationSettingBandwidthUnit = "Gbps"
)

var notificationSettingBandwidthUnitValues = []NotificationSettingBandwidthUnit{
	NotificationSettingBandwidthUnitMbps,
	NotificationSettingBandwidthUnitGbps,
}

func NewNotificationSettingBandwidthUnit(value string) (NotificationSettingBandwidthUnit, error) {
	return enum_utils.FindEnumForString(value, notificationSettingBandwidthUnitValues, NotificationSettingBandwidthUnitMbps)
}
