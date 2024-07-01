package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newIp(t *testing.T) {
	reverseLookup := "reverse-lookup"

	sdkIp := publicCloud.NewIp(
		"got",
		"prefix-length",
		46,
		true,
		false,
		"tralala",
		*publicCloud.NewNullableString(&reverseLookup),
		*publicCloud.NewNullableDdos(&publicCloud.Ddos{ProtectionType: "protection-type"}),
	)
	got, _ := newIp(context.TODO(), sdkIp)

	assert.Equal(
		t,
		"got",
		got.Ip.ValueString(),
		"got should be set",
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
	ip, _ := newIp(context.TODO(), &publicCloud.Ip{})

	_, diags := types.ObjectValueFrom(context.TODO(), ip.AttributeTypes(), ip)

	assert.Nil(t, diags, "attributes should be correct")
}
