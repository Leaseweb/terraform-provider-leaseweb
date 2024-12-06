package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptImageToImageResource(t *testing.T) {
	sdkImage := publiccloud.Image{
		Id:      "imageId",
		Name:    "name",
		Custom:  true,
		Flavour: "flavour",
	}

	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})
	want := imageResourceModel{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		Flavour:      basetypes.NewStringValue("flavour"),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}
	got := adaptImageToImageResource(sdkImage)

	assert.Equal(t, want, got)
}

func Test_adaptImageDetailsToImageResource(t *testing.T) {
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

	marketApps, _ := basetypes.NewListValueFrom(
		context.TODO(),
		types.StringType,
		[]string{"CPANEL_30"},
	)
	storageTypes, _ := basetypes.NewListValueFrom(
		context.TODO(),
		types.StringType,
		[]string{"CENTRAL"},
	)

	want := imageResourceModel{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("READY"),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}
	got, err := adaptImageDetailsToImageResource(
		context.TODO(),
		sdkImageDetails,
	)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_imageResourceModel_getUpdateImageOpts(t *testing.T) {
	image := imageResourceModel{
		Name: basetypes.NewStringValue("name"),
	}
	got := image.getUpdateImageOpts()

	want := publiccloud.UpdateImageOpts{Name: "name"}

	assert.Equal(t, want, got)
}

func Test_imageResourceModel_getCreateImageOpts(t *testing.T) {
	image := imageResourceModel{
		InstanceID: basetypes.NewStringValue("instanceId"),
		Name:       basetypes.NewStringValue("name"),
	}
	got := image.getCreateImageOpts()

	want := publiccloud.CreateImageOpts{
		Name:       "name",
		InstanceId: "instanceId",
	}

	assert.Equal(t, want, got)
}
