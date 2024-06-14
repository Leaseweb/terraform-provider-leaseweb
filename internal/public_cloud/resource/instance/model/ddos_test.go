package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newDdos(t *testing.T) {
	sdkDdos := publicCloud.NewNullableDdos(publicCloud.NewDdos())
	sdkDdos.Get().SetDetectionProfile("detectionProfile")
	sdkDdos.Get().SetProtectionType("protectionType")

	dDos := newDdos(sdkDdos.Get())

	assert.Equal(
		t,
		"detectionProfile",
		dDos.DetectionProfile.ValueString(),
		"detectionProfile should be set",
	)
	assert.Equal(
		t,
		"protectionType",
		dDos.ProtectionType.ValueString(),
		"protectionType should be set",
	)
}

func TestDdos_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(context.TODO(), Ddos{}.attributeTypes(), Ddos{})

	assert.Nil(t, diags, "attributes should be correct")
}