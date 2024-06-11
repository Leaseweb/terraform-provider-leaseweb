package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

func newContract(sdkContract publicCloud.Contract) contract {
	return contract{
		BillingFrequency: utils.GenerateInt(
			sdkContract.HasBillingFrequency(),
			sdkContract.GetBillingFrequency(),
		),
		Term: utils.GenerateInt(
			sdkContract.HasTerm(),
			sdkContract.GetTerm(),
		),
		Type: utils.GenerateString(
			sdkContract.HasType(),
			string(sdkContract.GetType()),
		),
		EndsAt:     utils.GenerateDateTime(sdkContract.GetEndsAt()),
		RenewalsAt: utils.GenerateDateTime(sdkContract.GetRenewalsAt()),
		CreatedAt:  utils.GenerateDateTime(sdkContract.GetCreatedAt()),
		State: utils.GenerateString(
			sdkContract.HasState(),
			string(sdkContract.GetState()),
		),
	}
}
