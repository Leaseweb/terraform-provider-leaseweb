package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			ID:       basetypes.NewStringValue("ID"),
			Name:     basetypes.NewStringValue("Name"),
			Protocol: basetypes.NewStringValue("HTTP"),
			Port:     basetypes.NewInt32Value(80),
			Region:   basetypes.NewStringValue("eu-central-1"),
			HealthCheck: basetypes.NewObjectNull(
				map[string]attr.Type{
					"protocol": types.StringType,
					"method":   types.StringType,
					"uri":      types.StringType,
					"host":     types.StringType,
					"port":     types.Int32Type,
				},
			),
		}

		require.NoError(t, err)
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

		require.NoError(t, err)
		assert.Equal(t, "HTTP", got.Protocol.ValueString())
	})
}
