package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newDdos(t *testing.T) {
	entityDdos := domain.NewDdos("detectionProfile", "protectionType")

	got, diags := newDdos(context.TODO(), entityDdos)

	assert.Nil(t, diags)

	assert.Equal(
		t,
		"detectionProfile",
		got.DetectionProfile.ValueString(),
		"detectionProfile should be set",
	)
	assert.Equal(
		t,
		"protectionType",
		got.ProtectionType.ValueString(),
		"protectionType should be set",
	)
}

func TestDdos_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Ddos{}.AttributeTypes(),
		Ddos{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
