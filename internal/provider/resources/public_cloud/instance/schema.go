package instance

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/handlers/public_cloud"
	customValidator "terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
)

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {

	handler := public_cloud.PublicCloudHandler{}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The instance unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"resources": schema.SingleNestedAttribute{
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
			"image": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Image ID",
						Validators: []validator.String{
							stringvalidator.OneOf(handler.GetImageIds()...),
						},
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
					stringvalidator.OneOf(enum.InstanceTypeC3Large.Values()...),
				},
			},
			"ssh_key": schema.StringAttribute{
				Optional:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "Public SSH key to be installed into the instance. Must be used only on Linux/FreeBSD instances",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(handler.GetSshKeyRegularExpression()),
						"Invalid ssh key",
					),
				},
			},
			"root_disk_size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances",
				Validators: []validator.Int64{
					int64validator.Between(
						handler.GetMinimumRootDiskSize(),
						handler.GetMaximumRootDiskSize(),
					),
				},
			},
			"root_disk_storage_type": schema.StringAttribute{
				Required:    true,
				Description: "The root disk's storage type",
				Validators: []validator.String{
					stringvalidator.OneOf(handler.GetRootDiskStorageTypes()...),
				},
			},
			"ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip":             schema.StringAttribute{Computed: true},
						"prefix_length":  schema.StringAttribute{Computed: true},
						"version":        schema.Int64Attribute{Computed: true},
						"null_routed":    schema.BoolAttribute{Computed: true},
						"main_ip":        schema.BoolAttribute{Computed: true},
						"network_type":   schema.StringAttribute{Computed: true},
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
							int64validator.OneOf(handler.GetBillingFrequencies()...),
						},
					},
					"term": schema.Int64Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when contract type is MONTHLY",
						Validators: []validator.Int64{
							int64validator.OneOf(handler.GetContractTerms()...),
						},
					},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(handler.GetContractTypes()...),
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
					},
				},
				Validators: []validator.Object{customValidator.ContractTermIsValid()},
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
			"auto_scaling_group": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The Auto Scaling Group unique identifier",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Auto Scaling Group type",
					},
					"state": schema.StringAttribute{
						Computed:    true,
						Description: "The Auto Scaling Group's current state.",
					},
					"desired_amount": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of instances that should be running",
					},
					"region": schema.StringAttribute{
						Computed:    true,
						Description: "The region in which the Auto Scaling Group was launched",
					},
					"reference": schema.StringAttribute{
						Computed:    true,
						Description: "The identifying name set to the auto scaling group",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "Date and time when the Auto Scaling Group was created",
					},
					"updated_at": schema.StringAttribute{
						Computed:    true,
						Description: "Date and time when the Auto Scaling Group was updated",
					},
					"starts_at": schema.StringAttribute{
						Computed:    true,
						Description: "Only for \"SCHEDULED\" auto scaling group. Date and time (UTC) that the instances need to be launched",
					},
					"ends_at": schema.StringAttribute{
						Computed:    true,
						Description: "Only for \"SCHEDULED\" auto scaling group. Date and time (UTC) that the instances need to be terminated",
					},
					"minimum_amount": schema.Int64Attribute{
						Computed:    true,
						Description: "The minimum number of instances that should be running",
					},
					"maximum_amount": schema.Int64Attribute{
						Computed:    true,
						Description: "Only for \"CPU_BASED\" auto scaling group. The maximum number of instances that can be running",
					},
					"cpu_threshold": schema.Int64Attribute{
						Computed:    true,
						Description: "Only for \"CPU_BASED\" auto scaling group. The target average CPU utilization for scaling",
					},
					"warmup_time": schema.Int64Attribute{
						Computed:    true,
						Description: "Only for \"CPU_BASED\" auto scaling group. Warm-up time in seconds for new instances",
					},
					"cooldown_time": schema.Int64Attribute{
						Computed:    true,
						Description: "Only for \"CPU_BASED\" auto scaling group. Cool-down time in seconds for new instances",
					},
					"load_balancer": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:    true,
								Description: "The load balancer unique identifier",
							},
							"type": schema.StringAttribute{
								Computed:    true,
								Description: "Load balancer type",
							},
							"resources": schema.SingleNestedAttribute{
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
							},
							"region": schema.StringAttribute{
								Computed:    true,
								Description: "The region where the load balancer was launched into",
							},
							"reference": schema.StringAttribute{
								Computed:    true,
								Description: "The identifying name set to the load balancer",
							},
							"state": schema.StringAttribute{
								Computed:    true,
								Description: "The load balancers current state",
							},
							"contract": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"billing_frequency": schema.Int64Attribute{
										Computed:    true,
										Description: "The billing frequency (in months) of the load balancer.",
									},
									"term": schema.Int64Attribute{
										Computed:    true,
										Description: "Contract term (in months). Used only when contract type is MONTHLY",
									},
									"type": schema.StringAttribute{
										Computed:    true,
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
									},
								},
							},
							"started_at": schema.StringAttribute{
								Computed:    true,
								Description: "Date and time when the load balancer was started for the first time, right after launching it",
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
							"load_balancer_configuration": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"balance": schema.StringAttribute{
										Computed: true,
									},
									"health_check": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"method": schema.StringAttribute{
												Computed: true,
											},
											"uri": schema.StringAttribute{
												Computed: true,
											},
											"host": schema.StringAttribute{
												Computed: true,
											},
											"port": schema.Int64Attribute{
												Computed: true,
											},
										},
									},
									"sticky_session": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Computed: true,
											},
											"max_life_time": schema.Int64Attribute{
												Computed: true,
											},
										},
									},
									"x_forwarded_for": schema.BoolAttribute{
										Computed: true,
									},
									"idle_timeout": schema.Int64Attribute{
										Computed: true,
									},
									"target_port": schema.Int64Attribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
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
