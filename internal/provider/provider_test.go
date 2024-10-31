package provider

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

const (
	providerConfig = `
provider "leaseweb" {
  host     = "localhost:8080"
  scheme = "http"
  token = "tralala"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (
		tfprotov6.ProviderServer,
		error,
	){
		"leaseweb": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestLeasewebProvider_Metadata(t *testing.T) {
	leasewebProvider := New("dev")
	metadataResponse := provider.MetadataResponse{}
	leasewebProvider().Metadata(
		context.TODO(),
		provider.MetadataRequest{},
		&metadataResponse,
	)

	want := "dev"
	got := metadataResponse.Version

	assert.Equal(
		t,
		want,
		got,
		"version should be passed to provider",
	)
}

func TestLeasewebProvider_Schema(t *testing.T) {
	leasewebProvider := New("dev")
	schemaResponse := provider.SchemaResponse{}
	leasewebProvider().Schema(
		context.TODO(),
		provider.SchemaRequest{},
		&schemaResponse,
	)

	assert.True(
		t,
		schemaResponse.Schema.Attributes["host"].IsOptional(),
		"host is optional",
	)
	assert.True(
		t,
		schemaResponse.Schema.Attributes["scheme"].IsOptional(),
		"scheme is optional",
	)
	assert.True(
		t,
		schemaResponse.Schema.Attributes["token"].IsSensitive(),
		"token is sensitive",
	)
}

func TestAccInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_public_cloud_instances" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.#",
						"4",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_instances.test",
						"instances.0.id",
						"ace712e9-a166-47f1-9065-4af0f7e7fce1",
					),
				),
			},
		},
	})
}

func TestAccPublicCloudCredentialResource(t *testing.T) {
	t.Run("creates and updates a credential", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
	resource "leaseweb_public_cloud_credential" "test" {
		instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
	   	username = "root"
	   	type = "OPERATING_SYSTEM"
	   	password = "12341234"
	}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"instance_id",
							"695ddd91-051f-4dd6-9120-938a927a47d0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"username",
							"root",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"type",
							"OPERATING_SYSTEM",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"password",
							"12341234",
						),
					),
				},
				// Update and Read testing
				{
					Config: providerConfig + `
				resource "leaseweb_public_cloud_credential" "test" {
					instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
				   	username = "root"
				   	type = "OPERATING_SYSTEM"
				   	password = "12341234"
				}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"instance_id",
							"695ddd91-051f-4dd6-9120-938a927a47d0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"username",
							"root",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"type",
							"OPERATING_SYSTEM",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_credential.test",
							"password",
							"12341234",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run(
		"username should not be empty",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `

		resource "leaseweb_public_cloud_credential" "test" {
			instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
		   	username = ""
		   	type = "OPERATING_SYSTEM"
		   	password = "blah"
		}`,
						ExpectError: regexp.MustCompile(
							`Attribute username string length must be at least 1, got: 0`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"password should not be empty",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `

		resource "leaseweb_public_cloud_credential" "test" {
			instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
		   	username = "root"
		   	type = "OPERATING_SYSTEM"
		   	password = ""
		}`,
						ExpectError: regexp.MustCompile(
							`Attribute password string length must be at least 1, got: 0`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"type must be a valid one",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `

		resource "leaseweb_public_cloud_credential" "test" {
			instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
		   	username = "root"
		   	type = "invalid"
		   	password = "12341234"
		}`,

						ExpectError: regexp.MustCompile(
							`Attribute type value must be one of:`,
						),
					},
				},
			})
		},
	)
}

func TestAccPublicCloudCredentialDataSource(t *testing.T) {
	t.Run("reading data for public cloud credential",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
        data "leaseweb_public_cloud_credential" "test" {
          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
          type                = "OPERATING_SYSTEM"
          username            = "root"
        }`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_public_cloud_credential.test",
								"instance_id",
								"695ddd91-051f-4dd6-9120-938a927a47d0",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_public_cloud_credential.test",
								"type",
								"OPERATING_SYSTEM",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_public_cloud_credential.test",
								"username",
								"root",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_public_cloud_credential.test",
								"password",
								"12341234",
							),
						),
					},
				},
			})
		})

	t.Run(
		"instance_id is required",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
        data "leaseweb_public_cloud_credential" "test" {
          type                = "OPERATING_SYSTEM"
          username            = "root"
        }`,
						ExpectError: regexp.MustCompile(
							"The argument \"instance_id\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"type is required",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
        data "leaseweb_public_cloud_credential" "test" {
          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
          username            = "root"
        }`,
						ExpectError: regexp.MustCompile(
							"The argument \"type\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"username is required",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
        data "leaseweb_public_cloud_credential" "test" {
          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
          type                = "OPERATING_SYSTEM"
        }`,
						ExpectError: regexp.MustCompile(
							"The argument \"username\" is required, but no definition was found",
						),
					},
				},
			})
		},
	)

	t.Run(
		"invalid type is not accepted",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
        data "leaseweb_public_cloud_credential" "test" {
          instance_id         = "695ddd91-051f-4dd6-9120-938a927a47d0"
          type                = "A_WRONG_TYPE"
          username            = "root"
        }`,
						ExpectError: regexp.MustCompile(
							"Attribute type value must be one of:",
						),
					},
				},
			})
		},
	)
}

func TestAccImagesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_public_cloud_images" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.#",
						"6",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_images.test",
						"images.0.id",
						"UBUNTU_24_04_64BIT",
					),
				),
			},
		},
	})
}

func TestAccInstanceImage(t *testing.T) {
	t.Run("creates & updates an image", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_image" "test" {
  id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
  name = "Custom image - 03"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"id",
							"ace712e9-a166-47f1-9065-4af0f7e7fce1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"name",
							"Custom image - 03",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:      "leaseweb_public_cloud_image.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_image" "test" {
  id = "ace712e9-a166-47f1-9065-4af0f7e7fce1"
  name = "Custom image - 03"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"id",
							"ace712e9-a166-47f1-9065-4af0f7e7fce1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_image.test",
							"name",
							"Custom image - 03",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})
}

func TestAccLoadBalancersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leaseweb_public_cloud_load_balancers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.leaseweb_public_cloud_load_balancers.test",
						"load_balancers.0.id",
						"5fd135a9-3ff6-4794-8b92-8cd8747a3ea3",
					),
				),
			},
		},
	})
}
