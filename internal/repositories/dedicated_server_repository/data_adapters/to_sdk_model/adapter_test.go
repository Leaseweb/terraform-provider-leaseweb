package to_sdk_model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/stretchr/testify/assert"
)

func TestAdaptToCreateNotificationSettingBandwidthOpts(t *testing.T) {
	domainNotificationSettingBandwidth := dedicated_server.NotificationSettingBandwidth{
		Frequency: "DAILY",
		Unit:      "Mbps",
		Threshold: "1",
	}
	got := AdaptToCreateNotificationSettingBandwidthOpts(domainNotificationSettingBandwidth)

	want := dedicatedServer.NewBandwidthNotificationSettingOpts("DAILY", "1", "Mbps")
	assert.Equal(t, want, got)
}
