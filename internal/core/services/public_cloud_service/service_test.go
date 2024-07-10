package public_cloud_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/ports"
)

var (
	_ ports.PublicCloudRepository = &repositorySpy{}
)

type repositorySpy struct {
	instances                entity.Instances
	instance                 *entity.Instance
	autoScalingGroup         *entity.AutoScalingGroup
	loadBalancer             *entity.LoadBalancer
	getAutoScalingGroupError error
	getLoadBalancerError     error
	getAllInstancesError     error
	getInstanceError         error
	createInstanceError      error
	updateInstanceError      error
	deleteInstanceError      error
}

func (r repositorySpy) GetAutoScalingGroup(
	id uuid.UUID,
	ctx context.Context,
) (*entity.AutoScalingGroup, error) {
	return r.autoScalingGroup, r.getAutoScalingGroupError
}

func (r repositorySpy) GetLoadBalancer(
	id uuid.UUID,
	ctx context.Context,
) (*entity.LoadBalancer, error) {
	return r.loadBalancer, r.getLoadBalancerError
}

func (r repositorySpy) GetAllInstances(ctx context.Context) (
	entity.Instances,
	error,
) {
	return r.instances, r.getAllInstancesError
}

func (r repositorySpy) GetInstance(
	id uuid.UUID,
	ctx context.Context,
) (*entity.Instance, error) {
	return r.instance, r.getInstanceError
}

func (r repositorySpy) CreateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	return r.instance, r.createInstanceError
}

func (r repositorySpy) UpdateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	return r.instance, r.updateInstanceError
}

func (r repositorySpy) DeleteInstance(
	id uuid.UUID,
	ctx context.Context,
) error {
	return r.deleteInstanceError
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{instances: entity.Instances{{Id: id}}})

			got, err := instanceService.GetAllInstances(context.TODO())

			assert.NoError(t, err)
			assert.Len(t, got, 1)
			assert.Equal(t, id, got[0].Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{getAllInstancesError: errors.New("some error")})

		_, err := instanceService.GetAllInstances(context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_GetInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{instance: &entity.Instance{Id: id}})

			got, err := instanceService.GetInstance(id, context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{getInstanceError: errors.New("some error")})

		_, err := instanceService.GetInstance(uuid.New(), context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_CreateInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{instance: &entity.Instance{Id: id}})

			got, err := instanceService.CreateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{createInstanceError: errors.New("some error")})

		_, err := instanceService.CreateInstance(entity.Instance{}, context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_UpdateInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{instance: &entity.Instance{Id: id}})

			got, err := instanceService.UpdateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{updateInstanceError: errors.New("some error")})

		_, err := instanceService.UpdateInstance(
			entity.Instance{},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}

func TestService_DeleteInstance(t *testing.T) {
	t.Run("service passes back nil from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{})

		err := instanceService.DeleteInstance(uuid.New(), context.TODO())

		assert.NoError(t, err)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{deleteInstanceError: errors.New("some error")})

		err := instanceService.DeleteInstance(uuid.New(), context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "some error")
	})
}
