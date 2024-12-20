package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource(t *testing.T) {
	sdkLoadBalancerListenerRule := publiccloud.LoadBalancerListenerRule{
		TargetGroupId: "targetGroupId",
	}

	got := adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource(sdkLoadBalancerListenerRule)
	want := loadBalancerListenerDefaultRuleResourceModel{
		TargetGroupID: basetypes.NewStringValue("targetGroupId"),
	}

	assert.Equal(t, want, got)
}

func Test_adaptLoadBalancerListenerToLoadBalancerListenerResource(t *testing.T) {
	t.Run("main values are set as expected", func(t *testing.T) {
		sdkLoadBalancerListener := publiccloud.LoadBalancerListener{
			Id:       "id",
			Protocol: publiccloud.PROTOCOL_HTTP,
			Port:     22,
		}

		got, err := adaptLoadBalancerListenerToLoadBalancerListenerResource(
			sdkLoadBalancerListener,
			context.TODO(),
		)

		want := loadBalancerListenerResourceModel{
			ListenerID: basetypes.NewStringValue("id"),
			Protocol:   basetypes.NewStringValue("HTTP"),
			Port:       basetypes.NewInt32Value(22),
		}

		require.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("first rule is set as defaultRule", func(t *testing.T) {
		sdkLoadBalancerListener := publiccloud.LoadBalancerListener{
			Rules: []publiccloud.LoadBalancerListenerRule{
				{
					TargetGroupId: "targetGroupId1",
				},
				{
					TargetGroupId: "targetGroupId2",
				},
			},
		}

		got, err := adaptLoadBalancerListenerToLoadBalancerListenerResource(
			sdkLoadBalancerListener,
			context.TODO(),
		)

		want, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			loadBalancerListenerDefaultRuleResourceModel{
				TargetGroupID: basetypes.NewStringValue("targetGroupId1"),
			},
		)

		require.NoError(t, err)
		assert.Equal(t, want, got.DefaultRule)
	})
}

func Test_loadBalancerListenerCertificateResourceModel_generateSslCertificate(t *testing.T) {
	t.Run("chain is set if model chain is set", func(t *testing.T) {
		certificate := loadBalancerListenerCertificateResourceModel{
			PrivateKey:  basetypes.NewStringValue("privateKey"),
			Certificate: basetypes.NewStringValue("certificate"),
			Chain:       basetypes.NewStringValue("chain"),
		}

		got := certificate.generateSslCertificate()

		chain := "chain"
		want := publiccloud.SslCertificate{
			PrivateKey:  "privateKey",
			Certificate: "certificate",
			Chain:       &chain,
		}

		assert.Equal(t, want, got)
	})

	t.Run("chain is set to nil if model chain is nil", func(t *testing.T) {
		certificate := loadBalancerListenerCertificateResourceModel{
			PrivateKey:  basetypes.NewStringValue("privateKey"),
			Certificate: basetypes.NewStringValue("certificate"),
			Chain:       basetypes.NewStringNull(),
		}

		got := certificate.generateSslCertificate()

		want := publiccloud.SslCertificate{
			PrivateKey:  "privateKey",
			Certificate: "certificate",
			Chain:       nil,
		}

		assert.Equal(t, want, got)
	})

	t.Run(
		"chain is set to nil if model chain is an empty string",
		func(t *testing.T) {
			certificate := loadBalancerListenerCertificateResourceModel{
				PrivateKey:  basetypes.NewStringValue("privateKey"),
				Certificate: basetypes.NewStringValue("certificate"),
				Chain:       basetypes.NewStringValue(""),
			}

			got := certificate.generateSslCertificate()

			want := publiccloud.SslCertificate{
				PrivateKey:  "privateKey",
				Certificate: "certificate",
				Chain:       nil,
			}

			assert.Equal(t, want, got)
		},
	)
}

func Test_loadBalancerListenerDefaultRuleResourceModel_generateLoadBalancerListenerDefaultRule(t *testing.T) {
	rule := loadBalancerListenerDefaultRuleResourceModel{
		TargetGroupID: basetypes.NewStringValue("targetGroupId"),
	}

	got := rule.generateLoadBalancerListenerDefaultRule()

	want := publiccloud.LoadBalancerListenerDefaultRule{
		TargetGroupId: "targetGroupId",
	}

	assert.Equal(t, want, got)
}
