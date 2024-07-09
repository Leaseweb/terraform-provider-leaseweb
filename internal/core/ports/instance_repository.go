package ports

import (
	"context"

	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type InstanceRepository interface {
	GetAllInstances(ctx context.Context) (entity.Instances, error)

	GetInstance(id uuid.UUID, ctx context.Context) (*entity.Instance, error)

	CreateInstance(
		instance entity.Instance,
		ctx context.Context,
	) (*entity.Instance, error)

	UpdateInstance(
		instance entity.Instance,
		ctx context.Context,
	) (*entity.Instance, error)

	DeleteInstance(id uuid.UUID, ctx context.Context) error
}
