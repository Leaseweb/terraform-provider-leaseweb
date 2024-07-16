package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type stickySession struct {
	Enabled     types.Bool  `tfsdk:"enabled"`
	MaxLifeTime types.Int64 `tfsdk:"max_life_time"`
}

func newStickySession(entityStickySession domain.StickySession) *stickySession {
	return &stickySession{
		Enabled:     basetypes.NewBoolValue(entityStickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(entityStickySession.MaxLifeTime)),
	}
}
