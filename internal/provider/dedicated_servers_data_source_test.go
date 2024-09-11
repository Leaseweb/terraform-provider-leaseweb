package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServersDataSource(t *testing.T) {

	t.Run(
		"getting dedicated servers by reference",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_servers" "test" {
			reference = "test-reference"
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"ids.#",
								"2",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"reference",
								"test-reference",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"getting dedicated servers",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_servers" "test" {
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"ids.#",
								"2",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"getting dedicated servers with all filters",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_servers" "filter" {
			reference = "test-reference"
			ip = "127.0.0.4"
			mac_address = "aa:bb:cc:dd:ee:ff"
			site = "ams-01"
			private_rack_id = "r id"
			private_network_capable = "true"
			private_network_enabled = "true"
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"ids.#",
								"2",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"reference",
								"test-reference",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"ip",
								"127.0.0.4",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"mac_address",
								"aa:bb:cc:dd:ee:ff",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"site",
								"ams-01",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"private_rack_id",
								"r id",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"private_network_capable",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"private_network_enabled",
								"true",
							),
						),
					},
				},
			})
		},
	)

}
