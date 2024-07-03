package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type stickySession struct {
	Enabled     types.Bool  `tfsdk:"enabled"`
	MaxLifeTime types.Int64 `tfsdk:"max_life_time"`
}

func newStickySession(sdkStickySession publicCloud.StickySession) *stickySession {
	return &stickySession{
		Enabled:     basetypes.NewBoolValue(sdkStickySession.GetEnabled()),
		MaxLifeTime: basetypes.NewInt64Value(int64(sdkStickySession.GetMaxLifeTime())),
	}
}
