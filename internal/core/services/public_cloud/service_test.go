package public_cloud

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	sharedRepository "github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.PublicCloudRepository = &repositorySpy{}
)

type repositorySpy struct {
	instances                       public_cloud.Instances
	instanceDetails                 map[string]*public_cloud.Instance
	createdInstance                 *public_cloud.Instance
	updatedInstance                 *public_cloud.Instance
	autoScalingGroup                *public_cloud.AutoScalingGroup
	loadBalancer                    *public_cloud.LoadBalancer
	availableInstanceTypesForUpdate public_cloud.InstanceTypes
	regions                         public_cloud.Regions
	instanceTypesForRegion          public_cloud.InstanceTypes
	images                          public_cloud.Images

	passedGetAvailableInstanceTypesForUpdateId string
	passedGetAutoScalingGroupId                string
	passedGetLoadBalancerId                    string
	passedGetInstanceId                        string
	passedDeleteInstanceId                     string
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
	getAllImagesSleep              time.Duration
	getRegionsSleep                time.Duration
	getAutoScalingGroupSleep       time.Duration
	getLoadBalancerSleep           time.Duration

	getInstanceTypesForRegionCount int
	getAllImagesCount              int
	getRegionsCount                int
	getAutoScalingGroupCount       int
	getLoadBalancerCount           int
}

func (r *repositorySpy) GetAllImages(ctx context.Context) (
	public_cloud.Images,
	*sharedRepository.RepositoryError,
) {
	time.Sleep(r.getAllImagesSleep)
	r.getAllImagesCount++

	return r.images, r.getAllImagesError
}

func (r *repositorySpy) GetInstanceTypesForRegion(
	region string,
	ctx context.Context,
) (public_cloud.InstanceTypes, *sharedRepository.RepositoryError) {
	time.Sleep(r.getInstanceTypesForRegionSleep)
	r.passedGetInstanceTypesForRegionRegion = region
	r.getInstanceTypesForRegionCount++

	return r.instanceTypesForRegion, r.getInstanceTypesForRegionError
}

func (r *repositorySpy) GetRegions(ctx context.Context) (
	public_cloud.Regions,
	*sharedRepository.RepositoryError,
) {
	time.Sleep(r.getRegionsSleep)
	r.getRegionsCount++

	return r.regions, r.getRegionsError
}

func (r *repositorySpy) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) (public_cloud.InstanceTypes, *sharedRepository.RepositoryError) {
	r.passedGetAvailableInstanceTypesForUpdateId = id

	return r.availableInstanceTypesForUpdate, r.getAvailableInstanceTypesForUpdateError
}

func (r *repositorySpy) GetAutoScalingGroup(
	id string,
	ctx context.Context,
) (*public_cloud.AutoScalingGroup, *sharedRepository.RepositoryError) {
	time.Sleep(r.getAutoScalingGroupSleep)

	r.passedGetAutoScalingGroupId = id
	r.getAutoScalingGroupCount++

	return r.autoScalingGroup, r.getAutoScalingGroupError
}

func (r *repositorySpy) GetLoadBalancer(
	id string,
	ctx context.Context,
) (*public_cloud.LoadBalancer, *sharedRepository.RepositoryError) {
	time.Sleep(r.getLoadBalancerSleep)

	r.passedGetLoadBalancerId = id
	r.getLoadBalancerCount++

	return r.loadBalancer, r.getLoadBalancerError
}

func (r *repositorySpy) GetAllInstances(ctx context.Context) (
	public_cloud.Instances,
	*sharedRepository.RepositoryError,
) {
	return r.instances, r.getAllInstancesError
}

func (r *repositorySpy) GetInstance(
	id string,
	ctx context.Context,
) (*public_cloud.Instance, *sharedRepository.RepositoryError) {
	r.passedGetInstanceId = id

	return r.instanceDetails[id], r.getInstanceError
}

func (r *repositorySpy) CreateInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *sharedRepository.RepositoryError) {
	return r.createdInstance, r.createInstanceError
}

func (r *repositorySpy) UpdateInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *sharedRepository.RepositoryError) {
	return r.updatedInstance, r.updateInstanceError
}

func (r *repositorySpy) DeleteInstance(
	id string,
	ctx context.Context,
) *sharedRepository.RepositoryError {
	r.passedDeleteInstanceId = id

	return r.deleteInstanceError
}

func newRepositorySpy() repositorySpy {
	return repositorySpy{
		images: public_cloud.Images{
			public_cloud.Image{Id: "imageId"},
		},
		instanceTypesForRegion: public_cloud.InstanceTypes{
			public_cloud.InstanceType{Name: "instanceType"},
		},
		getInstanceTypesForRegionSleep: 0,
	}
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			detailedInstance := generateInstance()
			detailedInstance.Region = "region"
			detailedInstance.Id = "instanceId"

			instance := make(map[string]*public_cloud.Instance)
			instance[detailedInstance.Id] = &detailedInstance

			returnedInstances := public_cloud.Instances{{Id: "instanceId"}}

			want := public_cloud.Instances{
				public_cloud.Instance{
					Id:     "instanceId",
					Region: "region",
					Image:  public_cloud.Image{Id: "imageId"},
					Type:   public_cloud.InstanceType{Name: "instanceType"},
				},
			}

			spy := newRepositorySpy()
			spy.instances = returnedInstances
			spy.instanceDetails = instance

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
					instances: public_cloud.Instances{{}, {}, {}},
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

	t.Run(
		"instances are returned in a predictable order",
		func(t *testing.T) {
			detailedInstance1 := generateInstance()
			detailedInstance1.Id = "b"
			detailedInstance1.Region = "region"
			detailedInstance2 := generateInstance()
			detailedInstance2.Id = "a"
			detailedInstance2.Region = "region"
			instanceDetails := make(map[string]*public_cloud.Instance)
			instanceDetails["b"] = &detailedInstance1
			instanceDetails["a"] = &detailedInstance2

			returnedInstances := public_cloud.Instances{{Id: "b"}, {Id: "a"}}

			spy := newRepositorySpy()
			spy.instances = returnedInstances
			spy.instanceDetails = instanceDetails

			service := New(&spy)

			got, err := service.GetAllInstances(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, "a", got[0].Id)
			assert.Equal(t, "b", got[1].Id)
		},
	)
}

func TestService_GetInstance(t *testing.T) {
	t.Run("passes back instance from repository", func(t *testing.T) {
		want := generateInstance()
		instanceDetails := make(map[string]*public_cloud.Instance)
		instanceDetails[want.Id] = &want

		spy := newRepositorySpy()
		spy.instanceDetails = instanceDetails
		service := New(&spy)

		got, err := service.GetInstance("", context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		instance := generateInstance()
		instance.Id = "id"
		instanceDetails := make(map[string]*public_cloud.Instance)
		instanceDetails[instance.Id] = &instance

		spy := newRepositorySpy()
		spy.instanceDetails = instanceDetails
		service := New(&spy)

		want := "id"

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

			_, err := service.GetInstance("", context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"bubbles up populateMissingInstanceAttributes error",
		func(t *testing.T) {
			instance := generateInstance()
			instance.AutoScalingGroup = &public_cloud.AutoScalingGroup{}

			instanceDetails := make(map[string]*public_cloud.Instance)
			instanceDetails[instance.Id] = &instance

			spy := newRepositorySpy()
			spy.getAutoScalingGroupError = sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			)
			spy.instanceDetails = instanceDetails
			service := New(&spy)

			_, err := service.GetInstance("", context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestService_CreateInstance(t *testing.T) {
	t.Run("passes back instance from repository", func(t *testing.T) {
		want := generateInstance()
		want.Id = "instanceId"
		instanceDetails := make(map[string]*public_cloud.Instance)
		instanceDetails[want.Id] = &want

		createdInstance := generateInstance()
		createdInstance.Id = "instanceId"
		createdInstance.Image.Id = "tralala"

		spy := newRepositorySpy()
		spy.createdInstance = &createdInstance
		spy.instanceDetails = instanceDetails

		service := New(&spy)

		got, err := service.CreateInstance(public_cloud.Instance{}, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("passes back error from repository", func(t *testing.T) {
		instanceService := New(
			&repositorySpy{
				createInstanceError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			},
		)

		_, err := instanceService.CreateInstance(public_cloud.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_UpdateInstance(t *testing.T) {
	t.Run("passes back instance from repository", func(t *testing.T) {
		want := generateInstance()

		spy := newRepositorySpy()
		spy.updatedInstance = &want

		service := New(&spy)

		got, err := service.UpdateInstance(public_cloud.Instance{}, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, &want, got)
	})

	t.Run("passes back error from repository", func(t *testing.T) {
		service := New(
			&repositorySpy{
				updateInstanceError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			},
		)

		_, err := service.UpdateInstance(
			public_cloud.Instance{},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_DeleteInstance(t *testing.T) {
	t.Run("passes back nil from repository", func(t *testing.T) {
		service := New(&repositorySpy{})

		err := service.DeleteInstance("", context.TODO())

		assert.Nil(t, err)
	})

	t.Run("passes back error from repository", func(t *testing.T) {
		service := New(
			&repositorySpy{
				deleteInstanceError: sharedRepository.NewGeneralError(
					"",
					errors.New("some error"),
				),
			},
		)

		err := service.DeleteInstance("", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := "id"

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
			want := public_cloud.InstanceTypes{{Name: "tralala"}}
			spy := &repositorySpy{availableInstanceTypesForUpdate: want}

			service := New(spy)
			got, err := service.GetAvailableInstanceTypesForUpdate(
				"",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("passes back error from repository", func(t *testing.T) {
		spy := &repositorySpy{
			getAvailableInstanceTypesForUpdateError: sharedRepository.NewGeneralError(
				"",
				errors.New("some error"),
			),
		}

		service := New(spy)
		_, err := service.GetAvailableInstanceTypesForUpdate(
			"",
			context.TODO(),
		)

		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := "id"

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
			want := public_cloud.Regions{{Name: "tralala"}}
			spy := &repositorySpy{regions: want}

			service := New(spy)
			got, err := service.GetRegions(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("passes back error from repository", func(t *testing.T) {
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

	t.Run(
		"does not query repository if a local cache exists",
		func(t *testing.T) {
			spy := newRepositorySpy()
			service := New(&spy)

			_, _ = service.GetRegions(context.TODO())
			_, _ = service.GetRegions(context.TODO())

			assert.Equal(t, 1, spy.getRegionsCount)
		},
	)
}

func TestService_getAutoScalingGroup(t *testing.T) {
	t.Run("populates autoScalingGroup from repository", func(t *testing.T) {
		autoScalingGroupId := "id"

		returnedAutoScalingGroup := public_cloud.AutoScalingGroup{
			Id: autoScalingGroupId,
		}
		spy := newRepositorySpy()
		spy.autoScalingGroup = &returnedAutoScalingGroup

		service := New(&spy)

		want := &public_cloud.AutoScalingGroup{Id: autoScalingGroupId}

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

		want := "id"

		_, _ = service.getAutoScalingGroup(want, context.TODO())

		assert.Equal(t, want, spy.passedGetAutoScalingGroupId)
	})

	t.Run("populates loadBalancer from repository", func(t *testing.T) {
		loadBalancerId := "loadBalancerId"
		autoScalingGroupId := "autoScalingGroupId"

		returnedLoadBalancer := public_cloud.LoadBalancer{Id: loadBalancerId}
		returnedAutoScalingGroup := public_cloud.AutoScalingGroup{
			Id:           autoScalingGroupId,
			LoadBalancer: &public_cloud.LoadBalancer{},
		}

		service := New(&repositorySpy{
			autoScalingGroup: &returnedAutoScalingGroup,
			loadBalancer:     &returnedLoadBalancer,
		})

		want := &public_cloud.AutoScalingGroup{
			Id: autoScalingGroupId, LoadBalancer: &public_cloud.LoadBalancer{
				Id: loadBalancerId,
			},
		}

		got, err := service.getAutoScalingGroup("", context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, got)
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

			_, err := service.getAutoScalingGroup("", context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"loadBalancer errors bubble up",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getLoadBalancerError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
					autoScalingGroup: &public_cloud.AutoScalingGroup{
						LoadBalancer: &public_cloud.LoadBalancer{},
					},
				},
			)

			_, err := service.getAutoScalingGroup("", context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"does not query repository if a local cache exists",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.autoScalingGroup = &public_cloud.AutoScalingGroup{}
			service := New(&spy)

			_, _ = service.getAutoScalingGroup("", context.TODO())
			_, _ = service.getAutoScalingGroup("", context.TODO())

			assert.Equal(t, 1, spy.getAutoScalingGroupCount)
		},
	)
}

func TestService_GetAvailableInstanceTypesForRegion(t *testing.T) {
	t.Run("instanceTypes are returned", func(t *testing.T) {
		wanted := public_cloud.InstanceTypes{public_cloud.InstanceType{Name: "tralala"}}

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

	t.Run(
		"does not query repository if a local cached instanceType exists",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.instanceTypesForRegion = public_cloud.InstanceTypes{
				public_cloud.InstanceType{Name: "tralala"},
			}
			service := New(&spy)
			_, _ = service.GetAvailableInstanceTypesForRegion(
				"region",
				context.TODO(),
			)
			_, _ = service.GetAvailableInstanceTypesForRegion(
				"region",
				context.TODO(),
			)

			assert.Equal(t, 1, spy.getInstanceTypesForRegionCount)
		},
	)
}

func TestService_getImage(t *testing.T) {
	t.Run("found image is returned", func(t *testing.T) {
		spy := &repositorySpy{
			images: public_cloud.Images{public_cloud.Image{Id: "tralala"}},
		}

		service := New(spy)
		got, err := service.getImage("tralala", context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, public_cloud.Image{Id: "tralala"}, *got)
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
			images: public_cloud.Images{public_cloud.Image{Id: "tralala"}},
		}
		service := New(spy)
		_, err := service.getImage("blaat", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "blaat")
	})

	t.Run("image is cached properly", func(t *testing.T) {
		spy := &repositorySpy{
			images: public_cloud.Images{public_cloud.Image{Id: "tralala"}},
		}
		service := New(spy)

		_, _ = service.getImage("tralala", context.TODO())
		_, _ = service.getImage("tralala", context.TODO())

		assert.Equal(t, 1, spy.getAllImagesCount)
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
			instance := generateInstance()
			instance.AutoScalingGroup = &public_cloud.AutoScalingGroup{}

			autoScalingGroupDetails := public_cloud.AutoScalingGroup{}

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
			instance.AutoScalingGroup = &public_cloud.AutoScalingGroup{}

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

		want := public_cloud.Image{Id: "imageId"}
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
		spy.images = public_cloud.Images{public_cloud.Image{Id: "tralala"}}
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

	t.Run("instanceType repository errors bubble up", func(t *testing.T) {
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
	})

	t.Run("gets instanceType", func(t *testing.T) {
		want := public_cloud.InstanceType{Name: "tralala"}

		instance := generateInstance()
		instance.Type = public_cloud.InstanceType{Name: "tralala"}

		spy := newRepositorySpy()
		spy.instanceTypesForRegion = public_cloud.InstanceTypes{want}
		service := New(&spy)

		got, err := service.populateMissingInstanceAttributes(
			instance,
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, got.Type)
	})
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
		want := public_cloud.InstanceType{Name: "tralala"}

		spy := newRepositorySpy()
		spy.instanceTypesForRegion = public_cloud.InstanceTypes{want}
		service := New(&spy)

		got, err := service.getInstanceType(
			"tralala",
			"",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, want, *got)
	})

}

func Benchmark_getInstanceType(b *testing.B) {
	spy := newRepositorySpy()
	spy.instanceTypesForRegion = public_cloud.InstanceTypes{
		public_cloud.InstanceType{Name: "tralala"},
	}
	spy.getInstanceTypesForRegionSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.getInstanceType("tralala", "", context.TODO())
	}
}

func Benchmark_getImage(b *testing.B) {
	spy := newRepositorySpy()
	spy.images = public_cloud.Images{
		public_cloud.Image{Id: "tralala"},
	}
	spy.getAllImagesSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.getImage("tralala", context.TODO())
	}
}

func BenchmarkService_GetRegions(b *testing.B) {
	spy := newRepositorySpy()
	spy.getRegionsSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.GetRegions(context.TODO())
	}
}

func BenchmarkService_getAutoScalingGroup(b *testing.B) {
	spy := newRepositorySpy()
	spy.autoScalingGroup = &public_cloud.AutoScalingGroup{}
	spy.getAutoScalingGroupSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.getAutoScalingGroup("", context.TODO())
	}
}

func BenchmarkService_getLoadBalancer(b *testing.B) {
	spy := newRepositorySpy()
	spy.loadBalancer = &public_cloud.LoadBalancer{}
	spy.getLoadBalancerSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.getLoadBalancer("id", context.TODO())
	}
}

func TestService_getLoadBalancer(t *testing.T) {
	t.Run("loadBalancerId is passed to repository", func(t *testing.T) {
		want := "id"

		spy := &repositorySpy{
			getLoadBalancerError: sharedRepository.NewGeneralError(
				"",
				errors.New(""),
			),
		}
		service := New(spy)

		_, _ = service.getLoadBalancer(want, context.TODO())

		assert.Equal(t, want, spy.passedGetLoadBalancerId)
	})

	t.Run("returns loadBalancer from repository", func(t *testing.T) {
		id := "id"
		returnedLoadBalancer := public_cloud.LoadBalancer{Id: id}
		service := New(&repositorySpy{
			loadBalancer: &returnedLoadBalancer,
		})

		want := &public_cloud.LoadBalancer{Id: id}
		got, err := service.getLoadBalancer(id, context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run(
		"does not query repository if a local cache exists",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.loadBalancer = &public_cloud.LoadBalancer{}
			service := New(&spy)

			_, _ = service.getLoadBalancer("", context.TODO())
			_, _ = service.getLoadBalancer("", context.TODO())

			assert.Equal(t, 1, spy.getLoadBalancerCount)
		},
	)
}

func generateInstance() public_cloud.Instance {
	return public_cloud.Instance{
		Image: public_cloud.Image{Id: "imageId"},
		Type:  public_cloud.InstanceType{Name: "instanceType"},
	}
}
