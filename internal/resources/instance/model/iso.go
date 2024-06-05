package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type Iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func newIso(sdkIso publicCloud.Iso) Iso {
	return Iso{
		Id:   resources.GetStringValue(sdkIso.HasId(), sdkIso.GetId()),
		Name: resources.GetStringValue(sdkIso.HasName(), sdkIso.GetName()),
	}
}
