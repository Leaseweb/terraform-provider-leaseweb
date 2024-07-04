package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInstanceResource(t *testing.T) {
	t.Run("creates and updates an instance", func(t *testing.T) {
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
	})

	t.Run("term must be 0 when contract type is HOURLY", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
			},
		})
	})

	t.Run("term must not be 0 when contract type is MONTHLY", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
			},
		})
	})
	t.Run("invalid instance type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
			},
		})
	})

	t.Run("invalid ssh key", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
			},
		})
	})

	t.Run("rootDiskSize is too small", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
			},
		})
	})

	t.Run("rootDiskSize is too big", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
			},
		})
	})

	t.Run("invalid rootDiskStorageType", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m4.2xlarge"
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
			},
		})
	})

	t.Run("invalid contract.billingFrequency", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.2xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 55
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile("Attribute contract.billing_frequency value must be one of"),
				},
			},
		})
	})

	t.Run("invalid contract.term", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.2xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 55
    type              = "MONTHLY"
  }
}`,
					ExpectError: regexp.MustCompile("Attribute contract.term value must be one of"),
				},
			},
		})
	})

	t.Run("invalid contract.type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.2xlarge"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 3
    type              = "tralala"
  }
}`,
					ExpectError: regexp.MustCompile("Attribute contract.type value must be one of"),
				},
			},
		})
	})

	t.Run("invalid operating_system_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region    = "eu-west-3"
  type      = "lsw.m3.large"
  reference = "my webserver"
  operating_system = {
    id = "tralala"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile("Attribute operating_system.id value must be one of"),
				},
			},
		})
	})

	t.Run("upgrading to invalid instanceType is not allowed", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
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
	})
}
