package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newStickySession(t *testing.T) {
	sdkStickySession := publicCloud.NewStickySession(false, 1)

	got, err := newStickySession(context.TODO(), *sdkStickySession)

	assert.Nil(t, err)
	assert.False(t, got.Enabled.ValueBool())
	assert.Equal(t, int64(1), got.MaxLifeTime.ValueInt64())
}

func TestStickySession_attributeTypes(t *testing.T) {
	stickySession, _ := newStickySession(context.TODO(), publicCloud.StickySession{})

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		stickySession.AttributeTypes(),
		stickySession,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
