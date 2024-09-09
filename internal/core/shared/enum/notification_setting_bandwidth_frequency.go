package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type NotificationSettingBandwidthFrequency string

func (f NotificationSettingBandwidthFrequency) String() string {
	return string(f)
}

func (f NotificationSettingBandwidthFrequency) Values() []string {
	return enum_utils.ConvertStringEnumToValues(notificationSettingBandwidthFrequencyValues)
}

const (
	NotificationSettingBandwidthFrequencyDaily   NotificationSettingBandwidthFrequency = "DAILY"
	NotificationSettingBandwidthFrequencyWeekly  NotificationSettingBandwidthFrequency = "WEEKLY"
	NotificationSettingBandwidthFrequencyMonthly NotificationSettingBandwidthFrequency = "MONTHLY"
)

var notificationSettingBandwidthFrequencyValues = []NotificationSettingBandwidthFrequency{
	NotificationSettingBandwidthFrequencyDaily,
	NotificationSettingBandwidthFrequencyWeekly,
	NotificationSettingBandwidthFrequencyMonthly,
}

func NewNotificationSettingBandwidthFrequency(value string) (NotificationSettingBandwidthFrequency, error) {
	return enum_utils.FindEnumForString(value, notificationSettingBandwidthFrequencyValues, NotificationSettingBandwidthFrequencyDaily)
}
