package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
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
		BillingFrequency: resources.GetIntValue(sdkContract.HasBillingFrequency(), sdkContract.GetBillingFrequency()),
		Term:             resources.GetIntValue(sdkContract.HasTerm(), sdkContract.GetTerm()),
		Type:             resources.GetStringValue(sdkContract.HasType(), string(sdkContract.GetType())),
		EndsAt:           resources.GetDateTime(sdkContract.GetEndsAt()),
		RenewalsAt:       resources.GetDateTime(sdkContract.GetRenewalsAt()),
		CreatedAt:        resources.GetDateTime(sdkContract.GetCreatedAt()),
		State:            resources.GetStringValue(sdkContract.HasState(), string(sdkContract.GetState())),
	}
}
