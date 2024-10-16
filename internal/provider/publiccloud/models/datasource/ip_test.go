package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newIp(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := Ip{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got := newIp(sdkIp)

	assert.Equal(t, want, got)
}
