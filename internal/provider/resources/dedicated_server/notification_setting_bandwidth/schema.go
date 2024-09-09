package notification_setting_bandwidth

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/notification_setting_bandwidth/validator"
)

func (n *notificationSettingBandwidthResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	facade := dedicated_server.DedicatedServerFacade{}
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_id": schema.StringAttribute{
				Required:    true,
				Description: "The server unique identifier",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The notification setting bandwidth unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"frequency": schema.StringAttribute{
				Required:    true,
				Description: "The notification frequency. Valid options can be *DAILY* or *WEEKLY* or *MONTHLY*.",
				Validators: []validator.String{
					stringvalidator.OneOf(facade.GetFrequencies()...),
				},
			},
			"last_checked_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date timestamp when the system last checked the server for threshold limit",
			},
			"threshold": schema.StringAttribute{
				Required:    true,
				Description: "Threshold Value. Value can be a number greater than 0.",
				Validators: []validator.String{
					customValidator.GreaterThanZero(),
				},
			},
			"threshold_exceeded_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date timestamp when the threshold exceeded the limit",
			},
			"unit": schema.StringAttribute{
				Required:    true,
				Description: "The notification unit. Valid options can be *Mbps* or *Gbps*.",
				Validators: []validator.String{
					stringvalidator.OneOf(facade.GetUnits()...),
				},
			},
			"actions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"last_triggered_at": schema.StringAttribute{Computed: true, Description: "Date timestamp when the last notification email triggered"},
						"type":              schema.StringAttribute{Computed: true, Description: "The type of the action"},
					},
				},
			},
		},
	}
}
