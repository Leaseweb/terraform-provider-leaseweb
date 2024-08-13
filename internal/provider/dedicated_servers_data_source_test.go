package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_dedicated_servers" "dtest" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_servers.dtest",
						"dedicated_servers.#",
						"6",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_servers.dtest",
						"dedicated_servers.0.id",
						"12345",
					),
				),
			},
		},
	})
}
