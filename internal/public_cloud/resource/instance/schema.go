package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	customerValidator "terraform-provider-leaseweb/internal/public_cloud/validator"
)

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The instance unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"equipment_id": schema.StringAttribute{
				Computed:    true,
				Description: "Equipment's UUID",
			},
			"sales_org_id": schema.StringAttribute{
				Computed: true,
			},
			"customer_id": schema.StringAttribute{
				Computed:    true,
				Description: "The customer ID who owns the instance",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Region to launch the instance into",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The identifying name set to the instance",
			},
			"resource": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"cpu": schema.SingleNestedAttribute{
						Description: "Number of cores",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"value": schema.Int64Attribute{Computed: true},
							"unit":  schema.StringAttribute{Computed: true},
						},
					},
					"memory": schema.SingleNestedAttribute{
						Description: "Total memory in GiB",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"value": schema.Float64Attribute{Computed: true},
							"unit":  schema.StringAttribute{Computed: true},
						},
					},
					"public_network_speed": schema.SingleNestedAttribute{
						Description: "Public network speed in Gbps",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"value": schema.Int64Attribute{Computed: true},
							"unit":  schema.StringAttribute{Computed: true},
						},
					},
					"private_network_speed": schema.SingleNestedAttribute{
						Description: "Private network speed in Gbps",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"value": schema.Int64Attribute{Computed: true},
							"unit":  schema.StringAttribute{Computed: true},
						},
					},
				},
				Description: "i available for the load balancer",
				Computed:    true,
			},
			"operating_system": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Operating System ID",
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.StringAttribute{
						Computed: true,
					},
					"family": schema.StringAttribute{
						Computed: true,
					},
					"flavour": schema.StringAttribute{
						Computed: true,
					},
					"architecture": schema.StringAttribute{
						Computed: true,
					},
					"market_apps": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"storage_types": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "The supported storage types for the instance type",
					},
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The instance's current state",
				Validators: []validator.String{
					stringvalidator.OneOf(
						[]string{"RUNNING", "STOPPED", "CREATING", "DESTROYING", "DESTROYED"}...,
					),
				},
			},
			"product_type": schema.StringAttribute{
				Computed:    true,
				Description: "The product type",
			},
			"has_public_ipv4": schema.BoolAttribute{
				Computed: true,
			},
			"has_private_network": schema.BoolAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Instance type",
				Validators: []validator.String{
					stringvalidator.OneOf(
						[]string{
							"lsw.m3.large",
							"lsw.m3.xlarge",
							"lsw.m3.2xlarge",
							"lsw.m4.large",
							"lsw.m4.xlarge",
							"lsw.m4.2xlarge",
							"lsw.m4.4xlarge",
							"lsw.m5.large",
							"lsw.m5.xlarge",
							"lsw.m5.2xlarge",
							"lsw.m5.4xlarge",
							"lsw.m5a.large",
							"lsw.m5a.xlarge",
							"lsw.m5a.2xlarge",
							"lsw.m5a.4xlarge",
							"lsw.m5a.8xlarge",
							"lsw.m5a.12xlarge",
							"lsw.m6a.large",
							"lsw.m6a.xlarge",
							"lsw.m6a.2xlarge",
							"lsw.m6a.4xlarge",
							"lsw.m6a.8xlarge",
							"lsw.m6a.12xlarge",
							"lsw.m6a.16xlarge",
							"lsw.m6a.24xlarge",
							"lsw.c3.large",
							"lsw.c3.xlarge",
							"lsw.c3.2xlarge",
							"lsw.c3.4xlarge",
							"lsw.c4.large",
							"lsw.c4.xlarge",
							"lsw.c4.2xlarge",
							"lsw.c4.4xlarge",
							"lsw.c5.large",
							"lsw.c5.xlarge",
							"lsw.c5.2xlarge",
							"lsw.c5.4xlarge",
							"lsw.c5a.large",
							"lsw.c5a.xlarge",
							"lsw.c5a.2xlarge",
							"lsw.c5a.4xlarge",
							"lsw.c5a.9xlarge",
							"lsw.c5a.12xlarge",
							"lsw.c6a.large",
							"lsw.c6a.xlarge",
							"lsw.c6a.2xlarge",
							"lsw.c6a.4xlarge",
							"lsw.c6a.8xlarge",
							"lsw.c6a.12xlarge",
							"lsw.c6a.16xlarge",
							"lsw.c6a.24xlarge",
							"lsw.r3.large",
							"lsw.r3.xlarge",
							"lsw.r3.2xlarge",
							"lsw.r4.large",
							"lsw.r4.xlarge",
							"lsw.r4.2xlarge",
							"lsw.r5.large",
							"lsw.r5.xlarge",
							"lsw.r5.2xlarge",
							"lsw.r5a.large",
							"lsw.r5a.xlarge",
							"lsw.r5a.2xlarge",
							"lsw.r5a.4xlarge",
							"lsw.r5a.8xlarge",
							"lsw.r5a.12xlarge",
							"lsw.r6a.large",
							"lsw.r6a.xlarge",
							"lsw.r6a.2xlarge",
							"lsw.r6a.4xlarge",
							"lsw.r6a.8xlarge",
							"lsw.r6a.12xlarge",
							"lsw.r6a.16xlarge",
							"lsw.r6a.24xlarge",
						}...,
					),
				},
			},
			"ssh_key": schema.StringAttribute{
				Optional:    true,
				Description: "Public SSH key to be installed into the instance. Must be used only on Linux/FreeBSD instances",
			},
			"root_disk_size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances",
				Validators: []validator.Int64{
					int64validator.Between(5, 1000),
				},
			},
			"root_disk_storage_type": schema.StringAttribute{
				Required:    true,
				Description: "The root disk's storage type",
				Validators: []validator.String{
					stringvalidator.OneOf(
						[]string{"LOCAL", "CENTRAL"}...,
					),
				},
			},
			"ips": schema.ListNestedAttribute{
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
							Validators: []validator.String{
								stringvalidator.OneOf(
									[]string{"INTERNAL", "PUBLIC"}...,
								),
							},
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
			},
			"started_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date and time when the instance was started for the first time, right after launching it",
			},
			"contract": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"billing_frequency": schema.Int64Attribute{
						Required:    true,
						Description: "The billing frequency (in months) of the instance.",
						Validators: []validator.Int64{
							int64validator.OneOf(
								[]int64{0, 1, 3, 6, 12}...,
							),
						},
					},
					"term": schema.Int64Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when contract type is MONTHLY",
						Validators: []validator.Int64{
							int64validator.OneOf(
								[]int64{0, 1, 3, 6, 12}...,
							),
						},
					},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								[]string{"HOURLY", "MONTHLY"}...,
							),
						},
						Description: "Select HOURLY for billing based on hourly usage, else MONTHLY for billing per month usage",
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
						Validators: []validator.String{
							stringvalidator.OneOf(
								[]string{"ACTIVE", "DELETE_SCHEDULED"}...,
							),
						},
					},
				},
				Validators: []validator.Object{customerValidator.ContractTermIsValid()},
			},
			"iso": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"market_app_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Market App ID that must be installed into the instance",
			},
			"private_network": schema.SingleNestedAttribute{
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
			},
		},
	}
}
