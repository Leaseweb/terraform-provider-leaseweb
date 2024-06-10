package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccInstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "leaseweb_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.large"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("leaseweb_instance.test", "id", "ace712e9-a166-47f1-9065-4af0f7e7fce1"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "region", "eu-west-3"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "type", "lsw.m3.large"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "reference", "my webserver"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "operating_system.id", "UBUNTU_22_04_64BIT"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "root_disk_storage_type", "CENTRAL"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "contract.billing_frequency", "1"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "contract.term", "0"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "contract.type", "HOURLY"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "leaseweb_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "leaseweb_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.large"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
			  }
			  `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("leaseweb_instance.test", "id", "ace712e9-a166-47f1-9065-4af0f7e7fce1"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "region", "eu-west-3"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "type", "lsw.m3.large"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "reference", "my webserver"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "operating_system.id", "UBUNTU_22_04_64BIT"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "root_disk_storage_type", "CENTRAL"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "contract.billing_frequency", "1"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "contract.term", "0"),
					resource.TestCheckResourceAttr("leaseweb_instance.test", "contract.type", "HOURLY"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
