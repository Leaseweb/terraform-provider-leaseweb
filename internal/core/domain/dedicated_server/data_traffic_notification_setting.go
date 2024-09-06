package dedicated_server

type DataTrafficNotificationSetting struct {
	Id        string
	Frequency string
	//LastCheckedAt       time.Time
	Threshold string
	//ThresholdExceededAt time.Time
	Unit string
	//Actions             DataTrafficNotificationSettingActions
}

func NewDataTrafficNotificationSetting(
	id string,
	frequency string,
	//lastCheckedAt time.Time,
	threshold string,
	//thresholdExceededAt time.Time,
	unit string,
	// actions DataTrafficNotificationSettingActions,
) DataTrafficNotificationSetting {
	return DataTrafficNotificationSetting{
		Id:        id,
		Frequency: frequency,
		//LastCheckedAt:       lastCheckedAt,
		Threshold: threshold,
		//ThresholdExceededAt: thresholdExceededAt,
		Unit: unit,
		//Actions:             actions,
	}
}

// TODO: these 2 functions can be merged together.
func NewCreateDataTrafficNotificationSetting(
	frequency string,
	threshold string,
	unit string,
) DataTrafficNotificationSetting {
	return DataTrafficNotificationSetting{
		Frequency: frequency,
		Threshold: threshold,
		Unit:      unit,
	}
}

func NewUpdateDataTrafficNotificationSetting(
	frequency string,
	threshold string,
	unit string,
) DataTrafficNotificationSetting {
	return DataTrafficNotificationSetting{
		Frequency: frequency,
		Threshold: threshold,
		Unit:      unit,
	}
}
