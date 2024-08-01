package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestPrivateNetwork_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		PrivateNetwork{}.AttributeTypes(),
		PrivateNetwork{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
