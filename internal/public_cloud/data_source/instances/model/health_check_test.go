package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newHealthCheck(t *testing.T) {
	host := "host"
	healthCheck := domain.NewHealthCheck(
		"method",
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got := newHealthCheck(healthCheck)

	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}
