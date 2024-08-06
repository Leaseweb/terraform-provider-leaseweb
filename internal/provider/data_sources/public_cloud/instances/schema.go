package instances

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	facade "github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared_schemas/public_cloud"
)

func (d *instancesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	publicCloudFacade := facade.PublicCloudFacade{}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The instance unique identifier",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "The region where the instance was launched",
						},
						"reference": schema.StringAttribute{
							Computed:    true,
							Description: "The identifying name set to the instance",
						},
						"resources": public_cloud.Resources(),
						"image": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "Image ID",
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
								"state": schema.StringAttribute{
									Computed: true,
								},
								"state_reason": schema.StringAttribute{
									Computed: true,
								},
								"region": schema.StringAttribute{
									Computed: true,
								},
								"created_at": schema.StringAttribute{
									Computed: true,
								},
								"updated_at": schema.StringAttribute{
									Computed: true,
								},
								"custom": schema.BoolAttribute{
									Computed: true,
								},
								"market_apps": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"storage_types": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "The supported storage types",
								},
								"storage_size": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"size": schema.Float64Attribute{
											Computed:    true,
											Description: "The storage size",
										},
										"unit": schema.StringAttribute{
											Computed:    true,
											Description: "The storage size unit",
										},
									},
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
						"type": public_cloud.InstanceType(false),
						"root_disk_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances",
						},
						"root_disk_storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "The root disk's storage type",
						},
						"ips": public_cloud.Ips(),
						"started_at": schema.StringAttribute{
							Computed:    true,
							Description: "Date and time when the instance was started for the first time, right after launching it",
						},
						"contract": public_cloud.Contract(false, publicCloudFacade),
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
							Description: "Market App ID",
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
										"type":      public_cloud.InstanceType(false),
										"resources": public_cloud.Resources(),
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
										"contract": public_cloud.Contract(false, publicCloudFacade),
										"started_at": schema.StringAttribute{
											Computed:    true,
											Description: "Date and time when the load balancer was started for the first time, right after launching it",
										},
										"ips":             public_cloud.Ips(),
										"private_network": public_cloud.Network(),
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
						"volume": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"size": schema.Float64Attribute{
									Computed:    true,
									Description: "The Volume Size",
								},
								"unit": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"private_network": public_cloud.Network(),
					},
				},
			},
		},
	}
}
