package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Ddos struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}
