package instance_service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type repositorySpy struct {
	Error     error
	Instances entity.Instances
	Instance  entity.Instance
}

func (r repositorySpy) GetAllInstances() (entity.Instances, error) {
	if r.Error != nil {
		return entity.Instances{}, r.Error
	}

	return r.Instances, nil
}

func (r repositorySpy) GetInstance(id uuid.UUID) (entity.Instance, error) {
	if r.Error != nil {
		return entity.Instance{}, r.Error
	}

	r.Instance.Id = id
	return r.Instance, nil
}

func (r repositorySpy) CreateInstance(instance entity.Instance) (
	entity.Instance,
	error,
) {
	if r.Error != nil {
		return entity.Instance{}, r.Error
	}

	return r.Instance, nil
}

func (r repositorySpy) UpdateInstance(instance entity.Instance) (
	entity.Instance,
	error,
) {
	if r.Error != nil {
		return entity.Instance{}, r.Error
	}

	return r.Instance, nil
}

func (r repositorySpy) DeleteInstance(id uuid.UUID) error {
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
			instanceService := New(repositorySpy{Instances: entity.Instances{{Id: id}}})

			got, err := instanceService.GetAllInstances()

			assert.Nil(t, err)
			assert.Equal(t, id, got[0].Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{Error: ErrFailedToUpdateInstance})

		_, err := instanceService.GetAllInstances()

		assert.NotNil(t, err)
	})
}

func TestService_GetInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{})

			got, err := instanceService.GetInstance(id)

			assert.Nil(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{Error: ErrFailedToUpdateInstance})

		_, err := instanceService.GetInstance(uuid.New())

		assert.NotNil(t, err)
	})
}

func TestService_CreateInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{Instance: entity.Instance{Id: id}})

			got, err := instanceService.CreateInstance(entity.Instance{})

			assert.Nil(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{Error: ErrFailedToUpdateInstance})

		_, err := instanceService.CreateInstance(entity.Instance{})

		assert.NotNil(t, err)
	})
}

func TestService_UpdateInstance(t *testing.T) {
	t.Run(
		"service passes back instance from repository",
		func(t *testing.T) {
			id := uuid.New()
			instanceService := New(repositorySpy{Instance: entity.Instance{Id: id}})

			got, err := instanceService.UpdateInstance(entity.Instance{})

			assert.Nil(t, err)
			assert.Equal(t, id, got.Id)
		},
	)

	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{Error: ErrFailedToUpdateInstance})

		_, err := instanceService.UpdateInstance(entity.Instance{})

		assert.NotNil(t, err)
	})
}

func TestService_DeleteInstance(t *testing.T) {
	t.Run("service passes back error from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{Error: ErrFailedToUpdateInstance})

		err := instanceService.DeleteInstance(uuid.New())

		assert.NotNil(t, err)
	})
	t.Run("service passes back nil from repository", func(t *testing.T) {
		instanceService := New(repositorySpy{})

		err := instanceService.DeleteInstance(uuid.New())

		assert.Nil(t, err)
	})
}
