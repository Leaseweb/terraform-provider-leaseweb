package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &dedicatedServersDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServersDataSource{}
)

type dedicatedServersDataSource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
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
	resp.TypeName = req.ProviderTypeName + "_dedicated_servers"
}

func (d *dedicatedServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var filter dedicatedServersDataSourceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &filter)...)
	// NOTE: we show only latest 50 items.
	request := d.client.GetServerList(d.authContext(ctx)).Limit(50)

	if !filter.Reference.IsNull() && !filter.Reference.IsUnknown() {
		request = request.Reference(filter.Reference.ValueString())
	}

	if !filter.Ip.IsNull() && !filter.Ip.IsUnknown() {
		request = request.Ip(filter.Ip.ValueString())
	}

	if !filter.MacAddress.IsNull() && !filter.MacAddress.IsUnknown() {
		request = request.MacAddress(filter.MacAddress.ValueString())
	}

	if !filter.Site.IsNull() && !filter.Site.IsUnknown() {
		request = request.Site(filter.Site.ValueString())
	}

	if !filter.PrivateRackId.IsNull() && !filter.PrivateRackId.IsUnknown() {
		request = request.PrivateRackId(filter.PrivateRackId.ValueString())
	}

	if !filter.PrivateNetworkCapable.IsNull() && !filter.PrivateNetworkCapable.IsUnknown() {
		request = request.PrivateNetworkCapable(filter.PrivateNetworkCapable.ValueString())
	}

	if !filter.PrivateNetworkEnabled.IsNull() && !filter.PrivateNetworkEnabled.IsUnknown() {
		request = request.PrivateNetworkEnabled(filter.PrivateNetworkEnabled.ValueString())
	}

	var Ids []types.String

	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading dedicated server", err.Error())
		return
	}
	for _, server := range result.GetServers() {
		Ids = append(Ids, types.StringValue(server.GetId()))
	}

	data := dedicatedServersDataSourceData{
		Ids:                   Ids,
		Reference:             filter.Reference,
		Ip:                    filter.Ip,
		MacAddress:            filter.MacAddress,
		Site:                  filter.Site,
		PrivateRackId:         filter.PrivateRackId,
		PrivateNetworkCapable: filter.PrivateNetworkCapable,
		PrivateNetworkEnabled: filter.PrivateNetworkEnabled,
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
			},
			"reference": schema.StringAttribute{
				Optional: true,
			},
			"ip": schema.StringAttribute{
				Optional: true,
			},
			"mac_address": schema.StringAttribute{
				Optional: true,
			},
			"site": schema.StringAttribute{
				Optional: true,
			},
			"private_rack_id": schema.StringAttribute{
				Optional: true,
			},
			"private_network_capable": schema.StringAttribute{
				Optional: true,
			},
			"private_network_enabled": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func NewDedicatedServersDataSource() datasource.DataSource {
	return &dedicatedServersDataSource{}
}
