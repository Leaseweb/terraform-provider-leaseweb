package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func Test_newStickySession(t *testing.T) {
	entityStickySession := entity.NewStickySession(false, 1)

	got := newStickySession(entityStickySession)

	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}
