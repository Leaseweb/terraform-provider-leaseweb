package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id: "imageId",
	}

	want := Image{
		Id: basetypes.NewStringValue("imageId"),
	}
	got := newImage(sdkImage)

	assert.Equal(t, want, got)
}
