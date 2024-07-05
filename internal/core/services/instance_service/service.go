package instance_service

import (
	"errors"

	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/ports"
)

var ErrFailedToGetInstancesFromRepository = errors.New("failed to retrieve instances from repository")
var ErrFailedToGetInstanceFromRepository = errors.New("failed to retrieve instance from repository")
var ErrFailedToCreateInstance = errors.New("failed to create instance")
var ErrFailedToUpdateInstance = errors.New("failed to update instance")
var ErrFailedToDeleteInstance = errors.New("failed to delete instance")

type Service struct {
	instanceRepository ports.InstanceRepository
}

func (srv Service) GetAllInstances() (entity.Instances, error) {
	instances, err := srv.instanceRepository.GetAllInstances()
	if err != nil {
		return entity.Instances{}, ErrFailedToGetInstancesFromRepository
	}

	return instances, nil
}

func (srv Service) GetInstance(id uuid.UUID) (entity.Instance, error) {
	instance, err := srv.instanceRepository.GetInstance(id)
	if err != nil {
		return entity.Instance{}, ErrFailedToGetInstanceFromRepository
	}

	return instance, nil
}

func (srv Service) CreateInstance(instance entity.Instance) (
	entity.Instance,
	error,
) {
	instance, err := srv.instanceRepository.CreateInstance(instance)
	if err != nil {
		return entity.Instance{}, ErrFailedToCreateInstance
	}

	return instance, nil
}

func (srv Service) UpdateInstance(instance entity.Instance) (
	entity.Instance,
	error,
) {
	instance, err := srv.instanceRepository.UpdateInstance(instance)
	if err != nil {
		return entity.Instance{}, ErrFailedToUpdateInstance
	}

	return instance, nil
}

func (srv Service) DeleteInstance(id uuid.UUID) error {
	err := srv.instanceRepository.DeleteInstance(id)
	if err != nil {
		return ErrFailedToDeleteInstance
	}

	return nil
}

func New(instanceRepository ports.InstanceRepository) Service {
	return Service{instanceRepository: instanceRepository}
}
