package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func Test_newIp(t *testing.T) {
	ip := domain.NewIp(
		"ip",
		"prefixLength",
		46,
		true,
		false,
		enum.NetworkTypeInternal,
		domain.OptionalIpValues{
			Ddos: &domain.Ddos{ProtectionType: "protection-type"},
		},
	)

	got := newIp(ip)

	assert.Equal(
		t,
		"ip",
		got.Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"prefixLength",
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
		"INTERNAL",
		got.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"protection-type",
		got.Ddos.ProtectionType.ValueString(),
		"ddos should be set",
	)
}
