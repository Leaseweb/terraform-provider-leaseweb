package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *dedicatedServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the server.",
			},
			"asset_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Asset Id of the server.",
			},
			"serial_number": schema.StringAttribute{
				Computed:    true,
				Description: "Serial number of server.",
			},
			"contract_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the contract.",
			},
			"rack_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Id of the rack.",
			},
			"rack_capacity": schema.StringAttribute{
				Computed:    true,
				Description: "The capacity of the rack.",
			},
			"rack_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the rack.",
			},
			"is_automation_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if automation feature is available for the server.",
			},
			"is_ipmi_reboot_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if ipmi_reboot feature is available for the server.",
			},
			"is_power_cycle_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if power_cycle feature is available for the server.",
			},
			"is_private_network_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if private network feature is available for the server.",
			},
			"is_remote_management_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if remote management feature is available for the server.",
			},
			"location_rack": schema.StringAttribute{
				Computed: true,
			},
			"location_site": schema.StringAttribute{
				Computed:    true,
				Description: "The site of the location.",
			},
			"location_suite": schema.StringAttribute{
				Computed:    true,
				Description: "The suite of the location.",
			},
			"location_unit": schema.StringAttribute{
				Computed:    true,
				Description: "The unit of the location.",
			},
			"public_mac": schema.StringAttribute{
				Computed:    true,
				Description: "Public mac address.",
			},
			"public_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Public ip address.",
			},
			"public_gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Public gateway.",
			},
			"internal_mac": schema.StringAttribute{
				Computed:    true,
				Description: "Internal mac address.",
			},
			"internal_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Internal ip address.",
			},
			"internal_gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Internal gateway.",
			},
			"remote_mac": schema.StringAttribute{
				Computed:    true,
				Description: "Remote mac address.",
			},
			"remote_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Remote ip address.",
			},
			"remote_gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Remote gateway.",
			},
			"ram_size": schema.Int32Attribute{
				Computed:    true,
				Description: "The size of the ram.",
			},
			"ram_unit": schema.StringAttribute{
				Computed:    true,
				Description: "The unit of the ram.",
			},
			"cpu_quantity": schema.Int32Attribute{
				Computed:    true,
				Description: "The quantity of the cpu.",
			},
			"cpu_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the cpu.",
			},
		},
	}
}
