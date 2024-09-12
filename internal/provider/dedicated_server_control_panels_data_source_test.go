package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccControlPanelsDataSource(t *testing.T) {

	t.Run(
		"getting all control panels",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
	data "leaseweb_dedicated_server_control_panels" "dtest" {
	}
`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_control_panels.dtest",
								"control_panels.#",
								"8",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_control_panels.dtest",
								"control_panels.0.id",
								"CPANEL_PREMIER_100",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"filterring control panels by operating_system_id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
	data "leaseweb_dedicated_server_control_panels" "dtest" {
	    operating_system_id = "ALMALINUX_8_64BIT"
	}
`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_control_panels.dtest",
								"control_panels.#",
								"8",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_control_panels.dtest",
								"control_panels.0.id",
								"CPANEL_PREMIER_100",
							),
						),
					},
				},
			})
		},
	)
}
