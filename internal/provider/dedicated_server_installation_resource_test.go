package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerInstallationResource(t *testing.T) {
	t.Run("install os on a dedicated server",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Create testing
					{
						Config: providerConfig + `
    resource "leaseweb_dedicated_server_installation" "test" {
      dedicated_server_id = "12345"
      operating_system_id = "UBUNTU_22_04_64BIT"
    }`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"leaseweb_dedicated_server_installation.test",
								"id",
								"bcf2bedf-8450-4b22-86a8-f30aeb3a38f9",
							),
							resource.TestCheckResourceAttr(
								"leaseweb_dedicated_server_installation.test",
								"dedicated_server_id",
								"12345",
							),
							resource.TestCheckResourceAttr(
								"leaseweb_dedicated_server_installation.test",
								"operating_system_id",
								"UBUNTU_22_04_64BIT",
							),
						),
					},
				},
			})
		})

	t.Run(
		"server id should be in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
    resource "leaseweb_dedicated_server_installation" "test" {
      operating_system_id = "UBUNTU_22_04_64BIT"
    }`,
						ExpectError: regexp.MustCompile(
							"The argument \"dedicated_server_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"operating system id should be in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
    resource "leaseweb_dedicated_server_installation" "test" {
      dedicated_server_id = "12345"
    }`,
						ExpectError: regexp.MustCompile(
							"The argument \"operating_system_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"raid.level should be one of these values '0', '1', '5', '10'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
		resource "leaseweb_dedicated_server_installation" "test" {
      dedicated_server_id = "12345"
      operating_system_id = "UBUNTU_22_04_64BIT"
      raid = {
        level = 11
      }
    }`,
						ExpectError: regexp.MustCompile(
							`Attribute raid.level value must be one of: \["0" "1" "5" "10"]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"raid.type should be one of these values 'HW', 'SW', 'NONE'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
		resource "leaseweb_dedicated_server_installation" "test" {
      dedicated_server_id = "12345"
      operating_system_id = "UBUNTU_22_04_64BIT"
      raid = {
        type = "TEST"
      }
    }`,
						ExpectError: regexp.MustCompile(
							`Attribute raid.type value must be one of: \["HW" "SW" "NONE"]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"ssh_keys should be set of string",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
		resource "leaseweb_dedicated_server_installation" "test" {
      dedicated_server_id = "12345"
      operating_system_id = "UBUNTU_22_04_64BIT"
      ssh_keys = "test keys"
    }`,
						ExpectError: regexp.MustCompile(
							`Inappropriate value for attribute "ssh_keys": set of string required`,
						),
					},
				},
			})
		},
	)
}
