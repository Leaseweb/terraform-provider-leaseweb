package instance

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
	_ resource.ResourceWithModifyPlan  = &instanceResource{}
)

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResource struct {
	client client.Client
}
