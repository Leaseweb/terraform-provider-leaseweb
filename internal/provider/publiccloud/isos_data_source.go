package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &ISOsDataSource{}
)

type ISODataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ISOsDataSourceModel struct {
	ISOs []ISODataSourceModel `tfsdk:"isos"`
}

func adaptIsosToISOsDataSource(sdkISOs []publiccloud.Iso) ISOsDataSourceModel {
	var isos ISOsDataSourceModel

	for _, iso := range sdkISOs {
		isos.ISOs = append(isos.ISOs, adaptIsoToISODataSource(iso))
	}

	return isos
}

func adaptIsoToISODataSource(iso publiccloud.Iso) ISODataSourceModel {
	return ISODataSourceModel{
		ID:   basetypes.NewStringValue(iso.GetId()),
		Name: basetypes.NewStringValue(iso.GetName()),
	}
}

func getISOs(
	ctx context.Context,
	api publiccloud.PubliccloudAPI,
) ([]publiccloud.Iso, *http.Response, error) {
	var isos []publiccloud.Iso
	var offset *int32

	request := api.GetIsoList(ctx)

	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("getISOs: %w", err)
		}
		isos = append(isos, result.Isos...)

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

	return isos, nil, nil
}

type ISOsDataSource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func (i *ISOsDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, i.name)
}

func (i *ISOsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"isos": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ISO ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of ISO",
						},
					},
				},
			},
		},
	}
}

func (i *ISOsDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	ISOs, httpResponse, err := getISOs(ctx, i.client)
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, fmt.Sprintf("Reading data %s", i.name), err, httpResponse)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			adaptIsosToISOsDataSource(ISOs),
		)...,
	)
}

func (i *ISOsDataSource) Configure(
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

	i.client = coreClient.PubliccloudAPI
}

func NewISOsDataSource() datasource.DataSource {
	return &ISOsDataSource{
		name: "public_cloud_isos",
	}
}
