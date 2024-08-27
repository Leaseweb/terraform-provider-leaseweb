package dedicated_server

import (
	"context"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
)

// Service fulfills the contract for ports.DedicatedServerService.
type Service struct {
	dedicatedServerRepository ports.DedicatedServerRepository
}

func (s Service) GetAllDedicatedServers(ctx context.Context) (
	*domain.DedicatedServers,
	*errors.ServiceError,
) {

	dedicatedServers, err := s.dedicatedServerRepository.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAllDedicatedServers",
			*err,
		)
	}

	return &dedicatedServers, nil
}

func (s Service) GetAllOperatingSystems(ctx context.Context) (
	domain.OperatingSystems,
	*errors.ServiceError,
) {

	operatingSystems, err := s.dedicatedServerRepository.GetAllOperatingSystems(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAllOperatingSystems",
			*err,
		)
	}

	return operatingSystems, nil
}

func New(dedicatedServerRepository ports.DedicatedServerRepository) Service {
	return Service{dedicatedServerRepository: dedicatedServerRepository}
}
