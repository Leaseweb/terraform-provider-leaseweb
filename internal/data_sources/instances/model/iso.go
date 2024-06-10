package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func newIso(sdkIso publicCloud.Iso) iso {
	return iso{
		Id:   utils.GenerateString(sdkIso.HasId(), sdkIso.GetId()),
		Name: utils.GenerateString(sdkIso.HasName(), sdkIso.GetName()),
	}
}
