package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/model"
)

type Reason string

const (
	ReasonContractTermCannotBeZero Reason = "contract.term cannot be 0 when contract type is MONTHLY"
	ReasonContractTermMustBeZero   Reason = "contract.term must be 0 when contract type is HOURLY"
	ReasonNone                     Reason = ""
)

type Contract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func (c Contract) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int64Type,
		"term":              types.Int64Type,
		"type":              types.StringType,
		"ends_at":           types.StringType,
		"state":             types.StringType,
	}
}

func (c Contract) IsContractTermValid() (bool, Reason) {
	if c.Type.ValueString() == string(publicCloud.CONTRACTTYPE_MONTHLY) && c.Term.ValueInt64() == 0 {
		return false, ReasonContractTermCannotBeZero
	}

	if c.Type.ValueString() == string(publicCloud.CONTRACTTYPE_HOURLY) && c.Term.ValueInt64() != 0 {
		return false, ReasonContractTermMustBeZero
	}

	return true, ReasonNone
}

func newContract(ctx context.Context, sdkContract publicCloud.Contract) (*Contract, error) {
	return &Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           model.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}, nil
}
