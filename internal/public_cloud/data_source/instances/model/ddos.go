package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type ddos struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}

func newDdos(sdkDdos publicCloud.Ddos) ddos {
	return ddos{
		DetectionProfile: basetypes.NewStringValue(sdkDdos.GetDetectionProfile()),
		ProtectionType:   basetypes.NewStringValue(sdkDdos.GetProtectionType()),
	}
}
