package instance_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type repositorySpy struct {
	Error     error
	Instances entity.Instances
	Instance  *entity.Instance
}

func newRepositorySpy(
	hasError bool,
	instances entity.Instances,
	instance *entity.Instance,
) repositorySpy {
	var returnedError error = nil

	if hasError {
		returnedError = errors.New("some error")
	}

	return repositorySpy{
		Error:     returnedError,
		Instances: instances,
		Instance:  instance,
	}
}

func (r repositorySpy) GetAllInstances(ctx context.Context) (
	entity.Instances,
	error,
) {
	if r.Error != nil {
		return nil, r.Error
	}

	return r.Instances, nil
}

func (r repositorySpy) GetInstance(
	id uuid.UUID,
	ctx context.Context,
) (*entity.Instance, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	r.Instance.Id = id
	return r.Instance, nil
}

func (r repositorySpy) CreateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	return r.Instance, nil
}

func (r repositorySpy) UpdateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	return r.Instance, nil
}

func (r repositorySpy) DeleteInstance(
	id uuid.UUID,
	ctx context.Context,
) error {
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func TestService_GetAllInstances(t *testing.T) {
	t.Run(
		"service passes back instances from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(newRepositorySpy(
				false,
				entity.Instances{{Id: id}},
				nil,
			))

			got, err := instanceService.GetAllInstances(context.TODO())

			assert.NoError(t, err)
			assert.Len(t, got, 1)
			assert.Equal(t, id, got[0].Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(newRepositorySpy(
			true,
			nil,
			nil,
		))

		_, err := instanceService.GetAllInstances(context.TODO())

		assert.Error(t, err)
	})
}

func TestService_GetInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(newRepositorySpy(
				false,
				nil,
				&entity.Instance{Id: id},
			))

			got, err := instanceService.GetInstance(id, context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(newRepositorySpy(
			true,
			nil,
			nil,
		))

		_, err := instanceService.GetInstance(uuid.New(), context.TODO())

		assert.Error(t, err)
	})
}

func TestService_CreateInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(newRepositorySpy(
				false,
				nil,
				&entity.Instance{Id: id},
			))

			got, err := instanceService.CreateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(newRepositorySpy(
			true,
			nil,
			nil,
		))

		_, err := instanceService.CreateInstance(entity.Instance{}, context.TODO())

		assert.Error(t, err)
	})
}

func TestService_UpdateInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(newRepositorySpy(
				false,
				nil,
				&entity.Instance{Id: id},
			))

			got, err := instanceService.UpdateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(newRepositorySpy(
			true,
			nil,
			nil,
		))

		_, err := instanceService.UpdateInstance(
			entity.Instance{},
			context.TODO(),
		)

		assert.Error(t, err)
	})
}

func TestService_DeleteInstance(t *testing.T) {
	t.Run("service passes back nil from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{})

		err := instanceService.DeleteInstance(uuid.New(), context.TODO())

		assert.NoError(t, err)
	})

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(newRepositorySpy(
			true,
			nil,
			nil,
		))

		err := instanceService.DeleteInstance(uuid.New(), context.TODO())

		assert.Error(t, err)
	})
}
