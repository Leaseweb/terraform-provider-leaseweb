package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func TestNewIp(t *testing.T) {
	t.Run("test required values", func(t *testing.T) {

		ip := NewIp(
			"ip",
			"prefixLength",
			1,
			false,
			true,
			enum.NetworkTypePublic,
			OptionalIpValues{},
		)

		assert.Equal(t, "ip", ip.Ip)
		assert.Equal(t, "prefixLength", ip.PrefixLength)
		assert.Equal(t, 1, ip.Version)
		assert.False(t, ip.NullRouted)
		assert.True(t, ip.MainIp)
		assert.Equal(t, enum.NetworkTypePublic, ip.NetworkType)

		assert.Nil(t, ip.Ddos)
		assert.Nil(t, ip.ReverseLookup)
	})

	t.Run("test optional values", func(t *testing.T) {
		reverseLookup := "reverseLookup"

		ip := NewIp(
			"",
			"",
			0,
			false,
			false,
			enum.NetworkTypePublic,
			OptionalIpValues{Ddos: &Ddos{}, ReverseLookup: &reverseLookup},
		)

		assert.Equal(t, Ddos{}, *ip.Ddos)
		assert.Equal(t, "reverseLookup", *ip.ReverseLookup)
	})
}
