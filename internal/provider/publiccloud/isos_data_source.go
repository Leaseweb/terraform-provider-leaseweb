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
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &isosDataSource{}
)

type isoDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type isosDataSourceModel struct {
	ISOs []isoDataSourceModel `tfsdk:"isos"`
}

func adaptIsosToISOsDataSource(sdkISOs []publiccloud.Iso) isosDataSourceModel {
	var isos isosDataSourceModel

	for _, iso := range sdkISOs {
		isos.ISOs = append(isos.ISOs, adaptIsoToISODataSource(iso))
	}

	return isos
}

func adaptIsoToISODataSource(iso publiccloud.Iso) isoDataSourceModel {
	return isoDataSourceModel{
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

type isosDataSource struct {
	utils.PubliccloudDataSourceAPI
}

func (i *isosDataSource) Schema(
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

func (i *isosDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	ISOs, httpResponse, err := getISOs(ctx, i.Client)
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			adaptIsosToISOsDataSource(ISOs),
		)...,
	)
}

func NewISOsDataSource() datasource.DataSource {
	return &isosDataSource{
		PubliccloudDataSourceAPI: utils.PubliccloudDataSourceAPI{
			Name: "public_cloud_isos",
		},
	}
}
