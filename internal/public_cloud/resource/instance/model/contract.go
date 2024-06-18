package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

func (c Contract) attributeTypes() map[string]attr.Type {
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

func newContract(sdkContract *publicCloud.Contract) Contract {
	return Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.GetBillingFrequency())),
		Term:             basetypes.NewInt64Value(int64(sdkContract.GetTerm())),
		Type:             basetypes.NewStringValue(string(sdkContract.GetType())),
		EndsAt:           basetypes.NewStringValue(sdkContract.GetEndsAt().String()),
		RenewalsAt:       basetypes.NewStringValue(sdkContract.GetRenewalsAt().String()),
		CreatedAt:        basetypes.NewStringValue(sdkContract.GetCreatedAt().String()),
		State:            basetypes.NewStringValue(string(sdkContract.GetState())),
	}
}
