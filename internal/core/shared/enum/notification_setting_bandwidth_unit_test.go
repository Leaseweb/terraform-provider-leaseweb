package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationSettingBandwidthUnit_Value(t *testing.T) {
	got := NotificationSettingBandwidthUnitMbps.String()

	assert.Equal(t, "Mbps", got)
}

func TestNewNotificationSettingBandwidthUnit(t *testing.T) {
	want := NotificationSettingBandwidthUnitMbps
	got, err := NewNotificationSettingBandwidthUnit("Mbps")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestNotificationSettingBandwidthUnit_Values(t *testing.T) {
	want := []string{"Mbps", "Gbps"}
	got := NotificationSettingBandwidthUnitMbps.Values()

	assert.Equal(t, want, got)
}
