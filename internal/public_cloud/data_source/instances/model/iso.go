package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func newIso(sdkIso publicCloud.Iso) *iso {
	return &iso{
		Id:   basetypes.NewStringValue(sdkIso.GetId()),
		Name: basetypes.NewStringValue(sdkIso.GetName()),
	}
}
