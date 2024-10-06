package instances

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
)

var (
	_ datasource.DataSource              = &instancesDataSource{}
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{}
}

type instancesDataSource struct {
	client client.Client
}

func (d *instancesDataSource) Configure(
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
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected provider.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = coreClient
}

func (d *instancesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_instances"
}

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	tflog.Info(ctx, "Read public cloud instances")
	instances, err := d.client.PublicCloudFacade.GetAllInstances(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read instances", err.Error())
		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read instances",
			err.Error(),
		)

		return
	}

	state := instances

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
