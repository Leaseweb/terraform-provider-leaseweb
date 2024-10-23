package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_mapSdkImageToResourceImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id:      "imageId",
		Name:    "name",
		Custom:  true,
		Flavour: "flavour",
	}

	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})
	want := resourceModelImage{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		Flavour:      basetypes.NewStringValue("flavour"),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}
	got, err := mapSdkImageToResourceImage(context.TODO(), sdkImage)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_mapSdkImageDetailsToResourceImage(t *testing.T) {
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

	want := resourceModelImage{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("READY"),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}
	got, err := mapSdkImageDetailsToResourceImage(
		context.TODO(),
		sdkImageDetails,
	)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_resourceModelImage_GetUpdateImageOpts(t *testing.T) {
	image := resourceModelImage{
		Name: basetypes.NewStringValue("name"),
	}
	got := image.GetUpdateImageOpts()

	want := publicCloud.UpdateImageOpts{Name: "name"}

	assert.Equal(t, want, got)
}

func Test_resourceModelImage_GetCreateImageOpts(t *testing.T) {
	image := resourceModelImage{
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

func Test_instanceIdForCustomImageValidator_ValidateString(t *testing.T) {
	t.Run("valid instanceId passes", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.False(t, idResponse.Diagnostics.HasError())
	})

	t.Run("non existent instanceId does not pass", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:    "tralala",
					State: publicCloud.STATE_STOPPED,
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 1)
		assert.Equal(
			t,
			`Attribute id value must be one of: ["tralala"], got: "id"`,
			idResponse.Diagnostics.Errors()[0].Detail(),
		)
	})

	t.Run("instance with state other than stopped does not pass", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_RUNNING,
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 1)
		assert.Equal(
			t,
			`Instance linked to attribute ID "id" does not have state "STOPPED", has state "RUNNING"`,
			idResponse.Diagnostics.Errors()[0].Detail(),
		)
	})

	t.Run("instance with rootDiskSize greater than 100 does not pass", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:           "id",
					State:        publicCloud.STATE_STOPPED,
					RootDiskSize: 101,
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 1)
		assert.Equal(
			t,
			`Instance linked to attribute ID "id" has rootDiskSize of 101 GB, maximum allowed size is 100 GB`,
			idResponse.Diagnostics.Errors()[0].Detail(),
		)
	})

	t.Run("instance with Windows OS does not pass", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
					Image: publicCloud.Image{
						Flavour: "windows",
					},
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 1)
		assert.Equal(
			t,
			`Instance linked to attribute ID "id" has OS "windows", only Linux & BSD are allowed`,
			idResponse.Diagnostics.Errors()[0].Detail(),
		)
	})

	t.Run("nothing is validated if id is unknown", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringUnknown()}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 0)
	})

	t.Run("nothing is validated if id is null", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringNull()}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdForCustomImageValidator(
			[]publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
				},
			},
		)
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 0)
	})
}

func Test_newInstanceIdForCustomImageValidator(t *testing.T) {
	t.Run(
		"only ids for instances with state `STOPPED` are set",
		func(t *testing.T) {
			instances := []publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
				},
				{
					Id:    "id2",
					State: publicCloud.STATE_RUNNING,
				},
			}
			instanceIdValidator := newInstanceIdForCustomImageValidator(instances)

			assert.Equal(t, []string{"id"}, instanceIdValidator.validIds)
		},
	)

	t.Run(
		"only ids for instances with rootDiskSize <= 100 are set",
		func(t *testing.T) {
			instances := []publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
				},
				{
					Id:           "id2",
					State:        publicCloud.STATE_STOPPED,
					RootDiskSize: 101,
				},
			}
			instanceIdValidator := newInstanceIdForCustomImageValidator(instances)

			assert.Equal(t, []string{"id"}, instanceIdValidator.validIds)
		},
	)

	t.Run(
		"only ids for instances with non windows OS are set",
		func(t *testing.T) {
			instances := []publicCloud.Instance{
				{
					Id:    "id",
					State: publicCloud.STATE_STOPPED,
				},
				{
					Id:    "id2",
					State: publicCloud.STATE_STOPPED,
					Image: publicCloud.Image{
						Flavour: publicCloud.FLAVOUR_WINDOWS,
					},
				},
			}
			instanceIdValidator := newInstanceIdForCustomImageValidator(instances)

			assert.Equal(t, []string{"id"}, instanceIdValidator.validIds)
		},
	)
}
