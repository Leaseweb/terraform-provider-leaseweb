package validator

import (
	"context"
	"errors"
	"testing"

	terraformValidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
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
			) (bool, *public_cloud.CannotBeTerminatedReason, error) {
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
			) (bool, *public_cloud.CannotBeTerminatedReason, error) {
				return false, nil, errors.New("something went wrong")
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
				) (bool, *public_cloud.CannotBeTerminatedReason, error) {
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
				) (bool, *public_cloud.CannotBeTerminatedReason, error) {
					reason := public_cloud.CannotBeTerminatedReason("reason")
					return false, &reason, nil
				},
			)

			validator.ValidateObject(context.TODO(), request, &response)

			assert.True(t, response.Diagnostics.HasError())
			assert.Contains(t, response.Diagnostics[0].Detail(), "reason")
		},
	)
}

func generateInstanceModel() model.Instance {
	return model.Instance{
		Id:        basetypes.NewStringUnknown(),
		Region:    basetypes.NewStringUnknown(),
		Reference: basetypes.NewStringUnknown(),
		Resources: basetypes.NewObjectUnknown(
			model.Resources{}.AttributeTypes(),
		),
		Image: basetypes.NewObjectUnknown(
			model.Image{}.AttributeTypes(),
		),
		State:               basetypes.NewStringUnknown(),
		ProductType:         basetypes.NewStringUnknown(),
		HasPublicIpv4:       basetypes.NewBoolUnknown(),
		HasPrivateNetwork:   basetypes.NewBoolUnknown(),
		Type:                basetypes.NewStringUnknown(),
		RootDiskSize:        basetypes.NewInt64Unknown(),
		RootDiskStorageType: basetypes.NewStringUnknown(),
		Ips: basetypes.NewListUnknown(
			types.ObjectType{
				AttrTypes: model.Ip{}.AttributeTypes(),
			},
		),
		StartedAt: basetypes.NewStringUnknown(),
		Contract: basetypes.NewObjectUnknown(
			model.Contract{}.AttributeTypes(),
		),
		MarketAppId: basetypes.NewStringUnknown(),
		AutoScalingGroup: basetypes.NewObjectUnknown(
			model.AutoScalingGroup{}.AttributeTypes(),
		),
		Iso: basetypes.NewObjectUnknown(model.Iso{}.AttributeTypes()),
		PrivateNetwork: basetypes.NewObjectUnknown(
			model.PrivateNetwork{}.AttributeTypes(),
		),
		// TODO Enable SSH key support
		//SshKey: basetypes.NewStringUnknown(),
	}
}
