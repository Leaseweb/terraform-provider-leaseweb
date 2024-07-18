package public_cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	dataSourceModel "terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var (
	_ ports.PublicCloudService = &serviceSpy{}
)

type serviceSpy struct {
	createInstanceOpts          *domain.Instance
	createdInstance             *domain.Instance
	createInstanceError         error
	getInstance                 *domain.Instance
	getInstanceError            error
	deleteInstanceError         error
	instanceTypesForUpdate      domain.InstanceTypes
	instanceTypesForUpdateError error
	getRegions                  domain.Regions
	getRegionsError             error
	getInstances                domain.Instances
	getInstancesError           error
	updateInstanceOpts          *domain.Instance
	updatedInstance             *domain.Instance
	updateInstanceError         error
}

func (s *serviceSpy) GetAllInstances(ctx context.Context) (
	domain.Instances,
	error,
) {
	return s.getInstances, s.getInstancesError
}

func (s *serviceSpy) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, error) {
	return s.getInstance, s.getInstanceError
}

func (s *serviceSpy) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	s.createInstanceOpts = &instance

	return s.createdInstance, s.createInstanceError
}

func (s *serviceSpy) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	s.updateInstanceOpts = &instance

	return s.updatedInstance, s.updateInstanceError
}

func (s *serviceSpy) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) error {
	return s.deleteInstanceError
}

func (s *serviceSpy) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, error) {
	return s.instanceTypesForUpdate, s.instanceTypesForUpdateError
}

func (s *serviceSpy) GetRegions(ctx context.Context) (domain.Regions, error) {
	return s.getRegions, s.getRegionsError
}

func TestPublicCloudHandler_NewPublicCloudHandler(t *testing.T) {
	service := &serviceSpy{}
	handler := NewPublicCloudHandler(service)

	assert.Equal(t, service, handler.publicCloudService)
}

func TestPublicCloudHandler_CreateInstance(t *testing.T) {
	t.Run("expected instance is returned", func(t *testing.T) {
		createdInstanceId := value_object.NewGeneratedUuid()
		createdInstance := domain.Instance{Id: createdInstanceId}

		service := &serviceSpy{createdInstance: &createdInstance}

		image, _ := basetypes.NewObjectValue(
			model.Image{}.AttributeTypes(),
			map[string]attr.Value{"Id": basetypes.NewStringValue("UBUNTU_20_04_64BIT")},
		)

		contract, _ := basetypes.NewObjectValue(
			model.Contract{}.AttributeTypes(),
			map[string]attr.Value{
				"Type":             basetypes.NewStringValue("MONTHLY"),
				"Term":             basetypes.NewInt64Value(3),
				"BillingFrequency": basetypes.NewInt64Value(3),
			},
		)

		instance := model.Instance{
			Region:              basetypes.NewStringValue("region"),
			Type:                basetypes.NewStringValue("lsw.m5a.4xlarge"),
			RootDiskStorageType: basetypes.NewStringValue("CENTRAL"),
			Image:               image,
			Contract:            contract,
		}

		handler := NewPublicCloudHandler(service)
		handler.convertInstanceResourceModelToCreateInstanceOpts = func(
			instance model.Instance,
			ctx context.Context,
		) (*domain.Instance, error) {
			return &domain.Instance{}, nil

		}

		got, err := handler.CreateInstance(instance, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, createdInstanceId.String(), got.Id.ValueString())
	})

	t.Run("error is returned if createInstanceOpts fails", func(t *testing.T) {
		handler := PublicCloudHandler{
			convertInstanceResourceModelToCreateInstanceOpts: func(
				instance model.Instance,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, errors.New("some error")
			},
		}

		_, err := handler.CreateInstance(model.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("error is returned if service CreateInstance fails", func(t *testing.T) {
		handler := PublicCloudHandler{
			convertInstanceResourceModelToCreateInstanceOpts: func(
				instance model.Instance,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, nil
			},
			publicCloudService: &serviceSpy{createInstanceError: errors.New("some error")},
		}

		_, err := handler.CreateInstance(model.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudHandler_DeleteInstance(t *testing.T) {
	t.Run("instance is deleted successfully", func(t *testing.T) {
		spy := &serviceSpy{}
		handler := PublicCloudHandler{publicCloudService: spy}

		err := handler.DeleteInstance(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			context.TODO(),
		)

		assert.NoError(t, err)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		handler := PublicCloudHandler{}

		err := handler.DeleteInstance("tralala", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := &serviceSpy{deleteInstanceError: errors.New("some errors")}
		handler := PublicCloudHandler{publicCloudService: spy}

		err := handler.DeleteInstance(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudHandler_GetAvailableInstanceTypesForUpdate(t *testing.T) {
	t.Run("expected instanceTypes are returned", func(t *testing.T) {
		want := domain.InstanceTypes{{Name: "tralala"}}

		spy := &serviceSpy{instanceTypesForUpdate: want}
		handler := PublicCloudHandler{publicCloudService: spy}

		got, err := handler.GetAvailableInstanceTypesForUpdate(
			"085075b0-a6ad-4026-a0d1-e3256d3f7c47",
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("invalid id returns an error", func(t *testing.T) {
		handler := PublicCloudHandler{}

		_, err := handler.GetAvailableInstanceTypesForUpdate(
			"tralala",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := &serviceSpy{
			instanceTypesForUpdateError: errors.New("some errors"),
		}
		handler := PublicCloudHandler{publicCloudService: spy}

		_, err := handler.GetAvailableInstanceTypesForUpdate(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudHandler_GetRegions(t *testing.T) {
	t.Run("expected regions are returned", func(t *testing.T) {
		want := domain.Regions{{Name: "region"}}

		spy := &serviceSpy{getRegions: want}
		handler := PublicCloudHandler{publicCloudService: spy}

		got, err := handler.GetRegions(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := &serviceSpy{
			getRegionsError: errors.New("some errors"),
		}
		handler := PublicCloudHandler{publicCloudService: spy}

		_, err := handler.GetRegions(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudHandler_GetInstance(t *testing.T) {
	t.Run("expected instance is returned", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()
		sdkInstance := domain.Instance{Id: instanceId}

		want := model.Instance{Id: basetypes.NewStringValue(instanceId.String())}

		spy := serviceSpy{getInstance: &sdkInstance}
		handler := PublicCloudHandler{
			publicCloudService: &spy,
			convertInstanceToResourceModel: func(
				instance domain.Instance,
				ctx context.Context,
			) (*model.Instance, error) {
				assert.Equal(t, instanceId, instance.Id)
				return &want, nil
			},
		}

		got, err := handler.GetInstance(instanceId.String(), context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, want, *got)

	})

	t.Run("error is returned if id is invalid", func(t *testing.T) {
		handler := PublicCloudHandler{}

		_, err := handler.GetInstance("tralala", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run(
		"error is returned if service GetInstance fails",
		func(t *testing.T) {
			handler := PublicCloudHandler{
				publicCloudService: &serviceSpy{getInstanceError: errors.New("some error")},
			}

			_, err := handler.GetInstance(
				"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestPublicCloudHandler_GetAllInstances(t *testing.T) {
	t.Run("expected instances are returned", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()
		domainInstances := domain.Instances{{Id: instanceId}}

		modelInstances := dataSourceModel.Instances{
			Instances: []dataSourceModel.Instance{
				{Id: basetypes.NewStringValue(instanceId.String())},
			},
		}

		spy := &serviceSpy{getInstances: domainInstances}

		handler := PublicCloudHandler{
			publicCloudService: spy,
			convertInstancesToDataSourceModel: func(instances domain.Instances) dataSourceModel.Instances {
				assert.Equal(t, instanceId, instances[0].Id)
				return modelInstances
			},
		}

		got, err := handler.GetAllInstances(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, instanceId.String(), got.Instances[0].Id.ValueString())
	})

	t.Run(
		"error is returned if service GetAllInstances fails",
		func(t *testing.T) {
			handler := PublicCloudHandler{
				publicCloudService: &serviceSpy{getInstancesError: errors.New("some error")},
			}

			_, err := handler.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

}

func TestPublicCloudHandler_UpdateInstance(t *testing.T) {
	t.Run("expected instance is returned", func(t *testing.T) {
		plan := model.Instance{Id: basetypes.NewStringValue("CreateInstance")}
		want := model.Instance{Id: basetypes.NewStringValue("tralala")}

		instanceOpts := domain.Instance{}
		updatedInstance := domain.Instance{Id: value_object.NewGeneratedUuid()}

		spy := serviceSpy{updatedInstance: &updatedInstance}
		handler := PublicCloudHandler{
			publicCloudService: &spy,
			convertInstanceResourceModelToUpdateInstanceOpts: func(
				instance model.Instance,
				ctx context.Context,
			) (*domain.Instance, error) {
				assert.Equal(
					t,
					"CreateInstance",
					instance.Id.ValueString(),
					"model is converted into opts",
				)

				return &instanceOpts, nil
			},
			convertInstanceToResourceModel: func(
				instance domain.Instance,
				ctx context.Context,
			) (*model.Instance, error) {
				assert.Equal(
					t,
					updatedInstance.Id,
					instance.Id,
					"instance from repository is converted into model")

				return &want, nil
			},
		}

		got, err := handler.UpdateInstance(plan, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("error is returned if updateInstanceOpts fails", func(t *testing.T) {
		handler := PublicCloudHandler{
			convertInstanceResourceModelToUpdateInstanceOpts: func(
				instance model.Instance,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, errors.New("some error")
			},
		}

		_, err := handler.UpdateInstance(model.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("error is returned if service UpdateInstance fails", func(t *testing.T) {
		handler := PublicCloudHandler{
			convertInstanceResourceModelToUpdateInstanceOpts: func(
				instance model.Instance,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, nil
			},
			publicCloudService: &serviceSpy{updateInstanceError: errors.New("some error")},
		}

		_, err := handler.UpdateInstance(model.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudHandler_GetImageIds(t *testing.T) {
	handler := PublicCloudHandler{}
	want := enum.Debian1064Bit.Values()
	got := handler.GetImageIds()

	assert.Equal(t, want, got)
}

func TestPublicCloudHandler_GetSshKeyRegularExpression(t *testing.T) {
	handler := PublicCloudHandler{}
	want := value_object.SshRegexp
	got := handler.GetSshKeyRegularExpression()

	assert.Equal(t, want, got)
}

func TestPublicCloudHandler_GetMinimumRootDiskSize(t *testing.T) {
	handler := PublicCloudHandler{}
	want := int64(value_object.MinRootDiskSize)
	got := handler.GetMinimumRootDiskSize()

	assert.Equal(t, want, got)
}

func TestPublicCloudHandler_GetMaximumRootDiskSize(t *testing.T) {
	handler := PublicCloudHandler{}
	want := int64(value_object.MaxRootDiskSize)
	got := handler.GetMaximumRootDiskSize()

	assert.Equal(t, want, got)
}
