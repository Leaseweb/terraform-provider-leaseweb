package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/stretchr/testify/assert"
)

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
