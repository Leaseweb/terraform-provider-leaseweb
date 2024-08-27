package ports

import (
	"context"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// DedicatedServerRepository is used to connect to dedicated_server api.
type DedicatedServerRepository interface {
	// GetAllDedicatedServers retrieve all dedicated_servers from the dedicated server api.
	GetAllDedicatedServers(ctx context.Context) (
		domain.DedicatedServers,
		*shared.RepositoryError,
	)

	// GetAllOperatingSystems retrieve all operating systems from the dedicated server api.
	GetAllOperatingSystems(ctx context.Context) (
		domain.OperatingSystems,
		*shared.RepositoryError,
	)
}
