package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newDDos(t *testing.T) {
	sdkDdos := publicCloud.Ddos{
		DetectionProfile: "detectionProfile",
		ProtectionType:   "protectionType",
	}
	got := newDdos(sdkDdos)

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
