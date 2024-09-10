package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataTrafficNotificationSettingResource(t *testing.T) {
	t.Run("creates and updates a data traffic notification setting", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "GB"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"server_id",
							"145406",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"unit",
							"GB",
						),
					),
				},
				// Update and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "GB"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"server_id",
							"145406",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_data_traffic_notification_setting.test",
							"unit",
							"GB",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}
