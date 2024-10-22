package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
		"leaseweb": providerserver.NewProtocol6WithError(NewProvider("test")()),
	}
)

func TestLeasewebProvider_Metadata(t *testing.T) {
	leasewebProvider := NewProvider("dev")
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
	leasewebProvider := NewProvider("dev")
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
						"1",
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

func TestAccInstanceResource(t *testing.T) {
	t.Run("creates and updates an instance", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
							"image.id",
							"UBUNTU_20_04_64BIT",
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
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
							"image.id",
							"UBUNTU_20_04_64BIT",
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

	t.Run(
		"term must be 0 when contract type is HOURLY",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 3
    type              = "HOURLY"
  }
}`,
						ExpectError: regexp.MustCompile(
							"Attribute contract.term must be 0 when contract.type is \"HOURLY\", got: 3",
						),
					},
				},
			})
		},
	)

	t.Run("term must not be 0 when contract type is MONTHLY", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "MONTHLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.term cannot be 0 when contract.type is \"MONTHLY\", got: 0",
					),
				},
			},
		})
	})
	t.Run("invalid instanceType", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "tralala"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute type value must be one of:",
					),
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
  region = "eu-west-3"
  type = "lsw.m4.4xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  root_disk_size = 1
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute root_disk_size value must be between",
					),
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
  region = "eu-west-3"
  type = "lsw.m4.4xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  root_disk_size = 1001
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute root_disk_size value must be between",
					),
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
  region = "eu-west-3"
  type = "lsw.m4.4xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "tralala"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute root_disk_storage_type value must be one of",
					),
				},
			},
		})
	})

	t.Run("invalid region", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "tralala"
  type = "lsw.m4.2xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile("Invalid Region"),
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
  region = "eu-west-3"
  type = "lsw.m3.2xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 55
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.billing_frequency value must be one of",
					),
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
  region = "eu-west-3"
  type = "lsw.m3.2xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 55
    type              = "MONTHLY"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.term value must be one of",
					),
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
  region = "eu-west-3"
  type = "lsw.m3.2xlarge"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 3
    type              = "tralala"
  }
}`,
					ExpectError: regexp.MustCompile(
						"Attribute contract.type value must be one of",
					),
				},
			},
		})
	})

	type errorTestCases struct {
		requiredField string
		expectedError string
	}

	for _, scenario := range []errorTestCases{
		{
			requiredField: "region",
			expectedError: fmt.Sprintf(
				"The argument %q is required, but no definition was found.",
				"region",
			),
		},
		{
			requiredField: "root_disk_storage_type",
			expectedError: fmt.Sprintf(
				"The argument %q is required, but no definition was",
				"root_disk_storage_type",
			),
		},
		{
			requiredField: "type",
			expectedError: fmt.Sprintf(
				"The argument %q is required, but no definition was found.",
				"type",
			),
		},
		{
			requiredField: "image.id",
			expectedError: "Inappropriate value for attribute \"image\": attribute \"id\"",
		},
		{
			requiredField: "contract.type|contract.term|contract.billing_frequency",
			expectedError: "Inappropriate value for attribute \"contract\": attributes \"billing_frequency\",\n\"term\", and \"type\" are required.",
		},
	} {
		t.Run(scenario.requiredField+" should be set", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  image = {}
  contract = {}
}`,
						ExpectError: regexp.MustCompile(scenario.expectedError),
					},
				},
			})
		})
	}

	t.Run(
		"upgrading to invalid instanceType is not allowed",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
  region = "eu-west-3"
  type = "lsw.m4.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
						ExpectError: regexp.MustCompile(
							"Attribute type value must be one of:",
						),
					},
				},
			})
		},
	)

	t.Run("changing the region triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_instance.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-2"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
				},
			},
		})
	})

	t.Run("changing the imageId triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"leaseweb_public_cloud_instance.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
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
			},
		})
	})

	t.Run(
		"changing the marketAppId triggers replacement",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(
									"leaseweb_public_cloud_instance.test",
									plancheck.ResourceActionDestroyBeforeCreate,
								),
							},
						},
						// Ignore the inconsistent result as prism returns the old result.
						ExpectError: regexp.MustCompile(
							"Provider produced inconsistent result after apply",
						),
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  market_app_id = "newValue"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					},
				},
			})
		},
	)

	t.Run(
		"changing the rootDiskStorageType triggers replacement",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
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
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(
									"leaseweb_public_cloud_instance.test",
									plancheck.ResourceActionDestroyBeforeCreate,
								),
							},
						},
						// Ignore the inconsistent result as prism returns the old result.
						ExpectError: regexp.MustCompile(
							"Provider produced inconsistent result after apply",
						),
						Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my webserver"
  image = {
    id = "UBUNTU_20_04_64BIT"
  }
  root_disk_storage_type = "LOCAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					},
				},
			})
		},
	)
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
							`Attribute type value must be one of: \["OPERATING_SYSTEM" "CONTROL_PANEL"\],(\s*)got: "invalid"`,
						),
					},
				},
			})
		},
	)
}
