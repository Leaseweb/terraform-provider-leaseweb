package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type StickySession struct {
	Enabled     types.Bool  `tfsdk:"enabled"`
	MaxLifeTime types.Int64 `tfsdk:"max_life_time"`
}

func (s StickySession) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":       types.BoolType,
		"max_life_time": types.Int64Type,
	}
}

func newStickySession(
	ctx context.Context,
	sdkStickySession publicCloud.StickySession,
) (*StickySession, diag.Diagnostics) {
	return &StickySession{
		Enabled:     basetypes.NewBoolValue(sdkStickySession.GetEnabled()),
		MaxLifeTime: basetypes.NewInt64Value(int64(sdkStickySession.GetMaxLifeTime())),
	}, nil
}
