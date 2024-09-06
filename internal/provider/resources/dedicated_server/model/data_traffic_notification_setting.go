package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataTrafficNotificationSetting struct {
	Id        types.String `tfsdk:"id"`
	ServerId  types.String `tfsdk:"server_id"`
	Frequency types.String `tfsdk:"frequency"`
	Threshold types.String `tfsdk:"threshold"`
	Unit      types.String `tfsdk:"unit"`
}

func (d DataTrafficNotificationSetting) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"server_id": types.StringType,
		"frequency": types.StringType,
		"threshold": types.StringType,
		"unit":      types.StringType,
	}
}
