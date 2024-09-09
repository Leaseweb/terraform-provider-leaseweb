// Package dedicated_server implements services related to dedicated_servers
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
	domain.DedicatedServers,
	*errors.ServiceError,
) {

	dedicatedServers, err := s.dedicatedServerRepository.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAllDedicatedServers",
			*err,
		)
	}

	return dedicatedServers, nil
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

func (s Service) GetAllControlPanels(ctx context.Context) (
	domain.ControlPanels,
	*errors.ServiceError,
) {

	controlPanels, err := s.dedicatedServerRepository.GetAllControlPanels(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAllControlPanels",
			*err,
		)
	}

	return controlPanels, nil
}

func (srv Service) GetDedicatedServer(ctx context.Context, id string) (
	*domain.DedicatedServer,
	*errors.ServiceError,
) {

	dedicatedServer, err := srv.dedicatedServerRepository.GetDedicatedServer(ctx, id)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetDedicatedServer",
			*err,
		)
	}

	return dedicatedServer, nil
}

func New(dedicatedServerRepository ports.DedicatedServerRepository) Service {
	return Service{dedicatedServerRepository: dedicatedServerRepository}
}
