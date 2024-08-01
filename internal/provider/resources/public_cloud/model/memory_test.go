package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestMemory_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Memory{}.AttributeTypes(),
		Memory{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
