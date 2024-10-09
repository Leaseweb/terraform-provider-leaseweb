// Package public_cloud implements schemas used multiple times in public_cloud data sources & resources.
package public_cloud

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	facade "github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
)

func Resources() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"cpu": schema.SingleNestedAttribute{
				Description: "Number of cores",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
			"memory": schema.SingleNestedAttribute{
				Description: "Total memory in GiB",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Float64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
			"public_network_speed": schema.SingleNestedAttribute{
				Description: "Public network speed in Gbps",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
			"private_network_speed": schema.SingleNestedAttribute{
				Description: "Private network speed in Gbps",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
		},
		Description: "Available resources",
		Computed:    true,
	}
}

func InstanceType(required bool) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: !required,
		Required: required,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    required,
				Computed:    !required,
				Description: "Type name",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
				},
			},
			"resources": Resources(),
			"storage_types": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The supported storage types",
			},
			"prices": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"currency": schema.StringAttribute{
						Computed: true,
					},
					"currency_symbol": schema.StringAttribute{
						Computed: true,
					},
					"compute": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"hourly_price": schema.StringAttribute{
								Computed: true,
							},
							"monthly_price": schema.StringAttribute{
								Computed: true,
							},
						},
						Computed: true,
					},
					"storage": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"local": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"hourly_price": schema.StringAttribute{
										Computed: true,
									},
									"monthly_price": schema.StringAttribute{
										Computed: true,
									},
								},
								Computed: true,
							},
							"central": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"hourly_price": schema.StringAttribute{
										Computed: true,
									},
									"monthly_price": schema.StringAttribute{
										Computed: true,
									},
								},
								Computed: true,
							},
						},
						Computed: true,
					},
				},
				Computed: true,
			},
		},
	}
}

func DataSourceIps() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"ip":            schema.StringAttribute{Computed: true},
				"prefix_length": schema.StringAttribute{Computed: true},
				"version":       schema.Int64Attribute{Computed: true},
				"null_routed":   schema.BoolAttribute{Computed: true},
				"main_ip":       schema.BoolAttribute{Computed: true},
				"network_type": schema.StringAttribute{
					Computed: true,
				},
				"reverse_lookup": schema.StringAttribute{Computed: true},
			},
		},
	}
}

func ResourceIps() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"ip":            schema.StringAttribute{Computed: true},
				"prefix_length": schema.StringAttribute{Computed: true},
				"version":       schema.Int64Attribute{Computed: true},
				"null_routed":   schema.BoolAttribute{Computed: true},
				"main_ip":       schema.BoolAttribute{Computed: true},
				"network_type": schema.StringAttribute{
					Computed: true,
				},
				"reverse_lookup": schema.StringAttribute{Computed: true},
				"ddos": schema.SingleNestedAttribute{
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"detection_profile": schema.StringAttribute{
							Computed: true,
						},
						"protection_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func Contract(
	required bool,
	facade facade.PublicCloudFacade,
) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: !required,
		Required: required,
		Attributes: map[string]schema.Attribute{
			"billing_frequency": schema.Int64Attribute{
				Computed:    !required,
				Required:    required,
				Description: "The billing frequency (in months). Valid options are " + facade.GetBillingFrequencies().Markdown(),
				Validators: []validator.Int64{
					int64validator.OneOf(facade.GetBillingFrequencies().ToInt64()...),
				},
			},
			"term": schema.Int64Attribute{
				Computed:    !required,
				Required:    required,
				Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + facade.GetContractTerms().Markdown(),
				Validators: []validator.Int64{
					int64validator.OneOf(facade.GetContractTerms().ToInt64()...),
				},
			},
			"type": schema.StringAttribute{
				Computed:    !required,
				Required:    required,
				Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
				Validators: []validator.String{
					stringvalidator.OneOf(facade.GetContractTypes()...),
				},
			},
			"ends_at": schema.StringAttribute{Computed: true},
			"renewals_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the contract will be automatically renewed",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the contract was created",
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
		},
		Validators: []validator.Object{customValidator.ContractTermIsValid()},
	}
}

func Network() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"subnet": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func ResourceRegion(required bool, warning string) resourceSchema.SingleNestedAttribute {
	printedWarning := ""

	if required {
		printedWarning = warning
	}

	return resourceSchema.SingleNestedAttribute{
		Required: required,
		Computed: !required,
		Attributes: map[string]resourceSchema.Attribute{
			"name": resourceSchema.StringAttribute{
				Required: required,
				Computed: !required,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Our current regions can be found in the [developer documentation](https://developer.leaseweb.com/api-docs/publiccloud_v1.html#tag/Instances/operation/launchInstance)" + printedWarning,
			},
			"location": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "The city where the region is located",
			},
		},
	}
}
