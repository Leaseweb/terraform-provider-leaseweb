package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id: "imageId",
	}

	want := Image{
		Id: basetypes.NewStringValue("imageId"),
	}
	got := NewImage(sdkImage)

	assert.Equal(t, want, got)
}
