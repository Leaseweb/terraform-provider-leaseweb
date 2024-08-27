package ports

import (
	"context"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// DedicatedServerRepository is used to connect to dedicated_server api.
type DedicatedServerRepository interface {
	// GetAllDedicatedServers retrieves all dedicated_servers from the dedicated server api.
	GetAllDedicatedServers(ctx context.Context) (
		domain.DedicatedServers,
		*shared.RepositoryError,
	)

	// GetAllOperatingSystems retrieves all operating systems from the dedicated server api.
	GetAllOperatingSystems(ctx context.Context) (
		domain.OperatingSystems,
		*shared.RepositoryError,
	)
}
