package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &dedicatedServersDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServersDataSource{}
)

type dedicatedServersDataSource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	name   string
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type dedicatedServersDataSourceData struct {
	Ids                   []types.String `tfsdk:"ids"`
	Reference             types.String   `tfsdk:"reference"`
	Ip                    types.String   `tfsdk:"ip"`
	MacAddress            types.String   `tfsdk:"mac_address"`
	Site                  types.String   `tfsdk:"site"`
	PrivateRackId         types.String   `tfsdk:"private_rack_id"`
	PrivateNetworkCapable types.String   `tfsdk:"private_network_capable"`
	PrivateNetworkEnabled types.String   `tfsdk:"private_network_enabled"`
}

func (d *dedicatedServersDataSource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configuration := dedicatedServer.NewConfiguration()

	// TODO: Refactor this part, ProviderData can be managed directly, not within client.
	coreClient, ok := req.ProviderData.(client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}
	d.apiKey = coreClient.ProviderData.ApiKey
	if coreClient.ProviderData.Host != nil {
		configuration.Host = *coreClient.ProviderData.Host
	}
	if coreClient.ProviderData.Scheme != nil {
		configuration.Scheme = *coreClient.ProviderData.Scheme
	}

	apiClient := dedicatedServer.NewAPIClient(configuration)
	d.client = apiClient.DedicatedServerAPI
}

func (d *dedicatedServersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.name)
}

func (d *dedicatedServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data dedicatedServersDataSourceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// NOTE: we show only latest 50 items.
	request := d.client.GetServerList(d.authContext(ctx)).Limit(50)

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
		summary := fmt.Sprintf("Reading data %s", d.name)
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}
	for _, server := range result.GetServers() {
		Ids = append(Ids, types.StringValue(server.GetId()))
	}

	data = dedicatedServersDataSourceData{
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

func (d *dedicatedServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func NewDedicatedServersDataSource() datasource.DataSource {
	return &dedicatedServersDataSource{
		name: "dedicated_servers",
	}
}
