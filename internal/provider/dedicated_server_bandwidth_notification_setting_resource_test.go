package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerNotificationSettingBandwidthResource(t *testing.T) {
	t.Run("creates a notification setting bandwidth", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_bandwidth_notification_setting" "test" {
    dedicated_server_id = "12345678"
    frequency = "WEEKLY"
    threshold = "1"
    unit = "Gbps"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_bandwidth_notification_setting.test",
							"id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_bandwidth_notification_setting.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_bandwidth_notification_setting.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_bandwidth_notification_setting.test",
							"unit",
							"Gbps",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_bandwidth_notification_setting.test",
							"dedicated_server_id",
							"12345678",
						),
					),
				},
			},
		})
	})

	t.Run(
		"server id should be there in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_bandwidth_notification_setting" "test" {
    frequency = "WEEKLY"
    threshold = "1"
    unit = "Gbps"
}`,
						ExpectError: regexp.MustCompile(
							"The argument \"dedicated_server_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"frequency should be one of these values 'DAILY', 'MONTHLY', 'WEEKLY'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
	resource "leaseweb_dedicated_server_bandwidth_notification_setting" "test" {
	   dedicated_server_id = "12345678"
	   frequency = "WRONG"
	   threshold = "1"
	   unit = "Gbps"
	}`,
						ExpectError: regexp.MustCompile(
							`Attribute frequency value must be one of: \["DAILY" "WEEKLY" "MONTHLY"\]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"threshold should be greater than 0",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_bandwidth_notification_setting" "test" {
    dedicated_server_id = "12345678"
    frequency = "DAILY"
    threshold = "0"
    unit = "Gbps"
}`,
						ExpectError: regexp.MustCompile(
							"The value must be greater than 0, but got 0",
						),
					},
				},
			})
		},
	)

	t.Run(
		"unit should be one of these values 'Mbps', 'Gbps'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
	resource "leaseweb_dedicated_server_bandwidth_notification_setting" "test" {
	   dedicated_server_id = "12345678"
	   frequency = "DAILY"
	   threshold = "0"
	   unit = "Kbps"
	}`,
						ExpectError: regexp.MustCompile(
							`Attribute unit value must be one of: \["Mbps" "Gbps"\], got: "Kbps"`,
						),
					},
				},
			})
		},
	)
}
