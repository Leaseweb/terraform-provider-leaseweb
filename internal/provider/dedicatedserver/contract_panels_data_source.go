package dedicatedserver

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &controlPanelsDataSource{}
	_ datasource.DataSourceWithConfigure = &controlPanelsDataSource{}
)

type controlPanelsDataSource struct {
	utils.DataSourceAPI
}

type controlPanelDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type controlPanelsDataSourceModel struct {
	ControlPanels     []controlPanelDataSourceModel `tfsdk:"control_panels"`
	OperatingSystemId types.String                  `tfsdk:"operating_system_id"`
}

func (c *controlPanelsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config controlPanelsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	var controlPanels []controlPanelDataSourceModel
	var result *dedicatedserver.ControlPanelList
	var response *http.Response
	var err error

	// NOTE: we show only the latest 50 items.
	if !config.OperatingSystemId.IsNull() && !config.OperatingSystemId.IsUnknown() {
		request := c.DedicatedserverAPI.GetControlPanelListByOperatingSystemId(
			ctx,
			config.OperatingSystemId.ValueString(),
		).Limit(50)
		result, response, err = request.Execute()
	} else {
		request := c.DedicatedserverAPI.GetControlPanelList(ctx).Limit(50)
		result, response, err = request.Execute()
	}

	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	for _, cp := range result.GetControlPanels() {
		controlPanels = append(controlPanels, controlPanelDataSourceModel{
			ID:   basetypes.NewStringValue(cp.GetId()),
			Name: basetypes.NewStringValue(cp.GetName()),
		})
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			controlPanelsDataSourceModel{
				ControlPanels:     controlPanels,
				OperatingSystemId: config.OperatingSystemId,
			},
		)...,
	)
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
							Description: "ID of the control panel.",
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
		DataSourceAPI: utils.DataSourceAPI{
			Name: "dedicated_server_control_panels",
		},
	}
}
