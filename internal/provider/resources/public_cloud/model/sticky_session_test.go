package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestStickySession_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		StickySession{}.AttributeTypes(),
		StickySession{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
