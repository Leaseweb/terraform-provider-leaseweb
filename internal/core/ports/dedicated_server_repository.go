package ports

import (
	"context"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// DedicatedServerRepository is used to connect to dedicated_server api.
type DedicatedServerRepository interface {
	// GetAllDedicatedServers retrieves dedicated_server.DedicatedServers from the dedicated_server api.
	GetAllDedicatedServers(ctx context.Context) (
		domain.DedicatedServers,
		*shared.RepositoryError,
	)

	// GetAllOperatingSystems retrieves dedicated_server.OperatingSystems from the dedicated_server api.
	GetAllOperatingSystems(ctx context.Context) (
		domain.OperatingSystems,
		*shared.RepositoryError,
	)

	// GetAllControlPanels retrieves dedicated_server.ControlPanels from the dedicated_server api.
	GetAllControlPanels(ctx context.Context) (
		domain.ControlPanels,
		*shared.RepositoryError,
	)
}
