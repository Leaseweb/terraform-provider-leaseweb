package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/customerror"
)

var (
	_ datasource.DataSource              = &dedicatedServerOperatingSystemsDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServerOperatingSystemsDataSource{}
)

type dedicatedServerOperatingSystemsDataSource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type operatingSystem struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type dedicatedServerOperatingSystemsDataSourceData struct {
	OperatingSystems []operatingSystem `tfsdk:"operating_systems"`
	ControlPanelId   types.String      `tfsdk:"control_panel_id"`
}

func (d *dedicatedServerOperatingSystemsDataSource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerOperatingSystemsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dedicatedServerOperatingSystemsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_operating_systems"
}

func (d *dedicatedServerOperatingSystemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data dedicatedServerOperatingSystemsDataSourceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	request := d.client.GetOperatingSystemList(d.authContext(ctx))
	if !data.ControlPanelId.IsNull() && !data.ControlPanelId.IsUnknown() {
		request = request.ControlPanelId(data.ControlPanelId.ValueString())
	}
	// NOTE: we show only latest 50 items.
	result, response, err := request.Limit(50).Execute()
	if err != nil {
		summary := "Error reading control panels"
		resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
		return
	}

	var operatingSystems []operatingSystem
	for _, os := range result.GetOperatingSystems() {
		operatingSystems = append(operatingSystems, operatingSystem{
			Id:   basetypes.NewStringValue(os.GetId()),
			Name: basetypes.NewStringValue(os.GetName()),
		})
	}

	newData := dedicatedServerOperatingSystemsDataSourceData{
		OperatingSystems: operatingSystems,
		ControlPanelId:   data.ControlPanelId,
	}

	diags := resp.State.Set(ctx, &newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dedicatedServerOperatingSystemsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"operating_systems": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Id of the operating system.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Id of the operating system.",
						},
					},
				},
			},
			"control_panel_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter operating systems by control panel id.",
			},
		},
	}
}

func NewDedicatedServerOperatingSystemsDataSource() datasource.DataSource {
	return &dedicatedServerOperatingSystemsDataSource{}
}
