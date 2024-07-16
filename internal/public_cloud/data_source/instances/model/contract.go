package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/utils"
)

type contract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	RenewalsAt       types.String `tfsdk:"renewals_at"`
	CreatedAt        types.String `tfsdk:"created_at"`
	State            types.String `tfsdk:"state"`
}

func newContract(entityContract domain.Contract) contract {
	return contract{
		BillingFrequency: basetypes.NewInt64Value(
			int64(entityContract.BillingFrequency),
		),
		Term:       basetypes.NewInt64Value(int64(entityContract.Term)),
		Type:       basetypes.NewStringValue(string(entityContract.Type)),
		EndsAt:     utils.ConvertNullableTimeToStringValue(entityContract.EndsAt),
		RenewalsAt: basetypes.NewStringValue(entityContract.RenewalsAt.String()),
		CreatedAt:  basetypes.NewStringValue(entityContract.CreatedAt.String()),
		State:      basetypes.NewStringValue(string(entityContract.State)),
	}
}
