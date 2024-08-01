package public_cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	sharedRepository "terraform-provider-leaseweb/internal/repositories/shared"
)

var (
	_ ports.PublicCloudRepository = &repositorySpy{}
)

type repositorySpy struct {
	instances                       domain.Instances
	instance                        *domain.Instance
	autoScalingGroup                *domain.AutoScalingGroup
	loadBalancer                    *domain.LoadBalancer
	availableInstanceTypesForUpdate domain.InstanceTypes
	regions                         domain.Regions
	instanceTypesForRegion          domain.InstanceTypes

	passedGetAvailableInstanceTypesForUpdateId value_object.Uuid
	passedGetAutoScalingGroupId                value_object.Uuid
	passedGetLoadBalancerId                    value_object.Uuid
	passedGetInstanceId                        value_object.Uuid
	passedDeleteInstanceId                     value_object.Uuid
	passedGetInstanceTypesForRegionRegion      string

	getAutoScalingGroupError                *sharedRepository.RepositoryError
	getLoadBalancerError                    *sharedRepository.RepositoryError
	getAllInstancesError                    *sharedRepository.RepositoryError
	getInstanceError                        *sharedRepository.RepositoryError
	createInstanceError                     *sharedRepository.RepositoryError
	updateInstanceError                     *sharedRepository.RepositoryError
	deleteInstanceError                     *sharedRepository.RepositoryError
	getAvailableInstanceTypesForUpdateError *sharedRepository.RepositoryError
	getRegionsError                         *sharedRepository.RepositoryError
	getInstanceTypesForRegionError          *sharedRepository.RepositoryError
}

func (r *repositorySpy) GetInstanceTypesForRegion(
	region string,
	ctx context.Context,
) (domain.InstanceTypes, *sharedRepository.RepositoryError) {
	r.passedGetInstanceTypesForRegionRegion = region

	return r.instanceTypesForRegion, r.getInstanceTypesForRegionError
}

func (r *repositorySpy) GetRegions(ctx context.Context) (
	domain.Regions,
	*sharedRepository.RepositoryError,
) {
	return r.regions, r.getRegionsError
}

func (r *repositorySpy) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, *sharedRepository.RepositoryError) {
	r.passedGetAvailableInstanceTypesForUpdateId = id

	return r.availableInstanceTypesForUpdate, r.getAvailableInstanceTypesForUpdateError
}

func (r *repositorySpy) GetAutoScalingGroup(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.AutoScalingGroup, *sharedRepository.RepositoryError) {
	r.passedGetAutoScalingGroupId = id

	return r.autoScalingGroup, r.getAutoScalingGroupError
}

func (r *repositorySpy) GetLoadBalancer(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.LoadBalancer, *sharedRepository.RepositoryError) {
	r.passedGetLoadBalancerId = id

	return r.loadBalancer, r.getLoadBalancerError
}

func (r *repositorySpy) GetAllInstances(ctx context.Context) (
	domain.Instances,
	*sharedRepository.RepositoryError,
) {
	return r.instances, r.getAllInstancesError
}

func (r *repositorySpy) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, *sharedRepository.RepositoryError) {
	r.passedGetInstanceId = id

	return r.instance, r.getInstanceError
}

func (r *repositorySpy) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *sharedRepository.RepositoryError) {
	return r.instance, r.createInstanceError
}

func (r *repositorySpy) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *sharedRepository.RepositoryError) {
	return r.instance, r.updateInstanceError
}

func (r *repositorySpy) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) *sharedRepository.RepositoryError {
	r.passedDeleteInstanceId = id

	return r.deleteInstanceError
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			id := value_object.NewGeneratedUuid()
			detailedInstance := domain.Instance{Id: id, Region: "region"}
			returnedInstances := domain.Instances{domain.Instance{Id: id}}

			want := domain.Instances{domain.Instance{Id: id, Region: "region"}}

			service := New(&repositorySpy{
				instances: returnedInstances,
				instance:  &detailedInstance,
			})

			got, err := service.GetAllInstances(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"error from repository getAllInstances bubbles up",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getAllInstancesError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error from repository getInstance bubbles up", func(t *testing.T) {
			service := New(
				&repositorySpy{
					instances: domain.Instances{
						{Id: value_object.NewGeneratedUuid()},
						{Id: value_object.NewGeneratedUuid()},
						{Id: value_object.NewGeneratedUuid()},
					},
					getInstanceError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestService_GetInstance(t *testing.T) {
	t.Run("passes back instance from repository", func(t *testing.T) {
		want := domain.Instance{}

		service := New(&repositorySpy{instance: &want})

		got, err := service.GetInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := value_object.NewGeneratedUuid()

		spy := &repositorySpy{instance: &domain.Instance{}}
		service := New(spy)

		_, _ = service.GetInstance(want, context.TODO())

		assert.Equal(t, want, spy.passedGetInstanceId)
	})

	t.Run(
		"bubbles up getInstance error from repository",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getInstanceError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetInstance(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"bubbles up populateMissingInstanceAttributes error",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getAutoScalingGroupError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
					instance: &domain.Instance{
						AutoScalingGroup: &domain.AutoScalingGroup{
							Id: value_object.NewGeneratedUuid(),
						},
					},
				},
			)

			_, err := service.GetInstance(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestService_CreateInstance(t *testing.T) {
	t.Run("service passes back instance from repository", func(t *testing.T) {
		want := domain.Instance{}

		service := New(&repositorySpy{instance: &want})

		got, err := service.CreateInstance(domain.Instance{}, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(
			&repositorySpy{
				createInstanceError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			},
		)

		_, err := instanceService.CreateInstance(domain.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_UpdateInstance(t *testing.T) {
	t.Run("service passes back instance from repository", func(t *testing.T) {
		want := domain.Instance{}

		service := New(&repositorySpy{instance: &want})

		got, err := service.UpdateInstance(domain.Instance{}, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, &want, got)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		service := New(
			&repositorySpy{
				updateInstanceError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			},
		)

		_, err := service.UpdateInstance(
			domain.Instance{},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_DeleteInstance(t *testing.T) {
	t.Run("service passes back nil from repository", func(t *testing.T) {
		service := New(&repositorySpy{})

		err := service.DeleteInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Nil(t, err)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		service := New(
			&repositorySpy{
				deleteInstanceError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			},
		)

		err := service.DeleteInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := value_object.NewGeneratedUuid()

		spy := &repositorySpy{}
		service := New(spy)

		_ = service.DeleteInstance(want, context.TODO())

		assert.Equal(t, want, spy.passedDeleteInstanceId)
	})
}

func TestService_GetAvailableInstanceTypesForUpdate(t *testing.T) {
	t.Run(
		"expected instance types returned from repository",
		func(t *testing.T) {
			want := domain.InstanceTypes{{Name: "tralala"}}
			spy := &repositorySpy{availableInstanceTypesForUpdate: want}

			service := New(spy)
			got, err := service.GetAvailableInstanceTypesForUpdate(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		spy := &repositorySpy{
			getAvailableInstanceTypesForUpdateError: sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			),
		}

		service := New(spy)
		_, err := service.GetAvailableInstanceTypesForUpdate(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := value_object.NewGeneratedUuid()

		spy := &repositorySpy{}
		service := New(spy)

		_, _ = service.GetAvailableInstanceTypesForUpdate(want, context.TODO())

		assert.Equal(t, want, spy.passedGetAvailableInstanceTypesForUpdateId)
	})
}

func TestService_GetRegions(t *testing.T) {
	t.Run(
		"expected regions returned from repository",
		func(t *testing.T) {
			want := domain.Regions{{Name: "tralala"}}
			spy := &repositorySpy{regions: want}

			service := New(spy)
			got, err := service.GetRegions(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		spy := &repositorySpy{
			getRegionsError: sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			),
		}

		service := New(spy)
		_, err := service.GetRegions(context.TODO())

		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_populateMissingInstanceAttributes(t *testing.T) {
	t.Run("populates autoScalingGroup from repository", func(t *testing.T) {
		autoScalingGroupId := value_object.NewGeneratedUuid()

		returnedAutoScalingGroup := domain.AutoScalingGroup{Id: autoScalingGroupId}
		instance := domain.Instance{
			AutoScalingGroup: &domain.AutoScalingGroup{
				Id: value_object.NewGeneratedUuid(),
			},
		}

		service := New(&repositorySpy{
			autoScalingGroup: &returnedAutoScalingGroup,
		})

		want := domain.Instance{
			AutoScalingGroup: &domain.AutoScalingGroup{Id: autoScalingGroupId},
		}

		got, err := service.populateMissingInstanceAttributes(
			instance,
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("autoScalingGroupId is passed to repository", func(t *testing.T) {
		want := value_object.NewGeneratedUuid()

		instance := domain.Instance{
			AutoScalingGroup: &domain.AutoScalingGroup{
				Id: want,
			},
		}

		spy := &repositorySpy{
			getAutoScalingGroupError: sharedRepository.NewGeneralError(
				"",
				errors.New(""),
			),
		}
		service := New(spy)

		_, _ = service.populateMissingInstanceAttributes(instance, context.TODO())

		assert.Equal(t, want, spy.passedGetAutoScalingGroupId)
	})

	t.Run("populates loadBalancer from repository", func(t *testing.T) {
		loadBalancerId := value_object.NewGeneratedUuid()
		autoScalingGroupId := value_object.NewGeneratedUuid()

		returnedLoadBalancer := domain.LoadBalancer{Id: loadBalancerId}
		returnedAutoScalingGroup := domain.AutoScalingGroup{
			Id:           autoScalingGroupId,
			LoadBalancer: &domain.LoadBalancer{Id: value_object.NewGeneratedUuid()},
		}
		instance := domain.Instance{
			AutoScalingGroup: &domain.AutoScalingGroup{
				Id: value_object.NewGeneratedUuid(),
			},
		}

		service := New(&repositorySpy{
			autoScalingGroup: &returnedAutoScalingGroup,
			loadBalancer:     &returnedLoadBalancer,
		})

		want := domain.Instance{
			AutoScalingGroup: &domain.AutoScalingGroup{
				Id: autoScalingGroupId, LoadBalancer: &domain.LoadBalancer{
					Id: loadBalancerId,
				},
			},
		}

		got, err := service.populateMissingInstanceAttributes(
			instance,
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("loadBalancerId is passed to repository", func(t *testing.T) {
		want := value_object.NewGeneratedUuid()

		instance := domain.Instance{
			AutoScalingGroup: &domain.AutoScalingGroup{
				LoadBalancer: &domain.LoadBalancer{Id: want},
			},
		}

		spy := &repositorySpy{
			autoScalingGroup: &domain.AutoScalingGroup{
				LoadBalancer: &domain.LoadBalancer{Id: want},
			},
			getLoadBalancerError: sharedRepository.NewGeneralError(
				"",
				errors.New(""),
			),
		}
		service := New(spy)

		_, _ = service.populateMissingInstanceAttributes(instance, context.TODO())

		assert.Equal(t, want, spy.passedGetLoadBalancerId)
	})

	t.Run(
		"bubbles up getAutoScalingGroup error from repository",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getAutoScalingGroupError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.populateMissingInstanceAttributes(
				domain.Instance{
					AutoScalingGroup: &domain.AutoScalingGroup{Id: value_object.NewGeneratedUuid()},
				},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"bubbles up getLoadBalancer error from repository",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getLoadBalancerError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
					autoScalingGroup: &domain.AutoScalingGroup{
						Id: value_object.NewGeneratedUuid(),
						LoadBalancer: &domain.LoadBalancer{
							Id: value_object.NewGeneratedUuid(),
						},
					},
				},
			)

			_, err := service.populateMissingInstanceAttributes(
				domain.Instance{
					AutoScalingGroup: &domain.AutoScalingGroup{
						Id: value_object.NewGeneratedUuid(),
						LoadBalancer: &domain.LoadBalancer{
							Id: value_object.NewGeneratedUuid(),
						},
					},
				},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestService_GetAvailableInstanceTypesForRegion(t *testing.T) {
	t.Run("instanceTypes are returned", func(t *testing.T) {
		wanted := domain.InstanceTypes{domain.InstanceType{Name: "tralala"}}

		spy := &repositorySpy{instanceTypesForRegion: wanted}
		service := New(spy)

		got, err := service.GetAvailableInstanceTypesForRegion(
			"region",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, wanted, got)
	})

	t.Run("region is passed to repository", func(t *testing.T) {
		spy := &repositorySpy{}
		service := New(spy)

		_, _ = service.GetAvailableInstanceTypesForRegion(
			"region",
			context.TODO(),
		)

		assert.Equal(t, "region", spy.passedGetInstanceTypesForRegionRegion)
	})

	t.Run("errors from repository bubble up", func(t *testing.T) {
		spy := &repositorySpy{
			getInstanceTypesForRegionError: sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			),
		}

		service := New(spy)

		_, err := service.GetAvailableInstanceTypesForRegion(
			"",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}
