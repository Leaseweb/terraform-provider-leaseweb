package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newFromIp(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := ResourceModelIp{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got, err := newResourceModelIpFromIp(context.TODO(), sdkIp)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_newFromIpDetails(t *testing.T) {
	sdkIpDetails := publicCloud.IpDetails{
		Ip: "127.0.0.1",
	}

	want := ResourceModelIp{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got, err := newResourceModelIpFromIpDetails(context.TODO(), sdkIpDetails)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}
