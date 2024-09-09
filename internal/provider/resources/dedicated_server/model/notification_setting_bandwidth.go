package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NotificationSettingBandwidth struct {
	ServerId            types.String `tfsdk:"server_id"`
	Id                  types.String `tfsdk:"id"`
	Frequency           types.String `tfsdk:"frequency"`
	LastCheckedAt       types.String `tfsdk:"last_checked_at"`
	Threshold           types.String `tfsdk:"threshold"`
	ThresholdExceededAt types.String `tfsdk:"threshold_exceeded_at"`
	Unit                types.String `tfsdk:"unit"`
	Actions             types.List   `tfsdk:"actions"`
}

func (n NotificationSettingBandwidth) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_id":             types.StringType,
		"id":                    types.StringType,
		"frequency":             types.StringType,
		"last_checked_at":       types.StringType,
		"threshold":             types.StringType,
		"threshold_exceeded_at": types.StringType,
		"unit":                  types.StringType,
		"actions":               types.ListType{ElemType: types.ObjectType{AttrTypes: Action{}.AttributeTypes()}},
	}
}
