package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		HealthCheck{}.AttributeTypes(),
		HealthCheck{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
