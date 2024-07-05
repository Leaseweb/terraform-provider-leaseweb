package ports

import (
	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type InstanceService interface {
	GetAllInstances() (entity.Instances, error)

	GetInstance(id uuid.UUID) (entity.Instance, error)

	CreateInstance(instance entity.Instance) (entity.Instance, error)

	UpdateInstance(instance entity.Instance) (entity.Instance, error)

	DeleteInstance(id uuid.UUID) error
}
