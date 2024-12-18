package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptImageToImageDataSource(t *testing.T) {
	sdkImage := publiccloud.Image{
		Id:      "imageId",
		Name:    "name",
		Custom:  true,
		Flavour: "flavour",
	}

	want := imageModelDataSource{
		ID:      basetypes.NewStringValue("imageId"),
		Name:    basetypes.NewStringValue("name"),
		Custom:  basetypes.NewBoolValue(true),
		Flavour: basetypes.NewStringValue("flavour"),
	}
	got := adaptImageToImageDataSource(sdkImage)

	assert.Equal(t, want, got)
}

func Test_adaptImageDetailsToImageDataSource(t *testing.T) {
	state := publiccloud.IMAGESTATE_READY
	region := publiccloud.REGIONNAME_EU_WEST_3

	sdkImageDetails := publiccloud.ImageDetails{
		Id:           "imageId",
		Name:         "name",
		Custom:       true,
		State:        *publiccloud.NewNullableImageState(&state),
		MarketApps:   []publiccloud.MarketAppId{publiccloud.MARKETAPPID_CPANEL_30},
		StorageTypes: []publiccloud.StorageType{publiccloud.STORAGETYPE_CENTRAL},
		Flavour:      "flavour",
		Region:       *publiccloud.NewNullableRegionName(&region),
	}

	want := imageModelDataSource{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("READY"),
		MarketApps:   []string{"CPANEL_30"},
		StorageTypes: []string{"CENTRAL"},
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}
	got := adaptImageDetailsToImageDataSource(sdkImageDetails)

	assert.Equal(t, want, got)
}

func Test_adaptImagesToImagesDataSource(t *testing.T) {
	sdkImages := []publiccloud.ImageDetails{
		{Id: "id"},
	}

	got := adaptImagesToImagesDataSource(sdkImages)

	assert.Len(t, got.Images, 1)
	assert.Equal(t, "id", got.Images[0].ID.ValueString())
}
