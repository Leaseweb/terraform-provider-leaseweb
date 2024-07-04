package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

	got, err := newHealthCheck(context.TODO(), *sdkHealthCheck)

	assert.Nil(t, err)
	assert.Equal(t, "method", got.Method.ValueString())
	assert.Equal(t, "uri", got.Uri.ValueString())
	assert.Equal(t, host, got.Host.ValueString())
	assert.Equal(t, int64(22), got.Port.ValueInt64())
}

func TestHealthCheck_attributeTypes(t *testing.T) {
	healthCheck, _ := newHealthCheck(context.TODO(), publicCloud.HealthCheck{})

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		healthCheck.AttributeTypes(),
		healthCheck,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
