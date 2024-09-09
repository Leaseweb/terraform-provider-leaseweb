package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewNotificationSettingBandwidth(t *testing.T) {
	date := "2024-03-02"
	got := NewNotificationSettingBandwidth(
		"12345",
		"123456",
		"DAILY",
		"1",
		"Mbps",
		[]Action{},
		OptionalNotificationSettingBandwidthValues{
			LastCheckedAt:       &date,
			ThresholdExceededAt: &date,
		},
	)
	want := NotificationSettingBandwidth{
		ServerId:            "12345",
		Id:                  "123456",
		Frequency:           "DAILY",
		Threshold:           "1",
		Unit:                "Mbps",
		LastCheckedAt:       &date,
		ThresholdExceededAt: &date,
		Actions:             []Action{},
	}
	assert.Equal(t, want, got)
}

func Test_NewCreateNotificationSettingBandwidth(t *testing.T) {
	got := NewCreateNotificationSettingBandwidth("12345", "DAILY", "1", "Mbps")

	assert.Equal(t, "12345", got.ServerId)
	assert.Equal(t, "DAILY", got.Frequency)
	assert.Equal(t, "1", got.Threshold)
	assert.Equal(t, "Mbps", got.Unit)
}
