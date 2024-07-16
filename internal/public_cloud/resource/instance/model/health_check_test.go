package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newHealthCheck(t *testing.T) {
	host := "host"
	entityHealthCheck := domain.NewHealthCheck(
		"method",
		"uri",
		22,
		domain.OptionalHealthCheckValues{Host: &host},
	)

	got, err := newHealthCheck(context.TODO(), entityHealthCheck)

	assert.Nil(t, err)
	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func TestHealthCheck_attributeTypes(t *testing.T) {
	healthCheck, _ := newHealthCheck(context.TODO(), domain.HealthCheck{})

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		healthCheck.AttributeTypes(),
		healthCheck,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
