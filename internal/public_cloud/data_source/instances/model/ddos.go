package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type ddos struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}

func newDdos(entityDdos entity.Ddos) *ddos {
	return &ddos{
		DetectionProfile: basetypes.NewStringValue(entityDdos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(entityDdos.ProtectionType),
	}
}
