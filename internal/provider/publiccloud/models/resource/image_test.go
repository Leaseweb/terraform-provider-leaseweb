package resource

import (
	"context"
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
	got, err := newImage(context.TODO(), sdkImage)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}
