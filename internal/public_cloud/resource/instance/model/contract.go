package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/utils"
)

type Contract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	RenewalsAt       types.String `tfsdk:"renewals_at"`
	CreatedAt        types.String `tfsdk:"created_at"`
	State            types.String `tfsdk:"state"`
}

func (c Contract) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int64Type,
		"term":              types.Int64Type,
		"type":              types.StringType,
		"ends_at":           types.StringType,
		"renewals_at":       types.StringType,
		"created_at":        types.StringType,
		"state":             types.StringType,
	}
}

func newContract(
	ctx context.Context,
	entityContract entity.Contract,
) (*Contract, diag.Diagnostics) {
	return &Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(entityContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(entityContract.Term)),
		Type:             basetypes.NewStringValue(string(entityContract.Type)),
		EndsAt:           utils.ConvertNullableTimeToStringValue(entityContract.EndsAt),
		RenewalsAt:       basetypes.NewStringValue(entityContract.RenewalsAt.String()),
		CreatedAt:        basetypes.NewStringValue(entityContract.CreatedAt.String()),
		State:            basetypes.NewStringValue(string(entityContract.State)),
	}, nil
}
