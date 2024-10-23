package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &dedicatedServerControlPanelsDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServerControlPanelsDataSource{}
)

type dedicatedServerControlPanelsDataSource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type controlPanel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type dedicatedServerControlPanelsDataSourceData struct {
	ControlPanels     []controlPanel `tfsdk:"control_panels"`
	OperatingSystemId types.String   `tfsdk:"operating_system_id"`
}

func (d *dedicatedServerControlPanelsDataSource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerControlPanelsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dedicatedServerControlPanelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_control_panels"
}

func (d *dedicatedServerControlPanelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data dedicatedServerControlPanelsDataSourceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	var controlPanels []controlPanel
	var result *dedicatedServer.ControlPanelList
	var response *http.Response
	var err error

	// NOTE: we show only latest 50 items.
	if !data.OperatingSystemId.IsNull() && !data.OperatingSystemId.IsUnknown() {
		request := d.client.GetControlPanelListByOperatingSystemId(d.authContext(ctx), data.OperatingSystemId.ValueString()).Limit(50)
		result, response, err = request.Execute()
	} else {
		request := d.client.GetControlPanelList(d.authContext(ctx)).Limit(50)
		result, response, err = request.Execute()
	}

	if err != nil {
		summary := "Error reading control panels"
		resp.Diagnostics.AddError(summary, NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, NewError(response, err).Error()))
		return
	}

	for _, cp := range result.GetControlPanels() {
		controlPanels = append(controlPanels, controlPanel{
			Id:   basetypes.NewStringValue(cp.GetId()),
			Name: basetypes.NewStringValue(cp.GetName()),
		})
	}

	newData := dedicatedServerControlPanelsDataSourceData{
		ControlPanels:     controlPanels,
		OperatingSystemId: data.OperatingSystemId,
	}

	diags := resp.State.Set(ctx, &newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dedicatedServerControlPanelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"control_panels": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Id of the control panel.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the control panel.",
						},
					},
				},
			},
			"operating_system_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter control panels by operating system id.",
			},
		},
	}
}

func NewDedicatedServerControlPanelsDataSource() datasource.DataSource {
	return &dedicatedServerControlPanelsDataSource{}
}
