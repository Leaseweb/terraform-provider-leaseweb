package dedicatedserverresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dedicatedserver"

	tfResource "github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ tfResource.Resource                = &resource{}
	_ tfResource.ResourceWithConfigure   = &resource{}
	_ tfResource.ResourceWithImportState = &resource{}
)

type resource struct {
	dedicatedserver.API
}

func New() tfResource.Resource {
	return &resource{}
}

func (d *resource) Metadata(_ context.Context, req tfResource.MetadataRequest, resp *tfResource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server"
}

func (d *resource) Read(ctx context.Context, req tfResource.ReadRequest, resp *tfResource.ReadResponse) {
	var data resourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dedicatedServer, err := d.getServer(ctx, data.ID.ValueString())
	if err != nil {
		summary := "Reading dedicated server"
		resp.Diagnostics.AddError(summary, err.Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &dedicatedServer)
	resp.Diagnostics.Append(diags...)
}

func (d *resource) ImportState(ctx context.Context, req tfResource.ImportStateRequest, resp *tfResource.ImportStateResponse) {

	tfResource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)

	dedicatedServer, err := d.getServer(ctx, req.ID)
	if err != nil {
		summary := "Importing dedicated server"
		resp.Diagnostics.AddError(summary, err.Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
		return
	}

	diags := resp.State.Set(ctx, dedicatedServer)
	resp.Diagnostics.Append(diags...)
}

func (d *resource) Create(ctx context.Context, req tfResource.CreateRequest, resp *tfResource.CreateResponse) {
	panic("unimplemented")
}

func (d *resource) Delete(ctx context.Context, req tfResource.DeleteRequest, resp *tfResource.DeleteResponse) {
	panic("unimplemented")
}
