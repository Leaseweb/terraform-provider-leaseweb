package to_resource_model

import (
	"context"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AdaptNotificationSettingBandwidth(t *testing.T) {
	notificationSettingBandwidth := dedicated_server.NotificationSettingBandwidth{
		ServerId: "123345",
	}

	got, _ := AdaptNotificationSettingBandwidth(context.TODO(), notificationSettingBandwidth)
	assert.Equal(t, "123345", got.ServerId.ValueString())
}

func Test_adaptAction(t *testing.T) {
	action := dedicated_server.Action{
		Type: "EMAIL",
	}

	got, _ := adaptAction(context.TODO(), action)
	assert.Equal(t, "EMAIL", got.Type.ValueString())
}
