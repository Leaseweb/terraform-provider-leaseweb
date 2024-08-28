package operating_systems

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
)

func (d *operatingSystemsDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Info(ctx, "Reading dedicated server operating systems")
	operatingSystems, err := d.client.DedicatedServerFacade.GetAllOperatingSystems(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read dedicated_server_operating_systems", err.Error())
		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read dedicated_server_operating_systems",
			err.Error(),
		)

		return
	}

	state := operatingSystems

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
