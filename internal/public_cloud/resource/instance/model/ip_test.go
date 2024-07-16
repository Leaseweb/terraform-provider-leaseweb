package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newIp(t *testing.T) {
	reverseLookup := "reverse-lookup"

	ip := domain.NewIp(
		"1.2.3.4",
		"prefix-length",
		46,
		true,
		false,
		"tralala",
		domain.OptionalIpValues{
			Ddos:          &domain.Ddos{ProtectionType: "protection-type"},
			ReverseLookup: &reverseLookup,
		},
	)
	got, diags := newIp(context.TODO(), ip)

	assert.Nil(t, diags)
	assert.Equal(
		t,
		"1.2.3.4",
		got.Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"prefix-length",
		got.PrefixLength.ValueString(),
		"prefix-length should be set",
	)
	assert.Equal(
		t,
		int64(46),
		got.Version.ValueInt64(),
		"version should be set",
	)
	assert.Equal(
		t,
		true,
		got.NullRouted.ValueBool(),
		"nullRouted should be set",
	)
	assert.Equal(
		t,
		false,
		got.MainIp.ValueBool(),
		"mainIp should be set",
	)
	assert.Equal(
		t,
		"tralala",
		got.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"reverse-lookup",
		got.ReverseLookup.ValueString(),
		"reverseLookup should be set",
	)

	ddos := Ddos{}
	got.Ddos.As(context.TODO(), &ddos, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"protection-type",
		ddos.ProtectionType.ValueString(),
		"ddos should be set",
	)
}

func TestIp_attributeTypes(t *testing.T) {
	ip, _ := newIp(context.TODO(), domain.Ip{})

	_, diags := types.ObjectValueFrom(context.TODO(), ip.AttributeTypes(), ip)

	assert.Nil(t, diags, "attributes should be correct")
}
