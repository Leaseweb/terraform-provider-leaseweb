package dedicatedserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &serversDataSource{}
	_ datasource.DataSourceWithConfigure = &serversDataSource{}
)

type serversDataSource struct {
	utils.DataSourceAPI
}

type serversDataSourceModel struct {
	Ids                   []types.String `tfsdk:"ids"`
	Reference             types.String   `tfsdk:"reference"`
	Ip                    types.String   `tfsdk:"ip"`
	MacAddress            types.String   `tfsdk:"mac_address"`
	Site                  types.String   `tfsdk:"site"`
	PrivateRackId         types.String   `tfsdk:"private_rack_id"`
	PrivateNetworkCapable types.String   `tfsdk:"private_network_capable"`
	PrivateNetworkEnabled types.String   `tfsdk:"private_network_enabled"`
}

func (s *serversDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config serversDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	// NOTE: we show only the latest 50 items.
	request := s.DedicatedserverAPI.GetServerList(ctx).Limit(50)

	if !config.Reference.IsNull() && !config.Reference.IsUnknown() {
		request = request.Reference(config.Reference.ValueString())
	}

	if !config.Ip.IsNull() && !config.Ip.IsUnknown() {
		request = request.Ip(config.Ip.ValueString())
	}

	if !config.MacAddress.IsNull() && !config.MacAddress.IsUnknown() {
		request = request.MacAddress(config.MacAddress.ValueString())
	}

	if !config.Site.IsNull() && !config.Site.IsUnknown() {
		request = request.Site(config.Site.ValueString())
	}

	if !config.PrivateRackId.IsNull() && !config.PrivateRackId.IsUnknown() {
		request = request.PrivateRackId(config.PrivateRackId.ValueString())
	}

	if !config.PrivateNetworkCapable.IsNull() && !config.PrivateNetworkCapable.IsUnknown() {
		request = request.PrivateNetworkCapable(config.PrivateNetworkCapable.ValueString())
	}

	if !config.PrivateNetworkEnabled.IsNull() && !config.PrivateNetworkEnabled.IsUnknown() {
		request = request.PrivateNetworkEnabled(config.PrivateNetworkEnabled.ValueString())
	}

	var Ids []types.String

	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}
	for _, server := range result.GetServers() {
		Ids = append(Ids, types.StringValue(server.GetId()))
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			serversDataSourceModel{
				Ids:                   Ids,
				Reference:             config.Reference,
				Ip:                    config.Ip,
				MacAddress:            config.MacAddress,
				Site:                  config.Site,
				PrivateRackId:         config.PrivateRackId,
				PrivateNetworkCapable: config.PrivateNetworkCapable,
				PrivateNetworkEnabled: config.PrivateNetworkEnabled,
			},
		)...,
	)
}

func (s *serversDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of the dedicated server IDs available to the account.",
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list of servers by reference.",
			},
			"ip": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list of servers by ip address.",
			},
			"mac_address": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list of servers by mac address.",
			},
			"site": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list of servers by site (location).",
			},
			"private_rack_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list of servers by dedicated rack id.",
			},
			"private_network_capable": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list for private network capable servers.",
			},
			"private_network_enabled": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the list for private network enabled servers.",
			},
		},
	}
}

func NewServersDataSource() datasource.DataSource {
	return &serversDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "dedicated_servers",
		},
	}
}
