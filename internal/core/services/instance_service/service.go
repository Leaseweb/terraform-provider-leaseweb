package instance_service

import (
	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type Service struct{}

func (srv Service) GetAllInstances() (entity.Instances, error) {
	return entity.Instances{}, nil
}

func (srv Service) GetInstance(id uuid.UUID) (entity.Instance, error) {
	return entity.Instance{}, nil
}

func (srv Service) CreateInstance(instance entity.Instance) (
	entity.Instance,
	error,
) {
	return entity.Instance{}, nil
}

func (srv Service) UpdateInstance(instance entity.Instance) (
	entity.Instance,
	error,
) {
	return entity.Instance{}, nil
}

func (srv Service) DeleteInstance(id uuid.UUID) error {
	return nil
}

func New() Service {
	return Service{}
}
