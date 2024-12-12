package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptIpToIPResource(t *testing.T) {
	t.Run("is set properly when reverseLookup is set", func(t *testing.T) {
		reverseLookup := "example.com"
		sdkIp := publiccloud.Ip{
			Ip:            "127.0.0.1",
			ReverseLookup: *publiccloud.NewNullableString(&reverseLookup),
		}

		want := ipResourceModel{
			IP:            basetypes.NewStringValue("127.0.0.1"),
			ReverseLookup: basetypes.NewStringPointerValue(&reverseLookup),
		}
		got := adaptIpToIPResource(sdkIp)

		assert.Equal(t, want, got)
	})

	t.Run("is set properly when reverseLookup is null", func(t *testing.T) {
		sdkIp := publiccloud.Ip{
			Ip:            "127.0.0.1",
			ReverseLookup: *publiccloud.NewNullableString(nil),
		}

		want := ipResourceModel{
			IP:            basetypes.NewStringValue("127.0.0.1"),
			ReverseLookup: basetypes.NewStringPointerValue(nil),
		}
		got := adaptIpToIPResource(sdkIp)

		assert.Equal(t, want, got)
	})
}

func Test_adaptIpDetailsToIPResource(t *testing.T) {
	t.Run("is set properly when reverseLookup is set", func(t *testing.T) {
		reverseLookup := "example.com"
		sdkIpDetails := publiccloud.IpDetails{
			Ip:            "127.0.0.1",
			ReverseLookup: *publiccloud.NewNullableString(&reverseLookup),
		}

		want := ipResourceModel{
			IP:            basetypes.NewStringValue("127.0.0.1"),
			ReverseLookup: basetypes.NewStringPointerValue(&reverseLookup),
		}
		got := adaptIpDetailsToIPResource(sdkIpDetails)

		assert.Equal(t, want, got)
	})

	t.Run("is set properly when reverseLookup is null", func(t *testing.T) {
		sdkIpDetails := publiccloud.IpDetails{
			Ip:            "127.0.0.1",
			ReverseLookup: *publiccloud.NewNullableString(nil),
		}

		want := ipResourceModel{
			IP:            basetypes.NewStringValue("127.0.0.1"),
			ReverseLookup: basetypes.NewStringPointerValue(nil),
		}
		got := adaptIpDetailsToIPResource(sdkIpDetails)

		assert.Equal(t, want, got)
	})
}

func Test_iPResourceModel_generateUpdateOpts(t *testing.T) {
	t.Run("is set properly when reverseLookup is set", func(t *testing.T) {
		reverseLookup := "example.com"
		ip := ipResourceModel{
			ReverseLookup: basetypes.NewStringPointerValue(&reverseLookup),
		}
		got := ip.generateUpdateOpts()

		want := publiccloud.UpdateIpOpts{
			ReverseLookup: "example.com",
		}

		assert.Equal(t, want, got)
	})

	t.Run("is set properly when reverseLookup is not set", func(t *testing.T) {
		ip := ipResourceModel{
			ReverseLookup: basetypes.NewStringPointerValue(nil),
		}
		got := ip.generateUpdateOpts()

		want := publiccloud.UpdateIpOpts{
			ReverseLookup: "",
		}

		assert.Equal(t, want, got)
	})

}
