package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newDataSourceModelIp(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := DataSourceModelIp{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got := newDataSourceModelIp(sdkIp)

	assert.Equal(t, want, got)
}
