package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
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
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"id",
						"ace712e9-a166-47f1-9065-4af0f7e7fce1",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"region",
						"eu-west-3",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"type",
						"lsw.m3.large",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"reference",
						"my webserver",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"operating_system.id",
						"UBUNTU_22_04_64BIT",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"root_disk_storage_type",
						"CENTRAL",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"contract.billing_frequency",
						"1",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"contract.term",
						"0",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"contract.type",
						"HOURLY",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "leaseweb_public_cloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
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
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"id",
						"ace712e9-a166-47f1-9065-4af0f7e7fce1",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"region",
						"eu-west-3",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"type",
						"lsw.m3.large",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"reference",
						"my webserver",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"operating_system.id",
						"UBUNTU_22_04_64BIT",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"root_disk_storage_type",
						"CENTRAL",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"contract.billing_frequency",
						"1",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"contract.term",
						"0",
					),
					resource.TestCheckResourceAttr(
						"leaseweb_public_cloud_instance.test",
						"contract.type",
						"HOURLY",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccInstanceResource_validationError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// term must be 0 when contract type is HOURLY
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.large"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 5
    type              = "HOURLY"
  }
}`,
				ExpectError: regexp.MustCompile("Attribute contract.term must be 0 when contract.type is \"HOURLY\", got: 5"),
			},
			// term must not be 0 when contract type is MONTHLY
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
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
    type              = "MONTHLY"
  }
}`,
				ExpectError: regexp.MustCompile("Attribute contract.term cannot be 0 when contract.type is \"MONTHLY\", got: 0"),
			},
			// Invalid instance type
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "tralala"
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
}`,
				ExpectError: regexp.MustCompile("Attribute type value must be one of:"),
			},
			// invalid ssh key
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.4xlarge"
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
  ssh_key = "tralala"
}`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			// root_disk_size is too small
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.4xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  root_disk_size = 1
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
				ExpectError: regexp.MustCompile("Attribute root_disk_size value must be between"),
			},
			// root_disk_size is too big
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.4xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  root_disk_size = 1001
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
				ExpectError: regexp.MustCompile("Attribute root_disk_size value must be between"),
			},
			// Invalid root_disk_storage_type
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.4xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "tralala"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
				ExpectError: regexp.MustCompile("Attribute root_disk_storage_type value must be one of"),
			},
			// Invalid contract.billing_frequency
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.4xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRY"
  contract = {
    billing_frequency = 55
    term              = 0
    type              = "HOURLY"
  }
}`,
				ExpectError: regexp.MustCompile("Attribute root_disk_storage_type value must be one of"),
			},
		},
	})
}

func TestAccInstanceResource_Choosing_Invalid_Instance_Type_Not_Allowed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
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
}`,
			},
			{
				Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.large"
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
}`,
				ExpectError: regexp.MustCompile("Invalid Instance Type"),
			},
		},
	})
}
