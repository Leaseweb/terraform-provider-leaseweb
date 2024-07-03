package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newHealthCheck(t *testing.T) {
	host := "host"
	sdkHealthCheck := publicCloud.NewHealthCheck(
		"method",
		"uri",
		*publicCloud.NewNullableString(&host),
		22,
	)

	got := newHealthCheck(*sdkHealthCheck)

	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}
