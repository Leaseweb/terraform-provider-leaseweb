package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func Test_newIp(t *testing.T) {
	entityIp := entity.NewIp(
		"ip",
		"prefixLength",
		46,
		true,
		false,
		enum.NetworkTypeInternal,
		entity.OptionalIpValues{
			Ddos: &entity.Ddos{ProtectionType: "protection-type"},
		},
	)

	ip := newIp(entityIp)

	assert.Equal(
		t,
		"ip",
		ip.Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"prefixLength",
		ip.PrefixLength.ValueString(),
		"prefix-length should be set",
	)
	assert.Equal(
		t,
		int64(46),
		ip.Version.ValueInt64(),
		"version should be set",
	)
	assert.Equal(
		t,
		true,
		ip.NullRouted.ValueBool(),
		"nullRouted should be set",
	)
	assert.Equal(
		t,
		false,
		ip.MainIp.ValueBool(),
		"mainIp should be set",
	)
	assert.Equal(
		t,
		"INTERNAL",
		ip.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"protection-type",
		ip.Ddos.ProtectionType.ValueString(),
		"ddos should be set",
	)
}
