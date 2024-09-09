package notification_setting_bandwidth

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func Test_notificationSettingBandwidthResource_Metadata(t *testing.T) {
	resp := resource.MetadataResponse{}
	notificationSettingBandwidthResource := New()

	notificationSettingBandwidthResource.Metadata(
		context.TODO(),
		resource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(t,
		"tralala_dedicated_server_notification_setting_bandwidth",
		resp.TypeName,
		"Type name should be tralala_dedicated_server_notification_setting_bandwidth",
	)
}
