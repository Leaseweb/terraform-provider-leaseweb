package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
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
		Id:   utils.GenerateString(sdkIso.HasId(), sdkIso.GetId()),
		Name: utils.GenerateString(sdkIso.HasName(), sdkIso.GetName()),
	}
}
