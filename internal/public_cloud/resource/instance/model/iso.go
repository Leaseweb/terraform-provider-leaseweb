package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (i Iso) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

func newIso(
	ctx context.Context,
	sdkIso publicCloud.Iso,
) (*Iso, diag.Diagnostics) {
	return &Iso{
		Id:   basetypes.NewStringValue(sdkIso.GetId()),
		Name: basetypes.NewStringValue(sdkIso.GetName()),
	}, nil
}
