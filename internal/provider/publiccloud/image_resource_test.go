package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptImageToImageResource(t *testing.T) {
	sdkImage := publicCloud.Image{
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

func Test_imageResourceModel_GetUpdateImageOpts(t *testing.T) {
	image := imageResourceModel{
		Name: basetypes.NewStringValue("name"),
	}
	got := image.GetUpdateImageOpts()

	want := publicCloud.UpdateImageOpts{Name: "name"}

	assert.Equal(t, want, got)
}

func Test_imageResourceModel_GetCreateImageOpts(t *testing.T) {
	image := imageResourceModel{
		ID:   basetypes.NewStringValue("instanceId"),
		Name: basetypes.NewStringValue("name"),
	}
	got := image.GetCreateImageOpts()

	want := publicCloud.CreateImageOpts{
		Name:       "name",
		InstanceId: "instanceId",
	}

	assert.Equal(t, want, got)
}
