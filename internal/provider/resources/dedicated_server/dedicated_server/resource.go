package dedicated_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource                = &dedicatedServerResource{}
	_ resource.ResourceWithConfigure   = &dedicatedServerResource{}
	_ resource.ResourceWithImportState = &dedicatedServerResource{}
)

func NewDedicatedServerResource() resource.Resource {
	return &dedicatedServerResource{}
}

type dedicatedServerResource struct {
	client client.Client
}

// Create implements resource.Resource.
func (d *dedicatedServerResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (d *dedicatedServerResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Update implements resource.Resource.
func (d *dedicatedServerResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
