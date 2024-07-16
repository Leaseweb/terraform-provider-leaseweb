package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newStickySession(t *testing.T) {
	stickySession := domain.NewStickySession(false, 1)

	got := newStickySession(stickySession)

	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}
