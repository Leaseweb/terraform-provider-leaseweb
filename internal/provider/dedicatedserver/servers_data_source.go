package dedicatedserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &serversDataSource{}
	_ datasource.DataSourceWithConfigure = &serversDataSource{}
)

type serversDataSource struct {
	name   string
	client dedicatedServer.DedicatedServerAPI
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

func (s *serversDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	coreClient, ok := utils.GetDataSourceClient(req, resp)
	if !ok {
		return
	}

	s.client = coreClient.DedicatedServerAPI
}

func (s *serversDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, s.name)
}

func (s *serversDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	var data serversDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// NOTE: we show only the latest 50 items.
	request := s.client.GetServerList(ctx).Limit(50)

	if !data.Reference.IsNull() && !data.Reference.IsUnknown() {
		request = request.Reference(data.Reference.ValueString())
	}

	if !data.Ip.IsNull() && !data.Ip.IsUnknown() {
		request = request.Ip(data.Ip.ValueString())
	}

	if !data.MacAddress.IsNull() && !data.MacAddress.IsUnknown() {
		request = request.MacAddress(data.MacAddress.ValueString())
	}

	if !data.Site.IsNull() && !data.Site.IsUnknown() {
		request = request.Site(data.Site.ValueString())
	}

	if !data.PrivateRackId.IsNull() && !data.PrivateRackId.IsUnknown() {
		request = request.PrivateRackId(data.PrivateRackId.ValueString())
	}

	if !data.PrivateNetworkCapable.IsNull() && !data.PrivateNetworkCapable.IsUnknown() {
		request = request.PrivateNetworkCapable(data.PrivateNetworkCapable.ValueString())
	}

	if !data.PrivateNetworkEnabled.IsNull() && !data.PrivateNetworkEnabled.IsUnknown() {
		request = request.PrivateNetworkEnabled(data.PrivateNetworkEnabled.ValueString())
	}

	var Ids []types.String

	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Reading data %s", s.name)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}
	for _, server := range result.GetServers() {
		Ids = append(Ids, types.StringValue(server.GetId()))
	}

	data = serversDataSourceModel{
		Ids:                   Ids,
		Reference:             data.Reference,
		Ip:                    data.Ip,
		MacAddress:            data.MacAddress,
		Site:                  data.Site,
		PrivateRackId:         data.PrivateRackId,
		PrivateNetworkCapable: data.PrivateNetworkCapable,
		PrivateNetworkEnabled: data.PrivateNetworkEnabled,
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
		name: "dedicated_servers",
	}
}
