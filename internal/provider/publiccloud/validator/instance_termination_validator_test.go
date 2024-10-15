package validator

import (
	"context"
	"errors"
	"testing"

	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
	"github.com/stretchr/testify/assert"
)

func TestInstanceTerminationValidator_ValidateObject(t *testing.T) {
	t.Run("ConfigValue populate errors bubble up", func(t *testing.T) {
		request := terraformValidator.ObjectRequest{}
		response := terraformValidator.ObjectResponse{}

		validator := ValidateInstanceTermination(
			func(
				instanceId string,
				ctx context.Context,
			) (bool, *string, *serviceErrors.ServiceError) {
				return false, nil, nil
			},
		)

		validator.ValidateObject(context.TODO(), request, &response)

		assert.True(t, response.Diagnostics.HasError())
		assert.Contains(
			t,
			response.Diagnostics[0].Summary(),
			"Value Conversion Error",
		)
	})

	t.Run("facade errors bubble up", func(t *testing.T) {
		instance := generateInstanceModel()
		instanceObject, _ := basetypes.NewObjectValueFrom(
			context.TODO(),
			instance.AttributeTypes(),
			instance,
		)
		request := terraformValidator.ObjectRequest{ConfigValue: instanceObject}
		response := terraformValidator.ObjectResponse{}

		validator := ValidateInstanceTermination(
			func(
				instanceId string,
				ctx context.Context,
			) (bool, *string, *serviceErrors.ServiceError) {
				return false, nil, serviceErrors.NewError(
					"",
					errors.New("something went wrong"),
				)
			},
		)

		validator.ValidateObject(context.TODO(), request, &response)

		assert.True(t, response.Diagnostics.HasError())
		assert.Contains(
			t,
			response.Diagnostics[0].Summary(),
			"ValidateObject",
		)
		assert.Contains(
			t,
			response.Diagnostics[0].Detail(),
			"something went wrong",
		)
	})

	t.Run(
		"does not set a diagnostics error if instance is allowed to be terminated",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instanceObject, _ := basetypes.NewObjectValueFrom(
				context.TODO(),
				instance.AttributeTypes(),
				instance,
			)
			request := terraformValidator.ObjectRequest{ConfigValue: instanceObject}
			response := terraformValidator.ObjectResponse{}

			validator := ValidateInstanceTermination(
				func(
					instanceId string,
					ctx context.Context,
				) (bool, *string, *serviceErrors.ServiceError) {
					return true, nil, nil
				},
			)

			validator.ValidateObject(context.TODO(), request, &response)

			assert.False(t, response.Diagnostics.HasError())
		},
	)

	t.Run(
		"sets a diagnostics error if instance is not allowed to be terminated",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instanceObject, _ := basetypes.NewObjectValueFrom(
				context.TODO(),
				instance.AttributeTypes(),
				instance,
			)
			request := terraformValidator.ObjectRequest{ConfigValue: instanceObject}
			response := terraformValidator.ObjectResponse{}

			validator := ValidateInstanceTermination(
				func(
					instanceId string,
					ctx context.Context,
				) (bool, *string, *serviceErrors.ServiceError) {
					reason := "reason"
					return false, &reason, nil
				},
			)

			validator.ValidateObject(context.TODO(), request, &response)

			assert.True(t, response.Diagnostics.HasError())
			assert.Contains(t, response.Diagnostics[0].Detail(), "reason")
		},
	)
}

func generateInstanceModel() resource.Instance {
	return resource.Instance{
		Id:        basetypes.NewStringUnknown(),
		Region:    basetypes.NewStringUnknown(),
		Reference: basetypes.NewStringUnknown(),
		Image: basetypes.NewObjectUnknown(
			resource.Image{}.AttributeTypes(),
		),
		State:               basetypes.NewStringUnknown(),
		Type:                basetypes.NewStringUnknown(),
		RootDiskSize:        basetypes.NewInt64Unknown(),
		RootDiskStorageType: basetypes.NewStringUnknown(),
		Ips: basetypes.NewListUnknown(
			types.ObjectType{
				AttrTypes: resource.Ip{}.AttributeTypes(),
			},
		),
		Contract: basetypes.NewObjectUnknown(
			resource.Contract{}.AttributeTypes(),
		),
		MarketAppId: basetypes.NewStringUnknown(),
		// TODO Enable SSH key support
		//SshKey: basetypes.NewStringUnknown(),
	}
}
