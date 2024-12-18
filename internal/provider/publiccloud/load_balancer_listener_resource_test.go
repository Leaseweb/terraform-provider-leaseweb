package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptSslCertificateToLoadBalancerListenerCertificateResource(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		sdkSslCertificate := publiccloud.SslCertificate{
			PrivateKey:  "privateKey",
			Certificate: "certificate",
		}

		want := loadBalancerListenerCertificateResourceModel{
			PrivateKey:  basetypes.NewStringValue("privateKey"),
			Certificate: basetypes.NewStringValue("certificate"),
		}

		got := adaptSslCertificateToLoadBalancerListenerCertificateResource(sdkSslCertificate)

		assert.Equal(t, want, got)
	})

	t.Run("chain is set if it's not an empty string", func(t *testing.T) {
		want := "chain"
		sdkSslCertificate := publiccloud.SslCertificate{
			Chain: &want,
		}

		got := adaptSslCertificateToLoadBalancerListenerCertificateResource(sdkSslCertificate)

		assert.Equal(t, want, got.Chain.ValueString())
	})

	t.Run("chain is not set if it's an empty string", func(t *testing.T) {
		sdkChain := ""
		sdkSslCertificate := publiccloud.SslCertificate{
			Chain: &sdkChain,
		}

		got := adaptSslCertificateToLoadBalancerListenerCertificateResource(sdkSslCertificate)

		assert.Nil(t, got.Chain.ValueStringPointer())
	})

}

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

func Test_adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(t *testing.T) {
	t.Run("main values are set as expected", func(t *testing.T) {
		sdkLoadBalancerListenerDetails := publiccloud.LoadBalancerListenerDetails{
			Id:       "id",
			Protocol: publiccloud.PROTOCOL_HTTP,
			Port:     22,
		}

		got, err := adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(
			sdkLoadBalancerListenerDetails,
			context.TODO(),
		)

		want := loadBalancerListenerResourceModel{
			ListenerID: basetypes.NewStringValue("id"),
			Protocol:   basetypes.NewStringValue("HTTP"),
			Port:       basetypes.NewInt32Value(22),
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("first sslCertificate is set as certificate", func(t *testing.T) {
		sdkLoadBalancerListenerDetails := publiccloud.LoadBalancerListenerDetails{
			SslCertificates: []publiccloud.SslCertificate{
				{
					PrivateKey:  "privateKey1",
					Certificate: "certificate1",
				},
				{

					PrivateKey: "privateKey2",
				},
			},
		}

		got, err := adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(
			sdkLoadBalancerListenerDetails,
			context.TODO(),
		)

		want, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			loadBalancerListenerCertificateResourceModel{
				PrivateKey:  basetypes.NewStringValue("privateKey1"),
				Certificate: basetypes.NewStringValue("certificate1"),
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, want, got.Certificate)
	})

	t.Run("first rule is set as defaultRule", func(t *testing.T) {
		sdkLoadBalancerListenerDetails := publiccloud.LoadBalancerListenerDetails{
			Rules: []publiccloud.LoadBalancerListenerRule{
				{
					TargetGroupId: "targetGroupId1",
				},
				{
					TargetGroupId: "targetGroupId2",
				},
			},
		}

		got, err := adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(
			sdkLoadBalancerListenerDetails,
			context.TODO(),
		)

		want, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			loadBalancerListenerDefaultRuleResourceModel{
				TargetGroupID: basetypes.NewStringValue("targetGroupId1"),
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, want, got.DefaultRule)
	})
}

func Test_loadBalancerListenerResourceModel_generateLoadBalancerListenerCreateOpts(t *testing.T) {
	t.Run("required fields are set", func(t *testing.T) {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{
			TargetGroupID: basetypes.NewStringValue("targetGroupId"),
		}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			defaultRule,
		)

		listener := loadBalancerListenerResourceModel{
			Protocol:    basetypes.NewStringValue("HTTPS"),
			Port:        basetypes.NewInt32Value(22),
			DefaultRule: defaultRuleObject,
		}

		got, err := listener.generateLoadBalancerListenerCreateOpts(context.TODO())

		want := publiccloud.NewLoadBalancerListenerCreateOpts(
			publiccloud.PROTOCOL_HTTPS,
			22,
			*publiccloud.NewLoadBalancerListenerDefaultRule("targetGroupId"),
		)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("invalid defaultRule returns an error", func(t *testing.T) {
		defaultRule := loadBalancerListenerCertificateResourceModel{}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			defaultRule,
		)

		listener := loadBalancerListenerResourceModel{
			Protocol:    basetypes.NewStringValue("HTTPS"),
			Port:        basetypes.NewInt32Value(22),
			DefaultRule: defaultRuleObject,
		}

		_, err := listener.generateLoadBalancerListenerCreateOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Object defines fields not found in struct")
	})

	t.Run("optional fields are set", func(t *testing.T) {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{
			TargetGroupID: basetypes.NewStringValue("targetGroupId"),
		}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			defaultRule,
		)

		certificate := loadBalancerListenerCertificateResourceModel{
			PrivateKey: basetypes.NewStringValue("privateKey"),
		}
		certificateObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			certificate,
		)

		listener := loadBalancerListenerResourceModel{
			Protocol:    basetypes.NewStringValue("HTTPS"),
			Port:        basetypes.NewInt32Value(22),
			DefaultRule: defaultRuleObject,
			Certificate: certificateObject,
		}

		got, err := listener.generateLoadBalancerListenerCreateOpts(context.TODO())

		want := publiccloud.NewLoadBalancerListenerCreateOpts(
			publiccloud.PROTOCOL_HTTPS,
			22,
			*publiccloud.NewLoadBalancerListenerDefaultRule("targetGroupId"),
		)
		want.SetCertificate(publiccloud.SslCertificate{PrivateKey: "privateKey"})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("invalid certificate returns an error", func(t *testing.T) {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{
			TargetGroupID: basetypes.NewStringValue("targetGroupId"),
		}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			defaultRule,
		)

		certificate := loadBalancerListenerDefaultRuleResourceModel{}
		certificateObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			certificate,
		)

		listener := loadBalancerListenerResourceModel{
			Protocol:    basetypes.NewStringValue("HTTPS"),
			Port:        basetypes.NewInt32Value(22),
			DefaultRule: defaultRuleObject,
			Certificate: certificateObject,
		}

		_, err := listener.generateLoadBalancerListenerCreateOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Struct defines fields not found in object")
	})
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

		assert.NoError(t, err)
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

		assert.NoError(t, err)
		assert.Equal(t, want, got.DefaultRule)
	})
}

func Test_loadBalancerListenerResourceModel_generateLoadBalancerListenerUpdateOpts(t *testing.T) {
	t.Run("optional fields are set", func(t *testing.T) {
		listener := loadBalancerListenerResourceModel{
			Protocol: basetypes.NewStringValue("HTTPS"),
			Port:     basetypes.NewInt32Value(22),
		}

		got, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		protocol := publiccloud.PROTOCOL_HTTPS
		port := int32(22)
		want := publiccloud.LoadBalancerListenerOpts{
			Protocol: &protocol,
			Port:     &port,
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("certificate.chain is passed if it's set", func(t *testing.T) {
		want := "chain"
		certificate := loadBalancerListenerCertificateResourceModel{
			Chain: basetypes.NewStringPointerValue(&want),
		}
		certificateObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			certificate,
		)

		listener := loadBalancerListenerResourceModel{
			Certificate: certificateObject,
		}

		got, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, want, *got.Certificate.Chain)
	})

	t.Run("optional defaultRule fields are set", func(t *testing.T) {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{
			TargetGroupID: basetypes.NewStringValue("targetGroupId"),
		}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			defaultRule,
		)

		listener := loadBalancerListenerResourceModel{
			DefaultRule: defaultRuleObject,
		}

		got, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		want := publiccloud.LoadBalancerListenerDefaultRule{
			TargetGroupId: "targetGroupId",
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got.DefaultRule)
	})

	t.Run("invalid certificate returns an error", func(t *testing.T) {
		certificate := loadBalancerListenerDefaultRuleResourceModel{}
		certificateObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			certificate,
		)

		listener := loadBalancerListenerResourceModel{
			Certificate: certificateObject,
		}

		_, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Struct defines fields not found in object")
	})

	t.Run("invalid defaultRule returns an error", func(t *testing.T) {
		defaultRule := loadBalancerListenerCertificateResourceModel{}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			defaultRule,
		)

		listener := loadBalancerListenerResourceModel{
			DefaultRule: defaultRuleObject,
		}

		_, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Struct defines fields not found in object")
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

		assert.Equal(t, got, want)
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

		assert.Equal(t, got, want)
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

			assert.Equal(t, got, want)
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

	assert.Equal(t, got, want)
}
