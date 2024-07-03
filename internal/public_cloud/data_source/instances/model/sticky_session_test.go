package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newStickySession(t *testing.T) {
	sdkStickySession := publicCloud.NewStickySession(false, 1)

	got := newStickySession(*sdkStickySession)

	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}
