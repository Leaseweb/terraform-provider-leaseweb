package public_cloud

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/contracts"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
	shared2 "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service"
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

func TestService_UpdateInstance(t *testing.T) {
	t.Run("passes back instance from repository", func(t *testing.T) {
		updatedInstance := generateInstanceDetails()
		updatedInstance.Id = "instanceId"

		spy := newRepositorySpy()
		spy.updatedInstance = &updatedInstance

		service := New(&spy)

		got, err := service.UpdateInstance(generateInstanceModel(), context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, "instanceId", got.Id)
	})

	t.Run("passes back error from repository", func(t *testing.T) {
		service := New(
			&repositorySpy{
				updateInstanceError: repository.NewSdkError(
					"",
					errors.New("some error"),
					nil,
				),
			},
		)

		_, err := service.UpdateInstance(generateInstanceModel(), context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("bubbles up error from adaptToUpdateInstanceOpts", func(t *testing.T) {
		spy := newRepositorySpy()
		service := New(&spy)
		service.adaptToUpdateInstanceOpts = func(
			instance resource.Instance,
			ctx context.Context,
		) (*publicCloud.UpdateInstanceOpts, error) {
			return nil, errors.New("some error")
		}

		_, err := service.UpdateInstance(generateInstanceModel(), context.TODO())

		assert.NotNil(t, err)
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
				deleteInstanceError: repository.NewSdkError(
					"",
					errors.New("some error"),
					nil,
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
			want := []string{"tralala"}
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
			getAvailableInstanceTypesForUpdateError: repository.NewSdkError(
				"",
				errors.New("some error"),
				nil,
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
			want := []string{"tralala"}
			spy := &repositorySpy{regions: want}

			service := New(spy)
			got, err := service.GetRegions(context.TODO())

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
		wants := []string{"tralala"}

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
			getInstanceTypesForRegionError: repository.NewSdkError(
				"",
				errors.New("some error"),
				nil,
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
			spy.instanceTypesForRegion = []string{"tralala"}
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

func TestService_CanInstanceBeTerminated(t *testing.T) {
	t.Run("instance can be terminated", func(t *testing.T) {
		instanceDetails := generateInstanceDetails()
		instanceDetails.State = publicCloud.STATE_UNKNOWN
		instanceDetails.Id = "instanceId"

		instanceDetailsById := make(map[string]*publicCloud.InstanceDetails)
		instanceDetailsById[instanceDetails.Id] = &instanceDetails

		spy := newRepositorySpy()
		spy.instanceDetailsById = instanceDetailsById

		service := New(&spy)
		got, reason, err := service.CanInstanceBeTerminated(
			"instanceId",
			context.TODO(),
		)

		assert.Nil(t, err)
		assert.Nil(t, reason)

		assert.True(t, got)
	})

	t.Run(
		"instance cannot be terminated if state is CREATING/DESTROYING/DESTROYED",
		func(t *testing.T) {
			tests := []struct {
				name           string
				state          publicCloud.State
				reasonContains string
			}{
				{
					name:           "state is CREATING",
					state:          publicCloud.STATE_CREATING,
					reasonContains: "CREATING",
				},
				{
					name:           "state is DESTROYING",
					state:          publicCloud.STATE_DESTROYING,
					reasonContains: "DESTROYING",
				},
				{
					name:           "state is DESTROYED",
					state:          publicCloud.STATE_DESTROYED,
					reasonContains: "DESTROYED",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					instanceDetails := generateInstanceDetails()
					instanceDetails.State = tt.state
					instanceDetails.Id = "instanceId"

					instanceDetailsById := make(map[string]*publicCloud.InstanceDetails)
					instanceDetailsById[instanceDetails.Id] = &instanceDetails

					spy := newRepositorySpy()
					spy.instanceDetailsById = instanceDetailsById

					service := New(&spy)
					got, reason, err := service.CanInstanceBeTerminated(
						"instanceId",
						context.TODO(),
					)

					assert.Nil(t, err)
					assert.NotNil(t, reason)
					assert.Contains(t, *reason, tt.reasonContains)

					assert.False(t, got)
				})
			}
		},
	)

	t.Run(
		"instance cannot be terminated if contract.endsAt is set",
		func(t *testing.T) {
			endsAt, _ := time.Parse(
				"2006-01-02 15:04:05",
				"2023-12-14 17:09:47",
			)

			instanceDetails := generateInstanceDetails()
			instanceDetails.State = publicCloud.STATE_UNKNOWN
			instanceDetails.Id = "instanceId"
			instanceDetails.Contract.EndsAt = *publicCloud.NewNullableTime(&endsAt)

			instanceDetailsById := make(map[string]*publicCloud.InstanceDetails)
			instanceDetailsById[instanceDetails.Id] = &instanceDetails

			spy := newRepositorySpy()
			spy.instanceDetailsById = instanceDetailsById

			service := New(&spy)
			got, reason, err := service.CanInstanceBeTerminated(
				"instanceId",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.NotNil(t, reason)
			assert.Contains(t, *reason, "2023-12-14 17:09:47 +0000 UTC")

			assert.False(t, got)
		},
	)

	t.Run("error from getSdkError bubbles up", func(t *testing.T) {
		service := New(
			&repositorySpy{
				getInstanceError: repository.NewSdkError(
					"",
					errors.New("some error"),
					nil,
				),
			},
		)

		_, _, err := service.CanInstanceBeTerminated("id", context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	},
	)
}

func TestService_GetBillingFrequencies(t *testing.T) {
	service := Service{}
	want := shared2.NewIntMarkdownList([]int{0, 1, 3, 6, 12})
	got := service.GetBillingFrequencies()

	assert.Equal(t, want, got)
}

func TestService_GetContractTerms(t *testing.T) {
	service := Service{}
	want := shared2.NewIntMarkdownList([]int{0, 1, 3, 6, 12})
	got := service.GetContractTerms()

	assert.Equal(t, want, got)
}

func TestService_GetContractTypes(t *testing.T) {
	service := Service{}
	want := []string{"HOURLY", "MONTHLY"}
	got := service.GetContractTypes()

	assert.Equal(t, want, got)
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

func TestService_GetMinimumRootDiskSize(t *testing.T) {
	service := Service{}
	want := int64(5)
	got := service.GetMinimumRootDiskSize()

	assert.Equal(t, want, got)
}

func TestService_GetMaximumRootDiskSize(t *testing.T) {
	service := Service{}
	got := service.GetMaximumRootDiskSize()

	assert.Equal(t, int64(1000), got)
}

func TestService_GetRootDiskStorageTypes(t *testing.T) {
	service := Service{}
	got := service.GetRootDiskStorageTypes()

	assert.Contains(t, got, "CENTRAL")
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

func TestService_IsInstanceTypeAvailableForRegion(t *testing.T) {
	t.Run(
		"return true when instanceType is available for region",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.instanceTypesForRegion = []string{"tralala"}

			service := New(&spy)

			got, instanceTypes, err := service.IsInstanceTypeAvailableForRegion(
				"tralala",
				"region",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"tralala"}, instanceTypes)
			assert.True(t, got)
		},
	)

	t.Run(
		"return true when instanceType is not available for region",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.instanceTypesForRegion = []string{"piet"}

			service := New(&spy)

			got, instanceTypes, err := service.IsInstanceTypeAvailableForRegion(
				"tralala",
				"region",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"piet"}, instanceTypes)
			assert.False(t, got)
		},
	)

	t.Run("region is passed to repository", func(t *testing.T) {
		spy := newRepositorySpy()
		service := New(&spy)

		_, _, _ = service.IsInstanceTypeAvailableForRegion(
			"tralala",
			"region",
			context.TODO(),
		)

		assert.Equal(
			t,
			"region",
			spy.passedGetInstanceTypesForRegionRegion,
		)
	})

	t.Run("errors from service bubble up", func(t *testing.T) {
		spy := newRepositorySpy()
		spy.getInstanceTypesForRegionError = repository.NewSdkError(
			"",
			errors.New("some error"),
			nil,
		)
		service := New(&spy)

		_, _, err := service.IsInstanceTypeAvailableForRegion(
			"tralala",
			"region",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_CanInstanceTypeBeUsedWithInstance(t *testing.T) {
	t.Run(
		"returns true if instanceType is equal to the current instanceType",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.availableInstanceTypesForUpdate = []string{}

			service := New(&spy)

			got, instanceTypes, err := service.CanInstanceTypeBeUsedWithInstance(
				"085075b0-a6ad-4026-a0d1-e3256d3f7c47",
				"tralala",
				"tralala",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"tralala"}, instanceTypes)
			assert.True(t, got)
		},
	)
	t.Run(
		"returns true if instanceType is in availableInstanceTypes",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.availableInstanceTypesForUpdate = []string{"tralala"}

			service := New(&spy)

			got, instanceTypes, err := service.CanInstanceTypeBeUsedWithInstance(
				"085075b0-a6ad-4026-a0d1-e3256d3f7c47",
				"",
				"tralala",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"tralala", ""}, instanceTypes)
			assert.True(t, got)
		},
	)

	t.Run(
		"returns false if instanceType cannot be used with instance",
		func(t *testing.T) {
			spy := newRepositorySpy()
			spy.availableInstanceTypesForUpdate = []string{"piet"}

			service := New(&spy)

			got, instanceTypes, err := service.CanInstanceTypeBeUsedWithInstance(
				"085075b0-a6ad-4026-a0d1-e3256d3f7c47",
				"",
				"tralala",
				context.TODO(),
			)

			assert.Nil(t, err)
			assert.Equal(t, []string{"piet", ""}, instanceTypes)
			assert.False(t, got)
		},
	)

	t.Run("errors from the service bubble up", func(t *testing.T) {
		spy := newRepositorySpy()
		spy.getAvailableInstanceTypesForUpdateError = repository.NewSdkError(
			"",
			errors.New("some error"),
			nil,
		)

		service := New(&spy)

		_, _, err := service.CanInstanceTypeBeUsedWithInstance(
			"3cf0ddcb-b375-45a8-b18a-1bdad52527f2",
			"",
			"",
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("id is passed to repository", func(t *testing.T) {
		want := "085075b0-a6ad-4026-a0d1-e3256d3f7c47"

		spy := newRepositorySpy()

		service := New(&spy)

		_, _, _ = service.CanInstanceTypeBeUsedWithInstance(
			want,
			"",
			"",
			context.TODO(),
		)

		assert.Equal(t, want, spy.passedGetAvailableInstanceTypesForUpdateId)
	})
}

func generateInstanceDetails() publicCloud.InstanceDetails {
	return publicCloud.InstanceDetails{
		Id:     "id",
		Image:  publicCloud.Image{Id: "imageId"},
		Type:   "instanceType",
		Region: "region",
	}
}

func generateInstanceModel() resource.Instance {
	image, _ := types.ObjectValueFrom(
		context.TODO(),
		resource.Image{}.AttributeTypes(),
		resource.Image{
			Id: basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
		},
	)

	contract, _ := types.ObjectValueFrom(
		context.TODO(),
		resource.Contract{}.AttributeTypes(),
		resource.Contract{
			BillingFrequency: basetypes.NewInt64Value(int64(1)),
			Term:             basetypes.NewInt64Value(int64(3)),
			Type:             basetypes.NewStringValue("MONTHLY"),
			State:            basetypes.NewStringUnknown(),
		},
	)

	return resource.Instance{
		Id:                  basetypes.NewStringValue("id"),
		Region:              basetypes.NewStringValue("eu-west-3"),
		Type:                basetypes.NewStringValue("lsw.m5a.4xlarge"),
		RootDiskStorageType: basetypes.NewStringValue("CENTRAL"),
		RootDiskSize:        basetypes.NewInt64Value(int64(55)),
		Image:               image,
		Contract:            contract,
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
	}
}