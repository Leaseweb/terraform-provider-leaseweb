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

func (srv Service) GetAllDedicatedServers(ctx context.Context) (
	domain.DedicatedServers,
	*errors.ServiceError,
) {

	dedicatedServers, err := srv.dedicatedServerRepository.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAllDedicatedServers",
			*err,
		)
	}

	return dedicatedServers, nil
}

func (srv Service) GetAllControlPanels(ctx context.Context) (
	domain.ControlPanels,
	*errors.ServiceError,
) {

	controlPanels, err := srv.dedicatedServerRepository.GetAllControlPanels(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAllControlPanels",
			*err,
		)
	}

	return controlPanels, nil
}

func New(dedicatedServerRepository ports.DedicatedServerRepository) Service {
	return Service{dedicatedServerRepository: dedicatedServerRepository}
}
