package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newMemory(t *testing.T) {
	sdkMemory := publicCloud.NewMemory(1, "unit")

	got := newMemory(*sdkMemory)

	assert.Equal(
		t,
		float64(1),
		got.Value.ValueFloat64(),
		"value should be set",
	)
	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
}
