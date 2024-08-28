package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccControlPanelsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_dedicated_server_control_panels" "dtest" {}`,
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
}
