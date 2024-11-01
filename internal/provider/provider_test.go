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
							"image.name",
							"Ubuntu 20.04 LTS (x86_64)",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"image.custom",
							"false",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"image.flavour",
							"ubuntu",
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
	t.Run("non existing instanceType is no accepted", func(t *testing.T) {
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

	t.Run("instanceType not in region is not accepted", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = "eu-west-3"
  type = "lsw.m5.large"
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
					ExpectError: regexp.MustCompile("Attribute region value must be one of"),
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

	t.Run(
		"instanceId must be valid when creating a custom image",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_image" "test" {
  id = "tralala"
  name = "Custom image"
}`,
						ExpectError: regexp.MustCompile("Attribute id value must be one of"),
					},
				},
			})
		},
	)

	t.Run(
		"instance connected to instanceId must have a `STOPPED` state",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_image" "test" {
  id = "f28ba2af-7508-4594-a63a-aa663db4fb3e"
  name = "Custom image"
}`,
						ExpectError: regexp.MustCompile("not have state"),
					},
				},
			})
		},
	)

	t.Run(
		"instance connected to instanceId must not have a large rootDiskSize",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_image" "test" {
  id = "6871686d-36c4-44f5-b692-a548e62dcf25"
  name = "Custom image"
}`,
						ExpectError: regexp.MustCompile(`rootDiskSize`),
					},
				},
			})
		},
	)

	t.Run(
		"instance connected to instanceId must not windows os",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_public_cloud_image" "test" {
  id = "9c095e3a-e9e3-403b-8d1b-37bb21b5598e"
  name = "Custom image"
}`,
						ExpectError: regexp.MustCompile(`windows`),
					},
				},
			})
		},
	)
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

func TestAccDedicatedServerNotificationSettingBandwidthResource(t *testing.T) {
	t.Run("creates a notification setting bandwidth", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
    dedicated_server_id = "12345678"
    frequency = "WEEKLY"
    threshold = "1"
    unit = "Gbps"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_bandwidth.test",
							"id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_bandwidth.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_bandwidth.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_bandwidth.test",
							"unit",
							"Gbps",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_bandwidth.test",
							"dedicated_server_id",
							"12345678",
						),
					),
				},
			},
		})
	})

	t.Run(
		"server id should be there in the request",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
    frequency = "WEEKLY"
    threshold = "1"
    unit = "Gbps"
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
		"frequency should be one of these values 'DAILY', 'MONTHLY', 'WEEKLY'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
	resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
	   dedicated_server_id = "12345678"
	   frequency = "WRONG"
	   threshold = "1"
	   unit = "Gbps"
	}`,
						ExpectError: regexp.MustCompile(
							`Attribute frequency value must be one of: \["DAILY" "WEEKLY" "MONTHLY"]`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"threshold should be greater than 0",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
    dedicated_server_id = "12345678"
    frequency = "DAILY"
    threshold = "0"
    unit = "Gbps"
}`,
						ExpectError: regexp.MustCompile(
							"The value must be greater than 0, but got 0",
						),
					},
				},
			})
		},
	)

	t.Run(
		"unit should be one of these values 'Mbps', 'Gbps'",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
	resource "leaseweb_dedicated_server_notification_setting_bandwidth" "test" {
	   dedicated_server_id = "12345678"
	   frequency = "DAILY"
	   threshold = "0"
	   unit = "Kbps"
	}`,
						ExpectError: regexp.MustCompile(
							`Attribute unit value must be one of: \["Mbps" "Gbps"], got: "Kbps"`,
						),
					},
				},
			})
		},
	)
}

func TestAccControlPanelsDataSource(t *testing.T) {

	t.Run(
		"getting all control panels",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
	data "leaseweb_dedicated_server_control_panels" "dtest" {
	}
`,
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
		},
	)

	t.Run(
		"filtering control panels by operating_system_id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
	data "leaseweb_dedicated_server_control_panels" "dtest" {
	    operating_system_id = "ALMALINUX_8_64BIT"
	}
`,
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
		},
	)
}

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
							`Attribute type value must be one of: \["OPERATING_SYSTEM" "CONTROL_PANEL"(\s*)"REMOTE_MANAGEMENT" "RESCUE_MODE" "SWITCH" "PDU" "FIREWALL" "LOAD_BALANCER"],(\s*)got: "invalid"`,
						),
					},
				},
			})
		},
	)
}

func TestAccDataTrafficNotificationSettingResource(t *testing.T) {
	t.Run("creates and updates a data traffic notification setting", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
  dedicated_server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "GB"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"dedicated_server_id",
							"145406",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"unit",
							"GB",
						),
					),
				},
				// Update and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
  dedicated_server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "GB"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"id",
							"12345",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"dedicated_server_id",
							"145406",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"frequency",
							"WEEKLY",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"threshold",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_dedicated_server_notification_setting_datatraffic.test",
							"unit",
							"GB",
						),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run(
		"threshold must be greater than 0",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
  dedicated_server_id = "145406"
  frequency = "WEEKLY"
  threshold = "-1"
  unit = "GB"
}`,
						ExpectError: regexp.MustCompile(
							"The value must be greater than 0, but got -1",
						),
					},
				},
			})
		},
	)

	t.Run(
		"unit must be one of GB,MB,TB",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
  dedicated_server_id = "145406"
  frequency = "WEEKLY"
  threshold = "1"
  unit = "blah"
}`,
						ExpectError: regexp.MustCompile(
							`Attribute unit value must be one of: \["MB" "GB" "TB"], got: "blah"`,
						),
					},
				},
			})
		},
	)

	t.Run(
		"frequency must be one of DAILY,WEEKLY,MONTHLY",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
resource "leaseweb_dedicated_server_notification_setting_datatraffic" "test" {
  dedicated_server_id = "145406"
  frequency = "blah"
  threshold = "1"
  unit = "GB"
}`,
						ExpectError: regexp.MustCompile(
							`Attribute frequency value must be one of: \["DAILY" "WEEKLY" "MONTHLY"], got:`,
						),
					},
				},
			})
		},
	)
}

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

func TestAccOperatingSystemsDataSource(t *testing.T) {

	t.Run(
		"getting all operating systems",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
	data "leaseweb_dedicated_server_operating_systems" "dtest" {
	}
`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.#",
								"24",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.0.id",
								"ALMALINUX_8_64BIT",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.0.name",
								"AlmaLinux 8 (x86_64)",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"filtering operating systems by control_panel_id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
	data "leaseweb_dedicated_server_operating_systems" "dtest" {
		control_panel_id = "CPANEL_PREMIER_100"
	}
`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.#",
								"24",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.0.id",
								"ALMALINUX_8_64BIT",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server_operating_systems.dtest",
								"operating_systems.0.name",
								"AlmaLinux 8 (x86_64)",
							),
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServerDataSource(t *testing.T) {
	t.Run(
		"getting dedicated server detail by id",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_server" "test" {
			id = "12345"
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"id",
								"12345",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"asset_id",
								"627294",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"serial_number",
								"JDK18291JK",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"contract_id",
								"674382",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"rack_type",
								"DEDICATED",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_automation_feature_available",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_ipmi_reboot_feature_available",
								"false",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_power_cycle_feature_available",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_private_network_feature_available",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"is_remote_management_feature_available",
								"false",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_rack",
								"13",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_site",
								"AMS-01",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_suite",
								"A6",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"location_unit",
								"16-17",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"public_mac",
								"AA:BB:CC:DD:EE:FF",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"public_ip",
								"123.123.123.123/27",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"public_gateway",
								"123.123.123.126",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"internal_mac",
								"AA:BB:CC:DD:EE:FF",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"internal_ip",
								"123.123.123.123/27",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"internal_gateway",
								"123.123.123.126",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"ram_size",
								"32",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"ram_unit",
								"GB",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"cpu_quantity",
								"4",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_server.test",
								"cpu_type",
								"Intel Xeon E3-1220",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"id is required for getting the dedicated server detail",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_server" "test" {
		}`,
						ExpectError: regexp.MustCompile(
							`The argument "id" is required, but no definition was found`,
						),
					},
				},
			})
		},
	)
}

func TestAccDedicatedServersDataSource(t *testing.T) {
	t.Run(
		"getting dedicated servers by reference",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_servers" "test" {
			reference = "test-reference"
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"ids.#",
								"2",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"reference",
								"test-reference",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"getting dedicated servers",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_servers" "test" {
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.test",
								"ids.#",
								"2",
							),
						),
					},
				},
			})
		},
	)

	t.Run(
		"getting dedicated servers with all filters",
		func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: providerConfig + `
		data "leaseweb_dedicated_servers" "filter" {
			reference = "test-reference"
			ip = "127.0.0.4"
			mac_address = "aa:bb:cc:dd:ee:ff"
			site = "ams-01"
			private_rack_id = "r id"
			private_network_capable = "true"
			private_network_enabled = "true"
		}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"ids.#",
								"2",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"reference",
								"test-reference",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"ip",
								"127.0.0.4",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"mac_address",
								"aa:bb:cc:dd:ee:ff",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"site",
								"ams-01",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"private_rack_id",
								"r id",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"private_network_capable",
								"true",
							),
							resource.TestCheckResourceAttr(
								"data.leaseweb_dedicated_servers.filter",
								"private_network_enabled",
								"true",
							),
						),
					},
				},
			})
		},
	)

}

func TestAccLoadBalancerResource(t *testing.T) {
	t.Run("creates and updates a loadBalancer", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my-loadbalancer1"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"id",
							"32082a93-d1e2-4bc0-8f5e-8fe4312b0844",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"region",
							"eu-west-3",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"type",
							"lsw.m3.large",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"reference",
							"my-loadbalancer1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"contract.billing_frequency",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"contract.term",
							"0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"contract.type",
							"HOURLY",
						),
					),
				},
				// ImportState testing
				{
					ResourceName:      "leaseweb_public_cloud_load_balancer.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my-loadbalancer1"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"id",
							"32082a93-d1e2-4bc0-8f5e-8fe4312b0844",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"region",
							"eu-west-3",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"type",
							"lsw.m3.large",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"reference",
							"my-loadbalancer1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"contract.billing_frequency",
							"1",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
							"contract.term",
							"0",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_load_balancer.test",
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
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my-loadbalancer1"
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
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my-loadbalancer1"
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
	t.Run("invalid type", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "tralala"
  reference = "my-loadbalancer1"
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

	t.Run("invalid region", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "tralala"
  type = "lsw.m4.2xlarge"
  reference = "my-loadbalancer1"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}`,
					ExpectError: regexp.MustCompile("Attribute region value must be one of"),
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
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.2xlarge"
  reference = "my-loadbalancer1"
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
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.2xlarge"
  reference = "my-loadbalancer1"
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
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.2xlarge"
  reference = "my-loadbalancer1"
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
			requiredField: "type",
			expectedError: fmt.Sprintf(
				"The argument %q is required, but no definition was found.",
				"type",
			),
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
resource "leaseweb_public_cloud_load_balancer" "test" {
  contract = {}
}`,
						ExpectError: regexp.MustCompile(scenario.expectedError),
					},
				},
			})
		})
	}

	t.Run("changing the region triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-3"
  type = "lsw.m3.large"
  reference = "my-loadbalancer1"
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
								"leaseweb_public_cloud_load_balancer.test",
								plancheck.ResourceActionDestroyBeforeCreate,
							),
						},
					},
					// Ignore the inconsistent result as prism returns the old result.
					ExpectError: regexp.MustCompile(
						"Provider produced inconsistent result after apply",
					),
					Config: providerConfig + `
resource "leaseweb_public_cloud_load_balancer" "test" {
  region = "eu-west-2"
  type = "lsw.m3.large"
  reference = "my-loadbalancer1"
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
}
