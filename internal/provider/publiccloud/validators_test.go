package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_contractTermValidator_ValidateObject(t *testing.T) {
	t.Run(
		"does not set error if contract term is correct",
		func(t *testing.T) {
			contract := contractResourceModel{}
			configValue, _ := types.ObjectValueFrom(
				context.TODO(),
				contract.AttributeTypes(),
				contract,
			)

			request := validator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := validator.ObjectResponse{}

			contractTermValidator := contractTermValidator{}
			contractTermValidator.ValidateObject(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 0)
		},
	)

	t.Run(
		"returns expected error if contract term cannot be 0",
		func(t *testing.T) {
			contract := contractResourceModel{
				Type: basetypes.NewStringValue("MONTHLY"),
				Term: basetypes.NewInt32Value(0),
			}
			configValue, _ := types.ObjectValueFrom(
				context.TODO(),
				contract.AttributeTypes(),
				contract,
			)

			request := validator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := validator.ObjectResponse{}

			contractTermValidator := contractTermValidator{}
			contractTermValidator.ValidateObject(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 1)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"MONTHLY",
			)
		},
	)

	t.Run(
		"returns expected error if contract term must be 0",
		func(t *testing.T) {
			contract := contractResourceModel{
				Type: basetypes.NewStringValue("HOURLY"),
				Term: basetypes.NewInt32Value(3),
			}
			configValue, _ := types.ObjectValueFrom(
				context.TODO(),
				contract.AttributeTypes(),
				contract,
			)

			request := validator.ObjectRequest{
				ConfigValue: configValue,
			}

			response := validator.ObjectResponse{}

			contractTermValidator := contractTermValidator{}
			contractTermValidator.ValidateObject(context.TODO(), request, &response)

			assert.Len(t, response.Diagnostics.Errors(), 1)
			assert.Contains(
				t,
				response.Diagnostics.Errors()[0].Detail(),
				"HOURLY",
			)
		},
	)
}

func Test_instanceTerminationValidator_ValidateObject(t *testing.T) {
	t.Run("ConfigValue populate errors bubble up", func(t *testing.T) {
		request := validator.ObjectRequest{}
		response := validator.ObjectResponse{}

		instanceTerminationValidator := instanceTerminationValidator{}
		instanceTerminationValidator.ValidateObject(context.TODO(), request, &response)

		assert.True(t, response.Diagnostics.HasError())
		assert.Contains(
			t,
			response.Diagnostics[0].Summary(),
			"Value Conversion Error",
		)
	})

	t.Run(
		"does not set a diagnostics error if instance is allowed to be terminated",
		func(t *testing.T) {
			instance := generateInstanceModelForValidator()
			instanceObject, _ := basetypes.NewObjectValueFrom(
				context.TODO(),
				instance.AttributeTypes(),
				instance,
			)
			request := validator.ObjectRequest{ConfigValue: instanceObject}
			response := validator.ObjectResponse{}

			instanceTerminationValidator := instanceTerminationValidator{}
			instanceTerminationValidator.ValidateObject(context.TODO(), request, &response)

			assert.False(t, response.Diagnostics.HasError())
		},
	)

	t.Run(
		"sets a diagnostics error if instance is not allowed to be terminated",
		func(t *testing.T) {
			instance := generateInstanceModelForValidator()
			instance.State = basetypes.NewStringValue("DESTROYED")
			instanceObject, _ := basetypes.NewObjectValueFrom(
				context.TODO(),
				instance.AttributeTypes(),
				instance,
			)
			request := validator.ObjectRequest{ConfigValue: instanceObject}
			response := validator.ObjectResponse{}

			instanceTerminationValidator := instanceTerminationValidator{}
			instanceTerminationValidator.ValidateObject(context.TODO(), request, &response)

			assert.True(t, response.Diagnostics.HasError())
			assert.Contains(t, response.Diagnostics[0].Detail(), "DESTROYED")
		},
	)
}

func generateInstanceModelForValidator() instanceResourceModel {
	contract := contractResourceModel{}
	contractObject, _ := types.ObjectValueFrom(
		context.TODO(),
		contract.AttributeTypes(),
		contract,
	)

	return instanceResourceModel{
		ID:        basetypes.NewStringUnknown(),
		Region:    basetypes.NewStringUnknown(),
		Reference: basetypes.NewStringUnknown(),
		Image: basetypes.NewObjectUnknown(
			imageResourceModel{}.AttributeTypes(),
		),
		State:               basetypes.NewStringUnknown(),
		Type:                basetypes.NewStringUnknown(),
		RootDiskSize:        basetypes.NewInt32Unknown(),
		RootDiskStorageType: basetypes.NewStringUnknown(),
		IPs: basetypes.NewListUnknown(
			types.ObjectType{
				AttrTypes: iPResourceModel{}.AttributeTypes(),
			},
		),
		Contract:    contractObject,
		MarketAppID: basetypes.NewStringUnknown(),
	}
}

func Test_instanceTypeValidator_ValidateString(t *testing.T) {
	t.Run("nothing happens if instanceType is unknown", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		instanceTypeValidator := instanceTypeValidator{}

		response := validator.StringResponse{}
		instanceTypeValidator.ValidateString(
			context.TODO(),
			validator.StringRequest{ConfigValue: basetypes.NewStringUnknown()},
			&response,
		)

		assert.Equal(t, 0, countIsInstanceTypeAvailableForRegionIsCalled)
		assert.Equal(t, 0, countCanInstanceTypeBeUsedWithInstanceIsCalled)
	})

	t.Run("nothing happens if instanceType does not change", func(t *testing.T) {
		countIsInstanceTypeAvailableForRegionIsCalled := 0
		countCanInstanceTypeBeUsedWithInstanceIsCalled := 0

		instanceTypeValidator := instanceTypeValidator{}

		response := validator.StringResponse{}
		instanceTypeValidator.ValidateString(
			context.TODO(),
			validator.StringRequest{
				ConfigValue: basetypes.NewStringNull(),
			},
			&response,
		)

		assert.Equal(t, 0, countIsInstanceTypeAvailableForRegionIsCalled)
		assert.Equal(t, 0, countCanInstanceTypeBeUsedWithInstanceIsCalled)
	})

	t.Run(
		"attributeError added to response if instanceType cannot be found",
		func(t *testing.T) {
			instanceTypeValidator := instanceTypeValidator{
				availableInstanceTypes: []string{"tralala"},
			}

			response := validator.StringResponse{}
			instanceTypeValidator.ValidateString(
				context.TODO(),
				validator.StringRequest{
					ConfigValue: basetypes.NewStringValue("doesNotExist"),
				},
				&response,
			)

			assert.Contains(
				t,
				response.Diagnostics[0].Detail(),
				"tralala",
			)
			assert.Contains(
				t,
				response.Diagnostics[0].Detail(),
				"doesNotExist",
			)
		},
	)

	t.Run(
		"attributeError not added to response if instanceType can be found",
		func(t *testing.T) {
			instanceTypeValidator := instanceTypeValidator{
				availableInstanceTypes: []string{"tralala"},
			}

			response := validator.StringResponse{}
			instanceTypeValidator.ValidateString(
				context.TODO(),
				validator.StringRequest{
					ConfigValue: basetypes.NewStringValue("tralala"),
				},
				&response,
			)

			assert.Len(t, response.Diagnostics, 0)
		},
	)
}

func Test_newInstanceTypeValidator(t *testing.T) {
	instanceTypeValidator := newInstanceTypeValidator(
		basetypes.NewStringValue("currentInstanceType"),
		[]string{"type1"},
	)

	assert.Equal(
		t,
		[]string{"type1", "currentInstanceType"},
		instanceTypeValidator.availableInstanceTypes,
	)
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
