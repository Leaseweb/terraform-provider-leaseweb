package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptSslCertificateToLoadBalancerListenerCertificateResource(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		sdkSslCertificate := publicCloud.SslCertificate{
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
		sdkSslCertificate := publicCloud.SslCertificate{
			Chain: &want,
		}

		got := adaptSslCertificateToLoadBalancerListenerCertificateResource(sdkSslCertificate)

		assert.Equal(t, want, got.Chain.ValueString())
	})

	t.Run("chain is not set if it's an empty string", func(t *testing.T) {
		sdkChain := ""
		sdkSslCertificate := publicCloud.SslCertificate{
			Chain: &sdkChain,
		}

		got := adaptSslCertificateToLoadBalancerListenerCertificateResource(sdkSslCertificate)

		assert.Nil(t, got.Chain.ValueStringPointer())
	})

}

func Test_adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource(t *testing.T) {
	sdkLoadBalancerListenerRule := publicCloud.LoadBalancerListenerRule{
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
		sdkLoadBalancerListenerDetails := publicCloud.LoadBalancerListenerDetails{
			Id:       "id",
			Protocol: publicCloud.PROTOCOL_HTTP,
			Port:     22,
		}

		got, err := adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(
			sdkLoadBalancerListenerDetails,
			context.TODO(),
		)

		want := LoadBalancerListenerResourceModel{
			ListenerID: basetypes.NewStringValue("id"),
			Protocol:   basetypes.NewStringValue("HTTP"),
			Port:       basetypes.NewInt32Value(22),
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("first sslCertificate is set as certificate", func(t *testing.T) {
		sdkLoadBalancerListenerDetails := publicCloud.LoadBalancerListenerDetails{
			SslCertificates: []publicCloud.SslCertificate{
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
		sdkLoadBalancerListenerDetails := publicCloud.LoadBalancerListenerDetails{
			Rules: []publicCloud.LoadBalancerListenerRule{
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

func TestLoadBalancerListenerResourceModel_generateLoadBalancerListenerCreateOpts(t *testing.T) {
	t.Run("required fields are set", func(t *testing.T) {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{
			TargetGroupID: basetypes.NewStringValue("targetGroupId"),
		}
		defaultRuleObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			defaultRule,
		)

		listener := LoadBalancerListenerResourceModel{
			Protocol:    basetypes.NewStringValue("HTTPS"),
			Port:        basetypes.NewInt32Value(22),
			DefaultRule: defaultRuleObject,
		}

		got, err := listener.generateLoadBalancerListenerCreateOpts(context.TODO())

		want := publicCloud.NewLoadBalancerListenerCreateOpts(
			publicCloud.PROTOCOL_HTTPS,
			22,
			*publicCloud.NewLoadBalancerListenerDefaultRule("targetGroupId"),
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

		listener := LoadBalancerListenerResourceModel{
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

		listener := LoadBalancerListenerResourceModel{
			Protocol:    basetypes.NewStringValue("HTTPS"),
			Port:        basetypes.NewInt32Value(22),
			DefaultRule: defaultRuleObject,
			Certificate: certificateObject,
		}

		got, err := listener.generateLoadBalancerListenerCreateOpts(context.TODO())

		want := publicCloud.NewLoadBalancerListenerCreateOpts(
			publicCloud.PROTOCOL_HTTPS,
			22,
			*publicCloud.NewLoadBalancerListenerDefaultRule("targetGroupId"),
		)
		want.SetCertificate(publicCloud.SslCertificate{PrivateKey: "privateKey"})

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

		listener := LoadBalancerListenerResourceModel{
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
		sdkLoadBalancerListener := publicCloud.LoadBalancerListener{
			Id:       "id",
			Protocol: publicCloud.PROTOCOL_HTTP,
			Port:     22,
		}

		got, err := adaptLoadBalancerListenerToLoadBalancerListenerResource(
			sdkLoadBalancerListener,
			context.TODO(),
		)

		want := LoadBalancerListenerResourceModel{
			ListenerID: basetypes.NewStringValue("id"),
			Protocol:   basetypes.NewStringValue("HTTP"),
			Port:       basetypes.NewInt32Value(22),
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("first rule is set as defaultRule", func(t *testing.T) {
		sdkLoadBalancerListener := publicCloud.LoadBalancerListener{
			Rules: []publicCloud.LoadBalancerListenerRule{
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

func TestLoadBalancerListenerResourceModel_generateLoadBalancerListenerUpdateOpts(t *testing.T) {
	t.Run("optional fields are set", func(t *testing.T) {
		listener := LoadBalancerListenerResourceModel{
			Protocol: basetypes.NewStringValue("HTTPS"),
			Port:     basetypes.NewInt32Value(22),
		}

		got, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		protocol := publicCloud.PROTOCOL_HTTPS
		port := int32(22)
		want := publicCloud.LoadBalancerListenerOpts{
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

		listener := LoadBalancerListenerResourceModel{
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

		listener := LoadBalancerListenerResourceModel{
			DefaultRule: defaultRuleObject,
		}

		got, err := listener.generateLoadBalancerListenerUpdateOpts(context.TODO())

		want := publicCloud.LoadBalancerListenerDefaultRule{
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

		listener := LoadBalancerListenerResourceModel{
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

		listener := LoadBalancerListenerResourceModel{
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
		want := publicCloud.SslCertificate{
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

		want := publicCloud.SslCertificate{
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

			want := publicCloud.SslCertificate{
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

	want := publicCloud.LoadBalancerListenerDefaultRule{
		TargetGroupId: "targetGroupId",
	}

	assert.Equal(t, got, want)
}
