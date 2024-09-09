package to_domain_entity

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdaptToCreateNotificationSettingBandwidthOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		notificationSettingBandwidth := model.NotificationSettingBandwidth{
			ServerId:  basetypes.NewStringValue("123456"),
			Frequency: basetypes.NewStringValue("DAILY"),
			Threshold: basetypes.NewStringValue("1"),
			Unit:      basetypes.NewStringValue("Gbps"),
		}

		got := AdaptToCreateNotificationSettingBandwidthOpts(notificationSettingBandwidth)

		assert.Equal(t, "123456", got.ServerId)
		assert.Equal(t, "DAILY", got.Frequency)
		assert.Equal(t, "1", got.Threshold)
		assert.Equal(t, "Gbps", got.Unit)
		assert.Nil(t, got.ThresholdExceededAt)
		assert.Nil(t, got.LastCheckedAt)
		assert.Nil(t, got.Actions)
	})
}
