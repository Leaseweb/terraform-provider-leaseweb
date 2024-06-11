package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type ddos struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}

func newDdos(sdkDdos publicCloud.Ddos) ddos {
	return ddos{
		DetectionProfile: utils.GenerateString(
			sdkDdos.HasDetectionProfile(),
			sdkDdos.GetDetectionProfile(),
		),
		ProtectionType: utils.GenerateString(
			sdkDdos.HasProtectionType(),
			sdkDdos.GetProtectionType(),
		),
	}
}
