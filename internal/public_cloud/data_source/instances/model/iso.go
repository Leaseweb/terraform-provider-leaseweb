package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func newIso(entityIso entity.Iso) *iso {
	return &iso{
		Id:   basetypes.NewStringValue(entityIso.Id),
		Name: basetypes.NewStringValue(entityIso.Name),
	}
}
