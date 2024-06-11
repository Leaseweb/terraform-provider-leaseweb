package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newDDos(t *testing.T) {
	sdkDdos := publicCloud.NewDdos()
	sdkDdos.SetDetectionProfile("detectionProfile")
	sdkDdos.SetProtectionType("protectionType")

	dDos := newDdos(*sdkDdos)

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
