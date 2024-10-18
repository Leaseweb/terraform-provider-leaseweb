package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newDataSourceModelImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id: "imageId",
	}

	want := DataSourceModelImage{
		Id: basetypes.NewStringValue("imageId"),
	}
	got := newDataSourceModelImage(sdkImage)

	assert.Equal(t, want, got)
}
