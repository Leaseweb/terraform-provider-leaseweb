package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Ddos struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}

func (d Ddos) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"detection_profile": types.StringType,
		"protection_type":   types.StringType,
	}
}
