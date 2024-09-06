package data_traffic_notification_setting

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func Test_dataTrafficNotificationSettingResource_Metadata(t *testing.T) {
	resp := resource.MetadataResponse{}
	dataTrafficNotificationSettingResource := NewDataTrafficNotificationSettingResource()

	dataTrafficNotificationSettingResource.Metadata(
		context.TODO(),
		resource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(t,
		"tralala_dedicated_server_data_traffic_notification_setting",
		resp.TypeName,
		"Type name should be tralala_dedicated_server_data_traffic_notification_setting",
	)
}
