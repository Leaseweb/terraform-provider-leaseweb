package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerCredentialDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
        data "leaseweb_dedicated_server_credential" "test" {
          dedicated_server_id = "12345"
          type                = "OPERATING_SYSTEM"
          username            = "root"
        }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_credential.test",
						"dedicated_server_id",
						"12345",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_credential.test",
						"type",
						"OPERATING_SYSTEM",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_credential.test",
						"username",
						"root",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_dedicated_server_credential.test",
						"password",
						"mys3cr3tp@ssw0rd",
					),
				),
			},
		},
	})
}
