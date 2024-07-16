package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type privateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func newPrivateNetwork(entityPrivateNetwork domain.PrivateNetwork) *privateNetwork {
	return &privateNetwork{
		Id:     basetypes.NewStringValue(entityPrivateNetwork.Id),
		Status: basetypes.NewStringValue(entityPrivateNetwork.Status),
		Subnet: basetypes.NewStringValue(entityPrivateNetwork.Subnet),
	}
}
