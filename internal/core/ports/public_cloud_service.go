package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type PublicCloudService interface {
	GetAllInstances(ctx context.Context) (entity.Instances, error)

	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*entity.Instance, error)

	CreateInstance(
		instance entity.Instance,
		ctx context.Context,
	) (*entity.Instance, error)

	UpdateInstance(
		instance entity.Instance,
		ctx context.Context,
	) (*entity.Instance, error)

	DeleteInstance(id value_object.Uuid, ctx context.Context) error

	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (entity.InstanceTypes, error)

	GetRegions(
		ctx context.Context,
	) (entity.Regions, error)
}
