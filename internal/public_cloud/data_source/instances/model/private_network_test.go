package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newPrivateNetwork(t *testing.T) {
	sdkPrivateNetwork := publicCloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)
	got := newPrivateNetwork(*sdkPrivateNetwork)

	assert.Equal(
		t,
		"id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"status",
		got.Status.ValueString(),
		"status should be set",
	)
	assert.Equal(
		t,
		"subnet",
		got.Subnet.ValueString(),
		"subnet should be set",
	)
}
