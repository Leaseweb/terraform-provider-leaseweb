package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerDataSource(t *testing.T) {

	t.Run(
		"getting dedicated server detail by id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_server" "test" {
			id = "12345"
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"id",
								"12345",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"asset_id",
								"627294",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"serial_number",
								"JDK18291JK",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"contract_id",
								"674382",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"rack_type",
								"DEDICATED",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_automation_feature_available",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_ipmi_reboot_feature_available",
								"false",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_power_cycle_feature_available",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_private_network_feature_available",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_remote_management_feature_available",
								"false",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_rack",
								"13",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_site",
								"AMS-01",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_suite",
								"A6",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_unit",
								"16-17",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"public_mac",
								"AA:BB:CC:DD:EE:FF",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"public_ip",
								"123.123.123.123/27",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"public_gateway",
								"123.123.123.126",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"internal_mac",
								"AA:BB:CC:DD:EE:FF",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"internal_ip",
								"123.123.123.123/27",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"internal_gateway",
								"123.123.123.126",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"ram_size",
								"32",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"ram_unit",
								"GB",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"cpu_quantity",
								"4",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"cpu_type",
								"Intel Xeon E3-1220",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"id is required for getting the dedicated server detail",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_server" "test" {
		}`,
						ExpectError: regexp.MustCompile(
							`The argument "id" is required, but no definition was found`,
						),
					},
				},
			})
		},
	)

}
