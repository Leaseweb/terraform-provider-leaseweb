package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataTrafficNotificationSettingResource(t *testing.T) {
	t.Run("creates and updates a data traffic notification setting", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "GB"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"server_id",
							"145406",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"unit",
							"GB",
						),
					),
				},
				// Update and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "GB"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"server_id",
							"145406",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_data_traffic_notification_setting.test",
							"unit",
							"GB",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run(
		"threshold must be greater than 0",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "WEEKLY"
  threshold = "-1"
  unit = "GB"
}`,
						ExpectError: regexp.MustCompile(
							"The value must be greater than 0, but got -1",
						),
					},
				},
			})
		},
	)

	t.Run(
		"unit must be one of GB,MB,TB",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "blah"
}`,
						ExpectError: regexp.MustCompile(
							`Attribute unit value must be one of: \["MB" "GB" "TB"\], got: "blah"`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"frequency must be one of DAILY,WEEKLY,MONTHLY",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_data_traffic_notification_setting" "test" {
  server_id = "145406"
  frequency = "blah"
  threshold = "1"
  unit = "GB"
}`,
						ExpectError: regexp.MustCompile(
							`Attribute frequency value must be one of: \["DAILY" "WEEKLY" "MONTHLY"\], got:`,
						),
					},
				},
			})
		},
	)
}
