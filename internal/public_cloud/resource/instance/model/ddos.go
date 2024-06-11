package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type Ddos struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}

func (d Ddos) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"detection_profile": types.StringType,
		"protection_type":   types.StringType,
	}
}

func newDdos(sdkDDos *publicCloud.Ddos) Ddos {
	return Ddos{
		DetectionProfile: utils.GenerateString(sdkDDos.HasDetectionProfile(), sdkDDos.GetDetectionProfile()),
		ProtectionType:   utils.GenerateString(sdkDDos.HasProtectionType(), sdkDDos.GetProtectionType()),
	}
}
