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

func Test_newResourceModelImageFromImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id:      "imageId",
		Name:    "name",
		Custom:  true,
		Flavour: "flavour",
	}

	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})
	want := resourceModelImage{
		Id:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		Flavour:      basetypes.NewStringValue("flavour"),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}
	got, err := newResourceModelImageFromImage(context.TODO(), sdkImage)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func Test_newResourceModelImageFromImageDetails(t *testing.T) {
	state := "RUNNING"
	region := publicCloud.REGIONNAME_EU_WEST_3

	sdkImageDetails := publicCloud.ImageDetails{
		Id:           "imageId",
		Name:         "name",
		Custom:       true,
		State:        *publicCloud.NewNullableString(&state),
		MarketApps:   []string{"marketApp"},
		StorageTypes: []string{"storageType"},
		Flavour:      "flavour",
		Region:       *publicCloud.NewNullableRegionName(&region),
	}

	marketApps, _ := basetypes.NewListValueFrom(
		context.TODO(),
		types.StringType,
		[]string{"marketApp"},
	)
	storageTypes, _ := basetypes.NewListValueFrom(
		context.TODO(),
		types.StringType,
		[]string{"storageType"},
	)

	want := resourceModelImage{
		Id:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("RUNNING"),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}
	got, err := newResourceModelImageFromImageDetails(
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
		Id:   basetypes.NewStringValue("instanceId"),
		Name: basetypes.NewStringValue("name"),
	}
	got := image.GetCreateImageOpts()

	want := publicCloud.CreateImageOpts{
		Name:       "name",
		InstanceId: "instanceId",
	}

	assert.Equal(t, want, got)
}

func Test_instanceIdValidator_ValidateString(t *testing.T) {
	t.Run("valid instanceId passes", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdValidator([]publicCloud.Instance{{Id: "id"}})
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 0)
	})

	t.Run("invalid instanceId does not pass", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringValue("id")}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdValidator([]publicCloud.Instance{{Id: "tralala"}})
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 1)
	})

	t.Run("nothing is validated if id is unknown", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringUnknown()}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdValidator([]publicCloud.Instance{{Id: "tralala"}})
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 0)
	})

	t.Run("nothing is validated if id is null", func(t *testing.T) {
		idRequest := validator.StringRequest{ConfigValue: basetypes.NewStringNull()}
		idResponse := validator.StringResponse{}

		instanceIdValidator := newInstanceIdValidator([]publicCloud.Instance{{Id: "tralala"}})
		instanceIdValidator.ValidateString(context.TODO(), idRequest, &idResponse)

		assert.Len(t, idResponse.Diagnostics.Errors(), 0)
	})
}
