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
	customerValidator "terraform-provider-leaseweb/internal/public_cloud/resource/instance/validator"
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
			"operating_system": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Operating System ID",
						Validators: []validator.String{
							stringvalidator.OneOf(
								[]string{
									"ALMALINUX_8_64BIT",
									"ALMALINUX_9_64BIT",
									"ARCH_LINUX_64BIT",
									"CENTOS_7_64BIT",
									"DEBIAN_10_64BIT",
									"DEBIAN_11_64BIT",
									"DEBIAN_12_64BIT",
									"FREEBSD_13_64BIT",
									"FREEBSD_14_64BIT",
									"ROCKY_LINUX_8_64BIT",
									"ROCKY_LINUX_9_64BIT",
									"UBUNTU_20_04_64BIT",
									"UBUNTU_22_04_64BIT",
									"UBUNTU_24_04_64BIT",
									"WINDOWS_SERVER_2016_STANDARD_64BIT",
									"WINDOWS_SERVER_2019_STANDARD_64BIT",
									"WINDOWS_SERVER_2022_STANDARD_64BIT",
								}...,
							),
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
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(
							`^(ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519|ssh-rsa)\s+(?:[a-zA-Z0-9+/]{4})*(?:|[a-zA-Z0-9+/]{3}=|[a-zA-Z0-9+/]{2}==|[a-zA-Z0-9+/]===)[\s+\x21-\x7F]+$`),
						"Invalid ssh key",
					),
				},
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
								Validators: []validator.String{
									stringvalidator.OneOf(
										[]string{"HOURLY", "MONTHLY"}...,
									),
								},
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
