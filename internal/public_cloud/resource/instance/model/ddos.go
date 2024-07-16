package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
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

func newDdos(
	ctx context.Context,
	entityDdos domain.Ddos,
) (*Ddos, diag.Diagnostics) {
	return &Ddos{
		DetectionProfile: basetypes.NewStringValue(entityDdos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(entityDdos.ProtectionType),
	}, nil
}
