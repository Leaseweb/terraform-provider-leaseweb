package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationSettingBandwidthFrequency_Value(t *testing.T) {
	got := NotificationSettingBandwidthFrequencyDaily.String()

	assert.Equal(t, "DAILY", got)
}

func TestNewNotificationSettingBandwidthFrequency(t *testing.T) {
	want := NotificationSettingBandwidthFrequencyDaily
	got, err := NewNotificationSettingBandwidthFrequency("DAILY")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNotificationSettingBandwidthFrequency_Values(t *testing.T) {
	want := []string{"DAILY", "WEEKLY", "MONTHLY"}
	got := NotificationSettingBandwidthFrequencyDaily.Values()

	assert.Equal(t, want, got)
}
