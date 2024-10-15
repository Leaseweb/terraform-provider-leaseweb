package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestNewIp(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := Ip{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got := NewIp(sdkIp)

	assert.Equal(t, want, got)
}
