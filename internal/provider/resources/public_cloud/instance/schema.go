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
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	sharedSchemas "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared_schemas/public_cloud"
)

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	facade := public_cloud.PublicCloudFacade{}
	warningError := "**WARNING!** Changing this value once running will cause this instance to be destroyed and a new one to be created."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The instance unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": sharedSchemas.ResourceRegion(true, warningError),
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The identifying name set to the instance",
			},
			"resources": sharedSchemas.Resources(),
			"image": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Can be either an Operating System or a UUID in case of a Custom Image ID." + warningError,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
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
					"state": schema.StringAttribute{
						Computed: true,
					},
					"state_reason": schema.StringAttribute{
						Computed: true,
					},
					"region": sharedSchemas.ResourceRegion(false, warningError),
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"updated_at": schema.StringAttribute{
						Computed: true,
					},
					"custom": schema.BoolAttribute{
						Computed:    true,
						Description: "Standard or Custom image",
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
			"has_user_data": schema.BoolAttribute{
				Computed: true,
			},
			"type": sharedSchemas.InstanceType(true),
			// TODO Enable SSH key support
			/**
			  "ssh_key": schema.StringAttribute{
			  	Optional:      true,
			  	Sensitive:     true,
			  	PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			  	Description:   "Public SSH key to be installed into the instance. Must be used only on Linux/FreeBSD instances",
			  	Validators: []validator.String{
			  		stringvalidator.RegexMatches(
			  			regexp.MustCompile(facade.GetSshKeyRegularExpression()),
			  			"Invalid ssh key",
			  		),
			  	},
			  },
			*/
			"root_disk_size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances. The maximum size is 1000 GB",
				Validators: []validator.Int64{
					int64validator.Between(
						facade.GetMinimumRootDiskSize(),
						facade.GetMaximumRootDiskSize(),
					),
				},
			},
			"root_disk_storage_type": schema.StringAttribute{
				Required:    true,
				Description: "The root disk's storage type. Can be *LOCAL* or *CENTRAL*. " + warningError,
				Validators: []validator.String{
					stringvalidator.OneOf(facade.GetRootDiskStorageTypes()...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ips": sharedSchemas.Ips(),
			"started_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date and time when the instance was started for the first time, right after launching it",
			},
			"contract": sharedSchemas.Contract(true, facade),
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
				Description: "Market App ID that must be installed into the instance." + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
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
					"region": sharedSchemas.ResourceRegion(false, warningError),
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
				},
			},
			"private_network": sharedSchemas.Network(),
		},
	}
}
