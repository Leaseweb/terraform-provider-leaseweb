package dedicatedserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &operatingSystemsDataSource{}
	_ datasource.DataSourceWithConfigure = &operatingSystemsDataSource{}
)

type operatingSystemsDataSource struct {
	utils.DedicatedserverDataSourceAPI
}

type operatingSystemDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type operatingSystemsDataSourceModel struct {
	OperatingSystems []operatingSystemDataSourceModel `tfsdk:"operating_systems"`
	ControlPanelId   types.String                     `tfsdk:"control_panel_id"`
}

func (o *operatingSystemsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config operatingSystemsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	request := o.Client.GetOperatingSystemList(ctx)
	if !config.ControlPanelId.IsNull() && !config.ControlPanelId.IsUnknown() {
		request = request.ControlPanelId(config.ControlPanelId.ValueString())
	}
	// NOTE: we show only the latest 50 items.
	result, response, err := request.Limit(50).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	var operatingSystems []operatingSystemDataSourceModel
	for _, os := range result.GetOperatingSystems() {
		operatingSystems = append(operatingSystems, operatingSystemDataSourceModel{
			Id:   basetypes.NewStringValue(os.GetId()),
			Name: basetypes.NewStringValue(os.GetName()),
		})
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			operatingSystemsDataSourceModel{
				OperatingSystems: operatingSystems,
				ControlPanelId:   config.ControlPanelId,
			},
		)...,
	)
}

func (o *operatingSystemsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
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

func NewOperatingSystemsDataSource() datasource.DataSource {
	return &operatingSystemsDataSource{
		DedicatedserverDataSourceAPI: utils.NewDedicatedserverDataSourceAPI("dedicated_server_operating_systems"),
	}
}
