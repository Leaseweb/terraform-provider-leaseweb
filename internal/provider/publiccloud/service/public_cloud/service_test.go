package public_cloud

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/contracts"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
	"github.com/stretchr/testify/assert"
)

var (
	_ contracts.PublicCloudRepository = &repositorySpy{}
)

type repositorySpy struct {
	instances                       []publicCloud.Instance
	instanceDetailsById             map[string]*publicCloud.InstanceDetails
	launchedInstance                *publicCloud.Instance
	updatedInstance                 *publicCloud.InstanceDetails
	availableInstanceTypesForUpdate []string
	regions                         []string
	instanceTypesForRegion          []string

	passedGetAvailableInstanceTypesForUpdateId string
	passedGetInstanceId                        string
	passedDeleteInstanceId                     string
	passedGetInstanceTypesForRegionRegion      string

	getAllInstancesError                    *repository.RepositoryError
	getInstanceError                        *repository.RepositoryError
	launchedInstanceError                   *repository.RepositoryError
	updateInstanceError                     *repository.RepositoryError
	deleteInstanceError                     *repository.RepositoryError
	getAvailableInstanceTypesForUpdateError *repository.RepositoryError
	getRegionsError                         *repository.RepositoryError
	getInstanceTypesForRegionError          *repository.RepositoryError

	getInstanceTypesForRegionSleep time.Duration
	getRegionsSleep                time.Duration

	getInstanceTypesForRegionCount int
	getRegionsCount                int
}

func (r *repositorySpy) GetInstanceTypesForRegion(
	region string,
	ctx context.Context,
) ([]string, *repository.RepositoryError) {
	time.Sleep(r.getInstanceTypesForRegionSleep)
	r.passedGetInstanceTypesForRegionRegion = region
	r.getInstanceTypesForRegionCount++

	return r.instanceTypesForRegion, r.getInstanceTypesForRegionError
}

func (r *repositorySpy) GetRegions(ctx context.Context) (
	[]string,
	*repository.RepositoryError,
) {
	time.Sleep(r.getRegionsSleep)
	r.getRegionsCount++

	return r.regions, r.getRegionsError
}

func (r *repositorySpy) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) ([]string, *repository.RepositoryError) {
	r.passedGetAvailableInstanceTypesForUpdateId = id

	return r.availableInstanceTypesForUpdate, r.getAvailableInstanceTypesForUpdateError
}

func (r *repositorySpy) GetAllInstances(ctx context.Context) (
	[]publicCloud.Instance,
	*repository.RepositoryError,
) {
	return r.instances, r.getAllInstancesError
}

func (r *repositorySpy) GetInstance(
	id string,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *repository.RepositoryError) {
	r.passedGetInstanceId = id

	return r.instanceDetailsById[id], r.getInstanceError
}

func (r *repositorySpy) LaunchInstance(
	opts publicCloud.LaunchInstanceOpts,
	ctx context.Context,
) (*publicCloud.Instance, *repository.RepositoryError) {
	return r.launchedInstance, r.launchedInstanceError
}

func (r *repositorySpy) UpdateInstance(
	id string,
	opts publicCloud.UpdateInstanceOpts,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *repository.RepositoryError) {
	return r.updatedInstance, r.updateInstanceError
}

func (r *repositorySpy) DeleteInstance(
	id string,
	ctx context.Context,
) *repository.RepositoryError {
	r.passedDeleteInstanceId = id

	return r.deleteInstanceError
}

func newRepositorySpy() repositorySpy {
	return repositorySpy{
		instanceTypesForRegion:         []string{"instanceType"},
		getInstanceTypesForRegionSleep: 0,
		regions:                        []string{"region"},
	}
}

func TestService_GetRegions(t *testing.T) {
	t.Run(
		"expected regions returned from repository",
		func(t *testing.T) {
			want := []string{"tralala"}
			spy := &repositorySpy{regions: want}

			service := New(spy)
			got, err := service.getRegions(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run("passes back error from repository", func(t *testing.T) {
		spy := &repositorySpy{
			getRegionsError: repository.NewSdkError(
				"",
				errors.New("some error"),
				nil,
			),
		}

		service := New(spy)
		_, err := service.getRegions(context.TODO())

		assert.ErrorContains(t, err, "some error")
	})

	t.Run(
		"does not query repository if a local cache exists",
		func(t *testing.T) {
			spy := newRepositorySpy()
			service := New(&spy)

			_, _ = service.getRegions(context.TODO())
			_, _ = service.getRegions(context.TODO())

			assert.Equal(t, 1, spy.getRegionsCount)
		},
	)
}

func BenchmarkService_GetRegions(b *testing.B) {
	spy := newRepositorySpy()
	spy.getRegionsSleep = 200 * time.Millisecond

	service := New(&spy)

	for i := 0; i < b.N; i++ {
		_, _ = service.getRegions(context.TODO())
	}
}

func TestService_ValidateContractTerm(t *testing.T) {
	t.Run(
		"ErrContractTermCannotBeZero is returned when contract term is monthly and contract term is 0",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(0, "MONTHLY")

			assert.ErrorIs(t, got, ErrContractTermCannotBeZero)
		},
	)

	t.Run(
		"ErrContractTermMustBeZero is returned when contract term is hourly and contract term is not 0",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(3, "HOURLY")

			assert.ErrorIs(t, got, ErrContractTermMustBeZero)
		},
	)

	t.Run("no error is returned when contract is valid", func(t *testing.T) {
		service := Service{}
		got := service.ValidateContractTerm(0, "HOURLY")

		assert.Nil(t, got)
	},
	)

	t.Run(
		"error is returned when invalid contractTerm is passed",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(55, "HOURLY")

			assert.ErrorContains(t, got, "55")
		},
	)

	t.Run(
		"error is returned when invalid contractType is passed",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(0, "tralala")

			assert.ErrorContains(t, got, "tralala")
		},
	)
}

func TestService_DoesRegionExist(t *testing.T) {
	t.Run("returns true if region exists", func(t *testing.T) {
		spy := newRepositorySpy()
		spy.regions = []string{"region"}

		service := New(&spy)

		got, validRegions, err := service.DoesRegionExist(
			"region",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, []string{"region"}, validRegions)
		assert.True(t, got)
	})

	t.Run("returns false if region does not exist", func(t *testing.T) {
		spy := newRepositorySpy()
		spy.regions = []string{"region"}

		service := New(&spy)

		got, validRegions, err := service.DoesRegionExist(
			"tralala",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Equal(t, []string{"region"}, validRegions)
		assert.False(t, got)
	})

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := newRepositorySpy()
		spy.getRegionsError = repository.NewSdkError(
			"",
			errors.New("some error"),
			nil,
		)

		service := New(&spy)

		_, _, err := service.DoesRegionExist("region", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}
