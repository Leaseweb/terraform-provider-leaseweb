package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestVolume_AttributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Volume{}.AttributeTypes(),
		Volume{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
