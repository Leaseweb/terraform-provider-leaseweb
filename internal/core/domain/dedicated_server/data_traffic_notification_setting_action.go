package dedicated_server

import "time"

type DataTrafficNotificationSettingAction struct {
	LastTriggeredAt time.Time
	Type            string
}

func NewDataTrafficNotificationSettingAction(Type string, lastTriggeredAt time.Time) DataTrafficNotificationSettingAction {
	return DataTrafficNotificationSettingAction{
		Type:            Type,
		LastTriggeredAt: lastTriggeredAt,
	}
}
