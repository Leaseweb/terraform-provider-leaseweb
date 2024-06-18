package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (i Iso) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

func newIso(sdkIso publicCloud.Iso) Iso {
	return Iso{
		Id:   basetypes.NewStringValue(sdkIso.GetId()),
		Name: basetypes.NewStringValue(sdkIso.GetName()),
	}
}
