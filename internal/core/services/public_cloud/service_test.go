package public_cloud

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	sharedRepository "github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
	"github.com/stretchr/testify/assert"
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
	images                          domain.Images

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
	getAllImagesError                       *sharedRepository.RepositoryError

	getInstanceTypesForRegionSleep time.Duration
	// How many times has getInstanceTypesForRegion been called.
	getInstanceTypesForRegionCount int
}

func (r *repositorySpy) GetAllImages(ctx context.Context) (
	domain.Images,
	*sharedRepository.RepositoryError,
) {
	return r.images, r.getAllImagesError
}

func (r *repositorySpy) GetInstanceTypesForRegion(
	region string,
	ctx context.Context,
) (domain.InstanceTypes, *sharedRepository.RepositoryError) {
	time.Sleep(r.getInstanceTypesForRegionSleep)
	r.passedGetInstanceTypesForRegionRegion = region
	r.getInstanceTypesForRegionCount++

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

func newRepositorySpy() repositorySpy {
	return repositorySpy{
		images: domain.Images{
			domain.Image{Id: "imageId"},
		},
		instanceTypesForRegion: domain.InstanceTypes{
			domain.InstanceType{Name: "instanceType"},
		},
		getInstanceTypesForRegionSleep: 0,
	}
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			id := value_object.NewGeneratedUuid()
			detailedInstance := generateInstance()
			detailedInstance.Region = "region"
			detailedInstance.Id = id

			returnedInstances := domain.Instances{
				domain.Instance{Id: id},
			}

			want := domain.Instances{
				domain.Instance{
					Id:     id,
					Region: "region",
					Image:  domain.Image{Id: "imageId"},
					Type:   domain.InstanceType{Name: "instanceType"},
				},
			}

			spy := newRepositorySpy()
			spy.instances = returnedInstances
			spy.instance = &detailedInstance

			service := New(&spy)

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
		want := generateInstance()

		spy := newRepositorySpy()
		spy.instance = &want

		service := New(&spy)

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
		want := generateInstance()

		spy := newRepositorySpy()
		spy.instance = &want

		service := New(&spy)

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
		want := generateInstance()

		spy := newRepositorySpy()
		spy.instance = &want

		service := New(&spy)

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

func TestService_getAutoScalingGroup(t *testing.T) {
	t.Run("populates autoScalingGroup from repository", func(t *testing.T) {
		autoScalingGroupId := value_object.NewGeneratedUuid()

		returnedAutoScalingGroup := domain.AutoScalingGroup{Id: autoScalingGroupId}
		spy := newRepositorySpy()
		spy.autoScalingGroup = &returnedAutoScalingGroup

		service := New(&spy)

		want := &domain.AutoScalingGroup{Id: autoScalingGroupId}

		got, err := service.getAutoScalingGroup(autoScalingGroupId, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("autoScalingGroupId is passed to repository", func(t *testing.T) {

		spy := &repositorySpy{
			getAutoScalingGroupError: sharedRepository.NewGeneralError(
				"",
				errors.New(""),
			),
		}
		service := New(spy)

		want := value_object.NewGeneratedUuid()

		_, _ = service.getAutoScalingGroup(want, context.TODO())

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

		service := New(&repositorySpy{
			autoScalingGroup: &returnedAutoScalingGroup,
			loadBalancer:     &returnedLoadBalancer,
		})

		want := &domain.AutoScalingGroup{
			Id: autoScalingGroupId, LoadBalancer: &domain.LoadBalancer{
				Id: loadBalancerId,
			},
		}

		got, err := service.getAutoScalingGroup(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("loadBalancerId is passed to repository", func(t *testing.T) {
		want := value_object.NewGeneratedUuid()

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

		_, _ = service.getAutoScalingGroup(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

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

			_, err := service.getAutoScalingGroup(
				value_object.NewGeneratedUuid(),
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

			_, err := service.getAutoScalingGroup(
				value_object.NewGeneratedUuid(),
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

func TestService_getImage(t *testing.T) {
	t.Run("found image is returned", func(t *testing.T) {
		spy := &repositorySpy{
			images: domain.Images{domain.Image{Id: "tralala"}},
		}

		service := New(spy)
		got, err := service.getImage("tralala", context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, domain.Image{Id: "tralala"}, *got)
	})

	t.Run("errors bubble up from repository", func(t *testing.T) {
		spy := &repositorySpy{
			getAllImagesError: sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			),
		}

		service := New(spy)
		_, err := service.getImage("tralala", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("errors i returned if image cannot be found", func(t *testing.T) {
		spy := &repositorySpy{
			images: domain.Images{domain.Image{Id: "tralala"}},
		}

		service := New(spy)
		_, err := service.getImage("blaat", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "blaat")
	})
}

func TestService_populateMissingInstanceAttributes(t *testing.T) {
	t.Run(
		"autoScalingGroupDetails are not retrieved if instance.autoScalingGroup is nil",
		func(t *testing.T) {
			instance := generateInstance()

			spy := newRepositorySpy()

			service := New(&spy)

			got, err := service.populateMissingInstanceAttributes(
				instance,
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Nil(t, got.AutoScalingGroup)
		},
	)

	t.Run(
		"sets autoScalingGroupDetails if autoScalingGroup is set",
		func(t *testing.T) {
			autoScalingGroupId := value_object.NewGeneratedUuid()
			instance := generateInstance()
			instance.AutoScalingGroup = &domain.AutoScalingGroup{
				Id: autoScalingGroupId,
			}

			autoScalingGroupDetails := domain.AutoScalingGroup{
				Id: value_object.NewGeneratedUuid(),
			}

			spy := newRepositorySpy()
			spy.autoScalingGroup = &autoScalingGroupDetails
			service := New(&spy)

			got, err := service.populateMissingInstanceAttributes(
				instance,
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, autoScalingGroupDetails, *got.AutoScalingGroup)
		},
	)

	t.Run(
		"autoScalingGroupErrors bubble up", func(t *testing.T) {
			instance := generateInstance()
			instance.AutoScalingGroup = &domain.AutoScalingGroup{}

			spy := &repositorySpy{
				getAutoScalingGroupError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			}
			service := New(spy)

			_, err := service.populateMissingInstanceAttributes(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		})

	t.Run("gets imageDetails", func(t *testing.T) {
		instance := generateInstance()

		spy := newRepositorySpy()
		service := New(&spy)

		want := domain.Image{Id: "imageId"}
		got, err := service.populateMissingInstanceAttributes(
			instance,
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, got.Image)
	})

	t.Run("imageDetails errors bubble up", func(t *testing.T) {
		instance := generateInstance()

		spy := newRepositorySpy()
		spy.images = domain.Images{domain.Image{Id: "tralala"}}
		spy.getAllImagesError = sharedRepository.NewGeneralError(
			"",
			errors.New("some error"),
		)

		service := New(&spy)

		_, err := service.populateMissingInstanceAttributes(
			instance,
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run(
		"instanceType repository errors bubble up",
		func(t *testing.T) {
			instance := generateInstance()

			spy := newRepositorySpy()
			spy.getInstanceTypesForRegionError = sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			)
			service := New(&spy)

			_, err := service.populateMissingInstanceAttributes(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run("gets instanceType", func(t *testing.T) {
		want := domain.InstanceType{Name: "tralala"}

		instance := generateInstance()
		instance.Type = domain.InstanceType{Name: "tralala"}

		spy := newRepositorySpy()
		spy.instanceTypesForRegion = domain.InstanceTypes{want}
		service := New(&spy)

		got, err := service.populateMissingInstanceAttributes(
			instance,
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, got.Type)
	})
}

func generateInstance() domain.Instance {
	return domain.Instance{
		Image: domain.Image{Id: "imageId"},
		Type:  domain.InstanceType{Name: "instanceType"},
	}
}

func TestService_getInstanceType(t *testing.T) {
	t.Run("errors bubble up from the repository", func(t *testing.T) {
		spy := newRepositorySpy()
		spy.getInstanceTypesForRegionError = sharedRepository.NewGeneralError(
			"",
			errors.New("some error"),
		)
		service := New(&spy)

		_, err := service.getInstanceType("", "", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("region is passed to repository", func(t *testing.T) {
		spy := newRepositorySpy()
		service := New(&spy)

		_, _ = service.getInstanceType("", "region", context.TODO())

		assert.Equal(t, "region", spy.passedGetInstanceTypesForRegionRegion)
	})

	t.Run(
		"error is returned if instanceType is not found",
		func(t *testing.T) {
			spy := newRepositorySpy()
			service := New(&spy)

			_, err := service.getInstanceType(
				"tralala",
				"",
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run("instanceType is returned if found", func(t *testing.T) {
		want := domain.InstanceType{Name: "tralala"}

		spy := newRepositorySpy()
		spy.instanceTypesForRegion = domain.InstanceTypes{want}
		service := New(&spy)

		got, err := service.getInstanceType(
			"tralala",
			"",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run(
		"does not query repository if a local cached instanceType exists",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.instanceTypesForRegion = domain.InstanceTypes{
				domain.InstanceType{Name: "tralala"},
			}
			service := New(&spy)
			_, _ = service.getInstanceType(
				"name",
				"region",
				context.TODO(),
			)
			_, _ = service.getInstanceType("name", "region", context.TODO())

			assert.Equal(t, 1, spy.getInstanceTypesForRegionCount)
		},
	)
}

func Benchmark_getInstanceType(b *testing.B) {
	spy := newRepositorySpy()
	spy.instanceTypesForRegion = domain.InstanceTypes{
		domain.InstanceType{Name: "tralala"},
	}
	spy.getInstanceTypesForRegionSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {

		_, _ = service.getInstanceType(
			"tralala",
			"",
			context.TODO(),
		)
	}
}
