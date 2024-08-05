package public_cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.PublicCloudService = &serviceSpy{}
)

type serviceSpy struct {
	createInstancePassedInstance *domain.Instance
	createdInstance              *domain.Instance
	createInstanceError          *serviceErrors.ServiceError

	getRegions      domain.Regions
	getRegionsError *serviceErrors.ServiceError

	instanceTypesForUpdate                  domain.InstanceTypes
	instanceTypesForUpdateError             *serviceErrors.ServiceError
	availableInstanceTypesForUpdatePassedId value_object.Uuid

	deleteInstancePassedId value_object.Uuid
	deleteInstanceError    *serviceErrors.ServiceError

	updatedInstancePassedInstance *domain.Instance
	updatedInstance               *domain.Instance
	updateInstanceError           *serviceErrors.ServiceError

	getInstances      domain.Instances
	getInstancesError *serviceErrors.ServiceError

	getInstancePassedId value_object.Uuid
	getInstance         *domain.Instance
	getInstanceError    *serviceErrors.ServiceError

	getAvailableInstanceTypesForRegion             domain.InstanceTypes
	getAvailableInstanceTypesForRegionError        *serviceErrors.ServiceError
	getAvailableInstanceTypesForRegionPassedRegion string
}

func (s *serviceSpy) GetAvailableInstanceTypesForRegion(
	region string,
	ctx context.Context,
) (domain.InstanceTypes, *serviceErrors.ServiceError) {
	s.getAvailableInstanceTypesForRegionPassedRegion = region

	return s.getAvailableInstanceTypesForRegion, s.getAvailableInstanceTypesForRegionError
}

func (s *serviceSpy) GetAllInstances(ctx context.Context) (
	domain.Instances,
	*serviceErrors.ServiceError,
) {
	return s.getInstances, s.getInstancesError
}

func (s *serviceSpy) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, *serviceErrors.ServiceError) {
	s.getInstancePassedId = id

	return s.getInstance, s.getInstanceError
}

func (s *serviceSpy) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *serviceErrors.ServiceError) {
	s.createInstancePassedInstance = &instance

	return s.createdInstance, s.createInstanceError
}

func (s *serviceSpy) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *serviceErrors.ServiceError) {
	s.updatedInstancePassedInstance = &instance

	return s.updatedInstance, s.updateInstanceError
}

func (s *serviceSpy) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) *serviceErrors.ServiceError {
	s.deleteInstancePassedId = id

	return s.deleteInstanceError
}

func (s *serviceSpy) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, *serviceErrors.ServiceError) {
	s.availableInstanceTypesForUpdatePassedId = id

	return s.instanceTypesForUpdate, s.instanceTypesForUpdateError
}

func (s *serviceSpy) GetRegions(ctx context.Context) (
	domain.Regions,
	*serviceErrors.ServiceError,
) {
	return s.getRegions, s.getRegionsError
}

func TestPublicCloudFacadeNewPublicCloudFacade(t *testing.T) {
	service := &serviceSpy{}
	facade := NewPublicCloudFacade(service)

	assert.Equal(t, service, facade.publicCloudService)
}

func TestPublicCloudFacade_CreateInstance(t *testing.T) {
	t.Run("expected instance is returned", func(t *testing.T) {
		createdInstanceId := value_object.NewGeneratedUuid()
		createdInstance := domain.Instance{Id: createdInstanceId}

		service := &serviceSpy{createdInstance: &createdInstance}

		image, _ := basetypes.NewObjectValue(
			model.Image{}.AttributeTypes(),
			map[string]attr.Value{
				"Id": basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
			},
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

		facade := NewPublicCloudFacade(service)
		facade.adaptToCreateInstanceOpts = func(
			instance model.Instance,
			allowedInstanceTypes []string,
			ctx context.Context,
		) (*domain.Instance, error) {
			return &domain.Instance{}, nil
		}

		got, err := facade.CreateInstance(instance, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, createdInstanceId.String(), got.Id.ValueString())
	})

	t.Run("error is returned if createInstanceOpts fails", func(t *testing.T) {
		spy := serviceSpy{}
		facade := NewPublicCloudFacade(&spy)
		facade.adaptToCreateInstanceOpts = func(
			instance model.Instance,
			allowedInstanceTypes []string,
			ctx context.Context,
		) (*domain.Instance, error) {
			return &domain.Instance{}, errors.New("some error")
		}

		_, err := facade.CreateInstance(model.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run(
		"error is returned if service CreateInstance fails",
		func(t *testing.T) {
			facade := PublicCloudFacade{
				adaptToCreateInstanceOpts: func(
					instance model.Instance,
					allowedInstanceTypes []string,
					ctx context.Context,
				) (*domain.Instance, error) {
					return &domain.Instance{}, nil
				},
				publicCloudService: &serviceSpy{
					createInstanceError: serviceErrors.NewError(
						"",
						errors.New("some error"),
					),
				},
			}

			_, err := facade.CreateInstance(model.Instance{}, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned if adaptInstanceToResourceModel fails",
		func(t *testing.T) {
			createdInstance := domain.Instance{Id: value_object.NewGeneratedUuid()}
			service := &serviceSpy{createdInstance: &createdInstance}
			instance := model.Instance{}

			facade := NewPublicCloudFacade(service)
			facade.adaptToCreateInstanceOpts = func(
				instance model.Instance,
				allowedInstanceTypes []string,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, nil
			}
			facade.adaptInstanceToResourceModel = func(
				instance domain.Instance,
				ctx context.Context,
			) (*model.Instance, error) {
				return nil, errors.New("some error")
			}

			_, err := facade.CreateInstance(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned if GetInstanceTypesForRegion fails",
		func(t *testing.T) {
			spy := serviceSpy{
				getAvailableInstanceTypesForRegionError: serviceErrors.NewError(
					"",
					errors.New("some error"),
				),
			}
			facade := NewPublicCloudFacade(&spy)

			_, err := facade.CreateInstance(model.Instance{}, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestPublicCloudFacade_DeleteInstance(t *testing.T) {
	t.Run("instance is deleted successfully", func(t *testing.T) {
		spy := &serviceSpy{}
		facade := PublicCloudFacade{publicCloudService: spy}

		err := facade.DeleteInstance(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			context.TODO(),
		)

		assert.Nil(t, err)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		facade := PublicCloudFacade{}

		err := facade.DeleteInstance("tralala", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := &serviceSpy{
			deleteInstanceError: serviceErrors.NewError(
				"",
				errors.New("some errors"),
			),
		}
		facade := PublicCloudFacade{publicCloudService: spy}

		err := facade.DeleteInstance(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		spy := &serviceSpy{}
		facade := PublicCloudFacade{publicCloudService: spy}
		wanted := "3cf0ddcb-b375-45a8-b18a-1bdad52527f2"

		_ = facade.DeleteInstance(wanted, context.TODO())

		assert.Equal(t, wanted, spy.deleteInstancePassedId.String())
	})
}

func TestPublicCloudFacade_DoesRegionExist(t *testing.T) {
	t.Run("returns true if region exists", func(t *testing.T) {
		want := domain.Regions{{Name: "region"}}

		spy := &serviceSpy{getRegions: want}
		facade := PublicCloudFacade{publicCloudService: spy}

		got, validRegions, err := facade.DoesRegionExist(
			"region",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, []string{"region"}, validRegions)
		assert.True(t, got)
	})

	t.Run("returns false if region does not exist", func(t *testing.T) {
		want := domain.Regions{{Name: "region"}}

		spy := &serviceSpy{getRegions: want}
		facade := PublicCloudFacade{publicCloudService: spy}

		got, validRegions, err := facade.DoesRegionExist(
			"tralala",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, []string{"region"}, validRegions)
		assert.False(t, got)
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := &serviceSpy{
			getRegionsError: serviceErrors.NewError(
				"",
				errors.New("some errors"),
			),
		}
		facade := PublicCloudFacade{publicCloudService: spy}

		_, _, err := facade.DoesRegionExist("region", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudFacade_GetInstance(t *testing.T) {
	t.Run("expected instance is returned", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()
		sdkInstance := domain.Instance{Id: instanceId}

		want := model.Instance{Id: basetypes.NewStringValue(instanceId.String())}

		spy := serviceSpy{getInstance: &sdkInstance}
		facade := PublicCloudFacade{
			publicCloudService: &spy,
			adaptInstanceToResourceModel: func(
				instance domain.Instance,
				ctx context.Context,
			) (*model.Instance, error) {
				assert.Equal(t, instanceId, instance.Id)
				return &want, nil
			},
		}

		got, err := facade.GetInstance(instanceId.String(), context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("error is returned if id is invalid", func(t *testing.T) {
		facade := PublicCloudFacade{}

		_, err := facade.GetInstance("tralala", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run(
		"error is returned if service GetInstance fails",
		func(t *testing.T) {
			facade := PublicCloudFacade{
				publicCloudService: &serviceSpy{
					getInstanceError: serviceErrors.NewError(
						"",
						errors.New("some error"),
					),
				},
			}

			_, err := facade.GetInstance(
				"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned if adaptInstanceToResourceModel fails",
		func(t *testing.T) {
			sdkInstance := domain.Instance{}

			spy := serviceSpy{getInstance: &sdkInstance}
			facade := PublicCloudFacade{
				publicCloudService: &spy,
				adaptInstanceToResourceModel: func(
					instance domain.Instance,
					ctx context.Context,
				) (*model.Instance, error) {
					return nil, errors.New("some error")
				},
			}

			_, err := facade.GetInstance(
				value_object.NewGeneratedUuid().String(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run("id is passed to repository", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()

		spy := serviceSpy{getInstanceError: &serviceErrors.ServiceError{}}
		facade := PublicCloudFacade{
			publicCloudService: &spy,
		}

		_, _ = facade.GetInstance(instanceId.String(), context.TODO())

		assert.Equal(t, instanceId, spy.getInstancePassedId)
	})
}

func TestPublicCloudFacade_GetAllInstances(t *testing.T) {
	t.Run("expected instances are returned", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()
		domainInstances := domain.Instances{{Id: instanceId}}

		modelInstances := dataSourceModel.Instances{
			Instances: []dataSourceModel.Instance{
				{Id: basetypes.NewStringValue(instanceId.String())},
			},
		}

		spy := &serviceSpy{getInstances: domainInstances}

		facade := PublicCloudFacade{
			publicCloudService: spy,
			adaptInstancesToDataSourceModel: func(instances domain.Instances) dataSourceModel.Instances {
				assert.Equal(t, instanceId, instances[0].Id)
				return modelInstances
			},
		}

		got, err := facade.GetAllInstances(context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, instanceId.String(), got.Instances[0].Id.ValueString())
	})

	t.Run(
		"error is returned if service GetAllInstances fails",
		func(t *testing.T) {
			facade := PublicCloudFacade{
				publicCloudService: &serviceSpy{
					getInstancesError: serviceErrors.NewError(
						"",
						errors.New("some error"),
					),
				},
			}

			_, err := facade.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

}

func TestPublicCloudFacade_UpdateInstance(t *testing.T) {
	t.Run("expected instance is returned", func(t *testing.T) {
		createInstanceId := value_object.NewGeneratedUuid()

		plan := model.Instance{
			Id: basetypes.NewStringValue(createInstanceId.String()),
		}
		want := model.Instance{Id: basetypes.NewStringValue("tralala")}

		instanceOpts := domain.Instance{}
		updatedInstance := domain.Instance{Id: value_object.NewGeneratedUuid()}

		spy := serviceSpy{updatedInstance: &updatedInstance}
		facade := PublicCloudFacade{
			publicCloudService: &spy,
			adaptToUpdateInstanceOpts: func(
				instance model.Instance,
				allowedInstanceTypes []string,
				ctx context.Context,
			) (*domain.Instance, error) {
				assert.Equal(
					t,
					createInstanceId.String(),
					instance.Id.ValueString(),
					"model is converted into opts",
				)

				return &instanceOpts, nil
			},
			adaptInstanceToResourceModel: func(
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

		got, err := facade.UpdateInstance(plan, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run(
		"error is returned if updatedInstancePassedInstance fails",
		func(t *testing.T) {
			spy := serviceSpy{}
			facade := NewPublicCloudFacade(&spy)
			facade.adaptToUpdateInstanceOpts = func(
				instance model.Instance,
				allowedInstanceTypes []string,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, errors.New("some error")
			}

			_, err := facade.UpdateInstance(
				model.Instance{
					Id: basetypes.NewStringValue("5072e822-485a-429a-878f-cfc42f81aca4"),
				},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned if service UpdateInstance fails",
		func(t *testing.T) {
			spy := serviceSpy{
				updateInstanceError: serviceErrors.NewError(
					"",
					errors.New("some error"),
				),
			}
			facade := NewPublicCloudFacade(&spy)
			facade.adaptToUpdateInstanceOpts = func(
				instance model.Instance,
				allowedInstanceTypes []string,
				ctx context.Context,
			) (*domain.Instance, error) {
				return &domain.Instance{}, nil
			}

			_, err := facade.UpdateInstance(
				model.Instance{
					Id: basetypes.NewStringValue("5072e822-485a-429a-878f-cfc42f81aca4"),
				},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned if adaptInstanceToResourceModel fails",
		func(t *testing.T) {
			spy := serviceSpy{updatedInstance: &domain.Instance{}}
			facade := PublicCloudFacade{
				publicCloudService: &spy,
				adaptToUpdateInstanceOpts: func(
					instance model.Instance,
					allowedInstanceTypes []string,
					ctx context.Context,
				) (*domain.Instance, error) {

					return &domain.Instance{}, nil
				},
				adaptInstanceToResourceModel: func(
					instance domain.Instance,
					ctx context.Context,
				) (*model.Instance, error) {
					return nil, errors.New("some error")
				},
			}

			_, err := facade.UpdateInstance(
				model.Instance{
					Id: basetypes.NewStringValue("5072e822-485a-429a-878f-cfc42f81aca4"),
				},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned if GetAvailableInstancesTypesForUpdate fails",
		func(t *testing.T) {
			spy := serviceSpy{
				instanceTypesForUpdateError: serviceErrors.NewError(
					"",
					errors.New("some error"),
				),
			}
			facade := NewPublicCloudFacade(&spy)

			_, err := facade.UpdateInstance(
				model.Instance{
					Id: basetypes.NewStringValue("5072e822-485a-429a-878f-cfc42f81aca4"),
				},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run("passing an invalid id returns an error", func(t *testing.T) {
		facade := PublicCloudFacade{}
		_, err := facade.UpdateInstance(
			model.Instance{
				Id: basetypes.NewStringValue("tralala"),
			},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestPublicCloudFacade_GetSshKeyRegularExpression(t *testing.T) {
	facade := PublicCloudFacade{}
	want := value_object.SshRegexp
	got := facade.GetSshKeyRegularExpression()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_GetMinimumRootDiskSize(t *testing.T) {
	facade := PublicCloudFacade{}
	want := int64(value_object.MinRootDiskSize)
	got := facade.GetMinimumRootDiskSize()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_GetMaximumRootDiskSize(t *testing.T) {
	facade := PublicCloudFacade{}
	want := int64(value_object.MaxRootDiskSize)
	got := facade.GetMaximumRootDiskSize()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_GetRootDiskStorageTypes(t *testing.T) {
	facade := PublicCloudFacade{}
	want := enum.RootDiskStorageTypeCentral.Values()
	got := facade.GetRootDiskStorageTypes()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_GetBillingFrequencies(t *testing.T) {
	facade := PublicCloudFacade{}
	want := []int64{0, 1, 3, 6, 12}
	got := facade.GetBillingFrequencies()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_GetContractTerms(t *testing.T) {
	facade := PublicCloudFacade{}
	want := []int64{0, 1, 3, 6, 12}
	got := facade.GetContractTerms()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_GetContractTypes(t *testing.T) {
	facade := PublicCloudFacade{}
	want := []string{"HOURLY", "MONTHLY"}
	got := facade.GetContractTypes()

	assert.Equal(t, want, got)
}

func TestPublicCloudFacade_ValidateContractTerm(t *testing.T) {
	t.Run(
		"ErrContractTermCannotBeZero is returned when contract returns ErrContractTermCannotBeZero",
		func(t *testing.T) {
			facade := PublicCloudFacade{}
			got := facade.ValidateContractTerm(0, "MONTHLY")

			assert.ErrorIs(t, got, ErrContractTermCannotBeZero)
		},
	)

	t.Run(
		"ErrContractTermMustBeZero is returned when contract returns ErrContractTermMustBeZero",
		func(t *testing.T) {
			facade := PublicCloudFacade{}
			got := facade.ValidateContractTerm(3, "HOURLY")

			assert.ErrorIs(t, got, ErrContractTermMustBeZero)
		},
	)

	t.Run(
		"no error is returned when contract does not return an error",
		func(t *testing.T) {
			facade := PublicCloudFacade{}
			got := facade.ValidateContractTerm(0, "HOURLY")

			assert.Nil(t, got)
		},
	)

	t.Run(
		"error is returned when invalid contractTerm is passed",
		func(t *testing.T) {
			facade := PublicCloudFacade{}
			got := facade.ValidateContractTerm(55, "HOURLY")

			assert.ErrorContains(t, got, "55")
		},
	)

	t.Run(
		"error is returned when invalid contractType is passed",
		func(t *testing.T) {
			facade := PublicCloudFacade{}
			got := facade.ValidateContractTerm(0, "tralala")

			assert.ErrorContains(t, got, "tralala")
		},
	)
}

func TestPublicCloudFacade_IsInstanceTypeAvailableForRegion(t *testing.T) {
	t.Run(
		"return true when instanceType is available for region",
		func(t *testing.T) {
			spy := serviceSpy{getAvailableInstanceTypesForRegion: domain.InstanceTypes{
				domain.InstanceType{Name: "tralala"}},
			}
			facade := NewPublicCloudFacade(&spy)

			got, instanceTypes, err := facade.IsInstanceTypeAvailableForRegion(
				"tralala",
				"region",
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, []string{"tralala"}, instanceTypes)
			assert.True(t, got)
		},
	)

	t.Run(
		"return true when instanceType is not available for region",
		func(t *testing.T) {
			spy := serviceSpy{getAvailableInstanceTypesForRegion: domain.InstanceTypes{
				domain.InstanceType{Name: "piet"}},
			}
			facade := NewPublicCloudFacade(&spy)

			got, instanceTypes, err := facade.IsInstanceTypeAvailableForRegion(
				"tralala",
				"region",
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, []string{"piet"}, instanceTypes)
			assert.False(t, got)
		},
	)

	t.Run("region is passed to service", func(t *testing.T) {
		spy := serviceSpy{}
		facade := NewPublicCloudFacade(&spy)

		_, _, _ = facade.IsInstanceTypeAvailableForRegion(
			"tralala",
			"region",
			context.TODO(),
		)

		assert.Equal(
			t,
			"region",
			spy.getAvailableInstanceTypesForRegionPassedRegion,
		)
	})

	t.Run("errors from service bubble up", func(t *testing.T) {
		spy := serviceSpy{
			getAvailableInstanceTypesForRegionError: serviceErrors.NewError(
				"prefix",
				errors.New("some error"),
			),
		}
		facade := NewPublicCloudFacade(&spy)

		_, _, err := facade.IsInstanceTypeAvailableForRegion(
			"tralala",
			"region",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestPublicCloudFacade_CanInstanceTypeBeUsedWithInstance(t *testing.T) {
	t.Run(
		"returns true if instanceType can be used with instance",
		func(t *testing.T) {
			spy := &serviceSpy{
				instanceTypesForUpdate: domain.InstanceTypes{{Name: "tralala"}},
			}
			facade := PublicCloudFacade{publicCloudService: spy}

			got, instanceTypes, err := facade.CanInstanceTypeBeUsedWithInstance(
				"085075b0-a6ad-4026-a0d1-e3256d3f7c47",
				"tralala",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"tralala"}, instanceTypes)
			assert.True(t, got)
		},
	)

	t.Run(
		"returns false if instanceType cannot be used with instance",
		func(t *testing.T) {
			spy := &serviceSpy{
				instanceTypesForUpdate: domain.InstanceTypes{{Name: "piet"}},
			}
			facade := PublicCloudFacade{publicCloudService: spy}

			got, instanceTypes, err := facade.CanInstanceTypeBeUsedWithInstance(
				"085075b0-a6ad-4026-a0d1-e3256d3f7c47",
				"tralala",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"piet"}, instanceTypes)
			assert.False(t, got)
		},
	)

	t.Run("invalid id returns an error", func(t *testing.T) {
		facade := PublicCloudFacade{}

		_, _, err := facade.CanInstanceTypeBeUsedWithInstance(
			"tralala",
			"",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := &serviceSpy{
			instanceTypesForUpdateError: serviceErrors.NewError(
				"",
				errors.New("some errors"),
			),
		}
		facade := PublicCloudFacade{publicCloudService: spy}

		_, _, err := facade.CanInstanceTypeBeUsedWithInstance(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			"",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := "085075b0-a6ad-4026-a0d1-e3256d3f7c47"

		spy := &serviceSpy{}
		facade := PublicCloudFacade{publicCloudService: spy}

		_, _, _ = facade.CanInstanceTypeBeUsedWithInstance(
			want,
			"",
			context.TODO(),
		)

		assert.Equal(t, want, spy.availableInstanceTypesForUpdatePassedId.String())
	})
}
