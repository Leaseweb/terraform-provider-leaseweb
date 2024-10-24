package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptSdkImageToDatasourceImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id:      "imageId",
		Name:    "name",
		Custom:  true,
		Flavour: "flavour",
	}

	want := dataSourceModelImage{
		ID:      basetypes.NewStringValue("imageId"),
		Name:    basetypes.NewStringValue("name"),
		Custom:  basetypes.NewBoolValue(true),
		Flavour: basetypes.NewStringValue("flavour"),
	}
	got := adaptSdkImageToDatasourceImage(sdkImage)

	assert.Equal(t, want, got)
}

func Test_adaptSdkImageDetailsToDatasourceImage(t *testing.T) {
	state := publicCloud.IMAGESTATE_READY
	region := publicCloud.REGIONNAME_EU_WEST_3

	sdkImageDetails := publicCloud.ImageDetails{
		Id:           "imageId",
		Name:         "name",
		Custom:       true,
		State:        *publicCloud.NewNullableImageState(&state),
		MarketApps:   []publicCloud.MarketAppId{publicCloud.MARKETAPPID_CPANEL_30},
		StorageTypes: []publicCloud.StorageType{publicCloud.STORAGETYPE_CENTRAL},
		Flavour:      "flavour",
		Region:       *publicCloud.NewNullableRegionName(&region),
	}

	want := dataSourceModelImage{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("READY"),
		MarketApps:   []string{"CPANEL_30"},
		StorageTypes: []string{"CENTRAL"},
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}
	got := adaptSdkImageDetailsToDatasourceImage(sdkImageDetails)

	assert.Equal(t, want, got)
}

func Test_adaptSdkImagesToDatasourceImages(t *testing.T) {
	sdkImages := []publicCloud.ImageDetails{
		{Id: "id"},
	}

	got := adaptSdkImagesToDatasourceImages(sdkImages)

	assert.Len(t, got.Images, 1)
	assert.Equal(t, "id", got.Images[0].ID.ValueString())
}
