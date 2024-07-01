package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newIp(t *testing.T) {

	sdkDdos := publicCloud.NewNullableDdos(&publicCloud.Ddos{})
	sdkDdos.Get().SetProtectionType("protection-type")

	sdkIp := publicCloud.Ip{}
	sdkIp.SetIp("ip")
	sdkIp.SetPrefixLength("prefix-length")
	sdkIp.SetVersion(46)
	sdkIp.SetNullRouted(true)
	sdkIp.SetMainIp(false)
	sdkIp.SetNetworkType("tralala")
	sdkIp.SetReverseLookup("reverse-lookup")
	sdkIp.Ddos = *sdkDdos

	ip := newIp(sdkIp)

	assert.Equal(
		t,
		"ip",
		ip.Ip.ValueString(),
		"ip should be set",
	)
	assert.Equal(
		t,
		"prefix-length",
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
		"tralala",
		ip.NetworkType.ValueString(),
		"networkType should be set",
	)
	assert.Equal(
		t,
		"reverse-lookup",
		ip.ReverseLookup.ValueString(),
		"reverseLookup should be set",
	)
	assert.Equal(
		t,
		"protection-type",
		ip.Ddos.ProtectionType.ValueString(),
		"ddos should be set",
	)
}
