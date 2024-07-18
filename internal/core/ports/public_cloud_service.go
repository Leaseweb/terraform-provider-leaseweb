package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/services/shared"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type PublicCloudService interface {
	GetAllInstances(ctx context.Context) (domain.Instances, *shared.ServiceError)

	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	DeleteInstance(id value_object.Uuid, ctx context.Context) *shared.ServiceError

	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, *shared.ServiceError)

	GetRegions(
		ctx context.Context,
	) (domain.Regions, *shared.ServiceError)
}
