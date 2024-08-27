package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
							"region.name",
							"eu-west-3",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"type.name",
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
							"region.name",
							"eu-west-3",
						),
						resource.TestCheckResourceAttr(
							"leaseweb_public_cloud_instance.test",
							"type.name",
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "tralala"
  }
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

	// TODO Enable SSH key support
	/**
		  	t.Run("invalid sshKey", func(t *testing.T) {
		  		resource.Test(t, resource.TestCase{
		  			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		  			Steps: []resource.TestStep{
		  				{
		  					Config: providerConfig + `
		  resource "leaseweb_public_cloud_instance" "test" {
		    region = {
	        name = "eu-west-3"
	      }
		    type = {
		      name = "lsw.m4.4xlarge"
		    }
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
		    ssh_key = "tralala"
		  }`,
		  					ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
		  				},
		  			},
		  		})
		  	})
	*/

	t.Run("rootDiskSize is too small", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m4.4xlarge"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m4.4xlarge"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m4.2xlarge"
  }
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
  region = {
    name = "tralala"
  }
  type = {
    name = "lsw.m4.2xlarge"
  }
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.2xlarge"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m3.2xlarge"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m3.2xlarge"
  }
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
			requiredField: "type.name",
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m4.large"
  }
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

	// TODO Enable SSH key support
	/**
		  	t.Run(
		  		"changing the sshKey is not allowed",
		  		func(t *testing.T) {
		  			resource.Test(t, resource.TestCase{
		  				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		  				Steps: []resource.TestStep{
		  					{
		  						Config: providerConfig + `
		  resource "leaseweb_public_cloud_instance" "test" {
	      region = {
	        name = "eu-west-3"
	      }
		    type = {
		      name = "lsw.m3.large"
		    }
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
	      region = {
	          name = "eu-west-3"
	        }
		    type = {
		      name = "lsw.m3.large"
		    }
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
		    ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCbRsxME8r5CbjnXcPj2IydrrDlqDjqvvK4vd4a6zDyP+Pu47HuBdbIqskdDviS/6ZHuMm7x9On/4VDRFaqVUSDAHqkJktBGgsrpoLxy5OMX2BUuxVZibW7US9hBukfi0qaBuk4P78e5ginZ+hXtZZx7li9yqs1Q27BkN+LmQ0Z6Zsbn/agq58GnuUGVwdlcilQ4WC6RoV7vtV/DIstAGDzuxA9ANrE6w6jOU25epq/OUvK7DNIm3U0PH3QK5wzYCubLuhH8tx9M7zcKJPodVPTOTsAO1RxTcwiyYTlNOg3yuubYPY+Lug1wpMPFR8WOfxSCSW9AUUTdm1Zfq7V5M99 "
		  }`,
		  						ExpectError: regexp.MustCompile(
		  							"Attribute value is not allowed to change",
		  						),
		  					},
		  				},
		  			})
		  		},
		  	)
	*/

	t.Run("changing the region triggers replacement", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + `
resource "leaseweb_public_cloud_instance" "test" {
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-2"
  }
  type = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  market_app_id = "newValue"
  type      = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type = {
    name = "lsw.m3.large"
  }
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
  region = {
    name = "eu-west-3"
  }
  type      = {
    name = "lsw.m3.large"
  }
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
