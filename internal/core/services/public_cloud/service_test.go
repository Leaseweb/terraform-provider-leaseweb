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
	availableInstanceTypesForUpdate public_cloud.InstanceTypes
	regions                         public_cloud.Regions
	instanceTypesForRegion          public_cloud.InstanceTypes

	passedGetAvailableInstanceTypesForUpdateId string
	passedGetInstanceId                        string
	passedDeleteInstanceId                     string
	passedGetInstanceTypesForRegionRegion      string

	getAllInstancesError                    *sharedRepository.RepositoryError
	getInstanceError                        *sharedRepository.RepositoryError
	createInstanceError                     *sharedRepository.RepositoryError
	updateInstanceError                     *sharedRepository.RepositoryError
	deleteInstanceError                     *sharedRepository.RepositoryError
	getAvailableInstanceTypesForUpdateError *sharedRepository.RepositoryError
	getRegionsError                         *sharedRepository.RepositoryError
	getInstanceTypesForRegionError          *sharedRepository.RepositoryError

	getInstanceTypesForRegionSleep time.Duration
	getRegionsSleep                time.Duration

	getInstanceTypesForRegionCount int
	getRegionsCount                int
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
		instanceTypesForRegion:         public_cloud.InstanceTypes{"instanceType"},
		getInstanceTypesForRegionSleep: 0,
		regions:                        public_cloud.Regions{"region"},
	}
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			want := public_cloud.Instances{
				public_cloud.Instance{
					Id: "instanceId",
				},
			}

			spy := newRepositorySpy()
			spy.instances = want
			spy.regions = public_cloud.Regions{"region"}

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
			want := public_cloud.InstanceTypes{"tralala"}
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
			want := public_cloud.Regions{"tralala"}
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

func TestService_GetAvailableInstanceTypesForRegion(t *testing.T) {
	t.Run("instanceTypes are returned", func(t *testing.T) {
		wants := public_cloud.InstanceTypes{"tralala"}

		spy := &repositorySpy{instanceTypesForRegion: wants}
		service := New(spy)

		got, err := service.GetAvailableInstanceTypesForRegion(
			"region",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, wants, got)
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
			spy.instanceTypesForRegion = public_cloud.InstanceTypes{"tralala"}
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

func BenchmarkService_GetRegions(b *testing.B) {
	spy := newRepositorySpy()
	spy.getRegionsSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.GetRegions(context.TODO())
	}
}

func generateInstance() public_cloud.Instance {
	return public_cloud.Instance{
		Image:  public_cloud.Image{Id: "imageId"},
		Type:   "instanceType",
		Region: "region",
	}
}
