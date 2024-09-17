package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerCredentialResource(t *testing.T) {
	t.Run("creates and updates a credential", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_credential" "test" {
	dedicated_server_id = "12345"
   	username = "root"
   	type = "OPERATING_SYSTEM"
   	password = "mys3cr3tp@ssw0rd"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"dedicated_server_id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"username",
							"root",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"type",
							"OPERATING_SYSTEM",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"password",
							"mys3cr3tp@ssw0rd",
						),
					),
				},
				// Update and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_credential" "test" {
	dedicated_server_id = "12345"
   	username = "root"
   	type = "OPERATING_SYSTEM"
   	password = "mys3cr3tp@ssw0rd"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"dedicated_server_id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"username",
							"root",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"type",
							"OPERATING_SYSTEM",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_credential.test",
							"password",
							"mys3cr3tp@ssw0rd",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run(
		"type must be a valid one",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_credential" "test" {
	dedicated_server_id = "12345"
   	username = "root"
   	type = "invalid"
   	password = "mys3cr3tp@ssw0rd"
}`,
						ExpectError: regexp.MustCompile(
							`Attribute type value must be one of: \["OPERATING_SYSTEM" "CONTROL_PANEL"(\s*)"REMOTE_MANAGEMENT" "RESCUE_MODE" "SWITCH" "PDU" "FIREWALL" "LOAD_BALANCER"\],(\s*)got: "invalid"`,
						),
					},
				},
			})
		},
	)
}
