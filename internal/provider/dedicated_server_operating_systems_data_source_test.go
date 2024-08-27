package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerOperatingSystemsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_dedicated_server_operating_systems" "dtest" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_operating_systems.dtest",
						"operating_systems.#",
						"24",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_operating_systems.dtest",
						"operating_systems.0.id",
						"ALMALINUX_8_64BIT",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_operating_systems.dtest",
						"operating_systems.0.name",
						"AlmaLinux 8 (x86_64)",
					),
				),
			},
		},
	})
}
