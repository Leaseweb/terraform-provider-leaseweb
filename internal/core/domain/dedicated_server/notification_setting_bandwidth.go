package dedicated_server

type NotificationSettingBandwidth struct {
	ServerId            string
	Id                  string
	Frequency           string
	LastCheckedAt       *string
	Threshold           string
	ThresholdExceededAt *string
	Unit                string
	Actions             Actions
}

type OptionalNotificationSettingBandwidthValues struct {
	LastCheckedAt       *string
	ThresholdExceededAt *string
}

func NewNotificationSettingBandwidth(
	serverId string,
	id string,
	frequency string,
	threshold string,
	unit string,
	actions Actions,
	optional OptionalNotificationSettingBandwidthValues,
) NotificationSettingBandwidth {
	notificationSettingBandwidth := NotificationSettingBandwidth{
		ServerId:  serverId,
		Id:        id,
		Frequency: frequency,
		Threshold: threshold,
		Unit:      unit,
		Actions:   actions,
	}

	if optional.LastCheckedAt != nil {
		notificationSettingBandwidth.LastCheckedAt = optional.LastCheckedAt
	}

	if optional.ThresholdExceededAt != nil {
		notificationSettingBandwidth.ThresholdExceededAt = optional.ThresholdExceededAt
	}

	return notificationSettingBandwidth
}

// NewCreateNotificationSettingBandwidth creates notification setting bandwidth with only all the supported fields.
func NewCreateNotificationSettingBandwidth(
	serverId string,
	frequency string,
	threshold string,
	unit string,
) *NotificationSettingBandwidth {
	notificationSettingBandwidth := NotificationSettingBandwidth{
		ServerId:  serverId,
		Frequency: frequency,
		Threshold: threshold,
		Unit:      unit,
	}

	return &notificationSettingBandwidth
}
