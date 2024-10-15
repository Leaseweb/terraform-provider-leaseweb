package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/dataadapters/shared"
)

type Contract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func NewContract(sdkContract publicCloud.Contract) Contract {
	return Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}
}
