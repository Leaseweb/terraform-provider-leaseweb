package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &targetGroupsDataSource{}
)

type targetGroupsDataSourceModel struct {
	ID           types.String                 `tfsdk:"id"`
	Name         types.String                 `tfsdk:"name"`
	Protocol     types.String                 `tfsdk:"protocol"`
	Port         types.Int32                  `tfsdk:"port"`
	Region       types.String                 `tfsdk:"region"`
	TargetGroups []targetGroupDataSourceModel `tfsdk:"target_groups"`
}

func (t targetGroupsDataSourceModel) generateRequest(
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ApiGetTargetGroupListRequest, error) {
	funcName := "generateRequest"

	request := api.GetTargetGroupList(ctx)
	if !t.ID.IsNull() {
		request = request.Id(t.ID.ValueString())
	}
	if !t.Name.IsNull() {
		request = request.Name(t.Name.ValueString())
	}
	if !t.Protocol.IsNull() {
		sdkProtocol, err := publicCloud.NewProtocolFromValue(t.Protocol.ValueString())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", funcName, err)
		}
		request = request.Protocol(*sdkProtocol)
	}
	if !t.Port.IsNull() {
		request = request.Port(t.Port.ValueInt32())
	}
	if !t.Region.IsNull() {
		sdkRegion, err := publicCloud.NewRegionNameFromValue(t.Region.ValueString())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", funcName, err)
		}
		request = request.Region(*sdkRegion)
	}

	return &request, nil
}

func adaptTargetGroupsToTargetGroupsDataSource(sdkTargetGroups []publicCloud.TargetGroup) targetGroupsDataSourceModel {
	targetGroups := targetGroupsDataSourceModel{}
	for _, sdkTargetGroup := range sdkTargetGroups {
		targetGroups.TargetGroups = append(
			targetGroups.TargetGroups,
			adaptTargetGroupToTargetGroupDataSource(sdkTargetGroup),
		)
	}

	return targetGroups
}

type targetGroupDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Protocol types.String `tfsdk:"protocol"`
	Port     types.Int32  `tfsdk:"port"`
	Region   types.String `tfsdk:"region"`
}

func adaptTargetGroupToTargetGroupDataSource(sdkTargetGroup publicCloud.TargetGroup) targetGroupDataSourceModel {
	return targetGroupDataSourceModel{
		ID:       basetypes.NewStringValue(sdkTargetGroup.GetId()),
		Name:     basetypes.NewStringValue(sdkTargetGroup.GetName()),
		Protocol: basetypes.NewStringValue(string(sdkTargetGroup.GetProtocol())),
		Port:     basetypes.NewInt32Value(sdkTargetGroup.GetPort()),
		Region:   basetypes.NewStringValue(string(sdkTargetGroup.GetRegion())),
	}
}

func getTargetGroups(request publicCloud.ApiGetTargetGroupListRequest) (
	[]publicCloud.TargetGroup,
	*http.Response,
	error,
) {
	var targetGroups []publicCloud.TargetGroup
	var offset *int32

	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("getTargetGroups: %w", err)
		}

		targetGroups = append(targetGroups, result.GetTargetGroups()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		request = request.Offset(*offset)
	}

	return targetGroups, nil, nil
}

type targetGroupsDataSource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func (t *targetGroupsDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, t.name)
}

func (t *targetGroupsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "Target group ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the target group",
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publicCloud.AllowedProtocolEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedProtocolEnumValues)...),
				},
			},
			"port": schema.Int32Attribute{
				Optional:    true,
				Description: "The port of the target group",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "Region name. Valid options are " + utils.StringTypeArrayToMarkdown(publicCloud.AllowedRegionNameEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedRegionNameEnumValues)...),
				},
			},
			"target_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Target group ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the target group",
						},
						"protocol": schema.StringAttribute{
							Computed: true,
						},
						"port": schema.Int32Attribute{
							Computed:    true,
							Description: "The port of the target group",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "Region name",
						},
					},
				},
			},
		},
	}
}

func (t *targetGroupsDataSource) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var config targetGroupsDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf(
		"Reading data %s for id %q",
		t.name,
		config.ID,
	)

	apiRequest, err := config.generateRequest(ctx, t.client)
	if err != nil {
		utils.HandleSdkError(
			summary,
			nil,
			err,
			&response.Diagnostics,
			ctx,
		)

		return
	}

	tflog.Info(ctx, summary)
	targetGroups, httpResponse, err := getTargetGroups(*apiRequest)
	if err != nil {
		utils.HandleSdkError(
			summary,
			httpResponse,
			err,
			&response.Diagnostics,
			ctx,
		)

		return
	}

	state := adaptTargetGroupsToTargetGroupsDataSource(targetGroups)
	state.ID = config.ID
	state.Name = config.Name
	state.Protocol = config.Protocol
	state.Port = config.Port
	state.Region = config.Region

	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (t *targetGroupsDataSource) Configure(
	_ context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected provider.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)

		return
	}

	t.client = coreClient.PublicCloudAPI
}

func NewTargetGroupsDataSource() datasource.DataSource {
	return &targetGroupsDataSource{
		name: "public_cloud_target_groups",
	}
}
