package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type PrivateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func (p PrivateNetwork) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.StringType,
		"status": types.StringType,
		"subnet": types.StringType,
	}
}

func newPrivateNetwork(
	ctx context.Context,
	privateNetwork entity.PrivateNetwork,
) (*PrivateNetwork, diag.Diagnostics) {
	return &PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}, nil
}
