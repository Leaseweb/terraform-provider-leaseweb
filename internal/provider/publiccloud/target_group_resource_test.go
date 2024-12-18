package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_targetGroupResourceModel_generateCreateOpts(t *testing.T) {
	t.Run("required fields are set", func(t *testing.T) {
		targetGroup := targetGroupResourceModel{
			Name:     basetypes.NewStringValue("Name"),
			Protocol: basetypes.NewStringValue("HTTP"),
			Port:     basetypes.NewInt32Value(80),
			Region:   basetypes.NewStringValue("region"),
		}

		got, err := targetGroup.generateCreateOpts(context.TODO())

		want := publiccloud.CreateTargetGroupOpts{
			Name:     "Name",
			Protocol: "HTTP",
			Port:     80,
			Region:   "region",
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("optional fields are set", func(t *testing.T) {
		healthCheckObject, _ := types.ObjectValueFrom(
			context.TODO(),
			healthCheckResourceModel{}.attributeTypes(),
			healthCheckResourceModel{
				Protocol: basetypes.NewStringValue("HTTP"),
			},
		)
		targetGroup := targetGroupResourceModel{
			HealthCheck: healthCheckObject,
		}

		got, err := targetGroup.generateCreateOpts(context.TODO())

		want := publiccloud.CreateTargetGroupOpts{
			HealthCheck: &publiccloud.HealthCheckOpts{
				Protocol: publiccloud.Protocol("HTTP"),
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("invalid healthCheck returns an error", func(t *testing.T) {
		type dummy struct{}

		healthCheckObject, _ := types.ObjectValueFrom(
			context.TODO(),
			map[string]attr.Type{},
			dummy{},
		)

		targetGroup := targetGroupResourceModel{
			HealthCheck: healthCheckObject,
		}

		_, err := targetGroup.generateCreateOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, ".healthCheckResourceModel")
	})
}

func Test_adaptHealthCheckToHealthCheckResource(t *testing.T) {
	t.Run("required fields are set", func(t *testing.T) {
		sdkHealthCheck := publiccloud.HealthCheck{
			Protocol: publiccloud.PROTOCOL_HTTP,
			Uri:      "/",
			Port:     80,
		}

		got := adaptHealthCheckToHealthCheckResource(sdkHealthCheck)

		want := healthCheckResourceModel{
			Protocol: basetypes.NewStringValue("HTTP"),
			URI:      basetypes.NewStringValue("/"),
			Port:     basetypes.NewInt32Value(80),
		}

		assert.Equal(t, want, got)
	})

	t.Run("optional fields are set", func(t *testing.T) {
		httpMethod := publiccloud.HTTPMETHOD_GET
		host := "example.com"
		sdkHealthCheck := publiccloud.HealthCheck{
			Method: *publiccloud.NewNullableHttpMethod(&httpMethod),
			Host:   *publiccloud.NewNullableString(&host),
		}

		got := adaptHealthCheckToHealthCheckResource(sdkHealthCheck)

		want := healthCheckResourceModel{
			Protocol: basetypes.NewStringValue(""),
			URI:      basetypes.NewStringValue(""),
			Port:     basetypes.NewInt32Value(0),
			Method:   basetypes.NewStringValue("GET"),
			Host:     basetypes.NewStringValue("example.com"),
		}

		assert.Equal(t, want, got)
	})
}

func Test_targetGroupResourceModel_generateUpdateOpts(t *testing.T) {
	t.Run("main fields are set", func(t *testing.T) {
		targetGroup := targetGroupResourceModel{
			Name:     basetypes.NewStringValue("Name"),
			Protocol: basetypes.NewStringValue("HTTP"),
			Port:     basetypes.NewInt32Value(80),
			Region:   basetypes.NewStringValue("eu-west-2"),
		}

		got, err := targetGroup.generateUpdateOpts(context.TODO())

		name := "Name"
		port := int32(80)
		want := publiccloud.UpdateTargetGroupOpts{
			Name: &name,
			Port: &port,
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("invalid healthCheck returns an error", func(t *testing.T) {
		type dummy struct{}

		healthCheckObject, _ := types.ObjectValueFrom(
			context.TODO(),
			map[string]attr.Type{},
			dummy{},
		)

		targetGroup := targetGroupResourceModel{
			HealthCheck: healthCheckObject,
		}

		_, err := targetGroup.generateUpdateOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, ".healthCheckResourceModel")
	})

	t.Run("healthCheck is set", func(t *testing.T) {
		healthCheckObject, _ := types.ObjectValueFrom(
			context.TODO(),
			healthCheckResourceModel{}.attributeTypes(),
			healthCheckResourceModel{
				Protocol: basetypes.NewStringValue("HTTP"),
			},
		)
		targetGroup := targetGroupResourceModel{
			HealthCheck: healthCheckObject,
		}

		got, err := targetGroup.generateUpdateOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, publiccloud.PROTOCOL_HTTP, got.HealthCheck.Protocol)
	})
}

func Test_healthCheckResourceModel_generateOpts(t *testing.T) {
	t.Run("required fields are set", func(t *testing.T) {
		healthCheck := healthCheckResourceModel{
			Protocol: basetypes.NewStringValue("HTTP"),
			URI:      basetypes.NewStringValue("/"),
			Port:     basetypes.NewInt32Value(80),
		}

		got := healthCheck.generateOpts()

		protocol := publiccloud.PROTOCOL_HTTP
		uri := "/"
		port := int32(80)
		want := publiccloud.HealthCheckOpts{
			Protocol: protocol,
			Uri:      uri,
			Port:     port,
		}

		assert.Equal(t, want, got)
	})

	t.Run("optional fields are set", func(t *testing.T) {
		healthCheck := healthCheckResourceModel{
			Method: basetypes.NewStringValue("GET"),
			Host:   basetypes.NewStringValue("example.com"),
		}

		got := healthCheck.generateOpts()

		method := publiccloud.HTTPMETHODOPT_GET
		host := "example.com"
		want := publiccloud.HealthCheckOpts{
			Method: &method,
			Host:   &host,
		}

		assert.Equal(t, *want.Method, *got.Method)
		assert.Equal(t, *want.Host, *got.Host)
	})
}

func Test_adaptTargetGroupToTargetGroupResource(t *testing.T) {
	t.Run("main fields are set", func(t *testing.T) {
		sdkTargetGroup := publiccloud.TargetGroup{
			Id:       "ID",
			Name:     "Name",
			Protocol: publiccloud.PROTOCOL_HTTP,
			Port:     80,
			Region:   publiccloud.REGIONNAME_EU_CENTRAL_1,
		}

		got, err := adaptTargetGroupToTargetGroupResource(
			sdkTargetGroup,
			context.TODO(),
		)

		want := targetGroupResourceModel{
			ID:          basetypes.NewStringValue("ID"),
			Name:        basetypes.NewStringValue("Name"),
			Protocol:    basetypes.NewStringValue("HTTP"),
			Port:        basetypes.NewInt32Value(80),
			Region:      basetypes.NewStringValue("eu-central-1"),
			HealthCheck: basetypes.NewObjectNull(healthCheckResourceModel{}.attributeTypes()),
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("healthCheck is set", func(t *testing.T) {
		sdkTargetGroup := publiccloud.TargetGroup{
			HealthCheck: *publiccloud.NewNullableHealthCheck(
				&publiccloud.HealthCheck{
					Protocol: publiccloud.PROTOCOL_HTTP,
				},
			),
		}

		targetGroup, err := adaptTargetGroupToTargetGroupResource(
			sdkTargetGroup,
			context.TODO(),
		)

		got := healthCheckResourceModel{}
		targetGroup.HealthCheck.As(context.TODO(), &got, basetypes.ObjectAsOptions{})

		assert.NoError(t, err)
		assert.Equal(t, "HTTP", got.Protocol.ValueString())
	})
}
