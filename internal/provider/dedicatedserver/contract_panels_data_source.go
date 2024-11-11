package dedicatedserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &controlPanelsDataSource{}
	_ datasource.DataSourceWithConfigure = &controlPanelsDataSource{}
)

type controlPanelsDataSource struct {
	name   string
	client dedicatedServer.DedicatedServerAPI
}

type controlPanelDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type controlPanelsDataSourceModel struct {
	ControlPanels     []controlPanelDataSourceModel `tfsdk:"control_panels"`
	OperatingSystemId types.String                  `tfsdk:"operating_system_id"`
}

func (c *controlPanelsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

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

	c.client = coreClient.DedicatedServerAPI
}

func (c *controlPanelsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	utils.SetDataSourceTypeName(resp, req, c.name)
}

func (c *controlPanelsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	var data controlPanelsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	var controlPanels []controlPanelDataSourceModel
	var result *dedicatedServer.ControlPanelList
	var response *http.Response
	var err error

	// NOTE: we show only the latest 50 items.
	if !data.OperatingSystemId.IsNull() && !data.OperatingSystemId.IsUnknown() {
		request := c.client.GetControlPanelListByOperatingSystemId(
			ctx,
			data.OperatingSystemId.ValueString(),
		).Limit(50)
		result, response, err = request.Execute()
	} else {
		request := c.client.GetControlPanelList(ctx).Limit(50)
		result, response, err = request.Execute()
	}

	if err != nil {
		summary := fmt.Sprintf("Reading data %s", c.name)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	for _, cp := range result.GetControlPanels() {
		controlPanels = append(controlPanels, controlPanelDataSourceModel{
			Id:   basetypes.NewStringValue(cp.GetId()),
			Name: basetypes.NewStringValue(cp.GetName()),
		})
	}

	newData := controlPanelsDataSourceModel{
		ControlPanels:     controlPanels,
		OperatingSystemId: data.OperatingSystemId,
	}

	diags := resp.State.Set(ctx, &newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *controlPanelsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
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

func NewControlPanelsDataSource() datasource.DataSource {
	return &controlPanelsDataSource{
		name: "dedicated_server_control_panels",
	}
}
