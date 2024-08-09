package ports

import (
	"context"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
)

// DedicatedServerService gets data associated with dedicated_server.
type DedicatedServerService interface {
	// GetAllDedicatedServers gets all dedicated servers.
	GetAllDedicatedServers(ctx context.Context) (domain.DedicatedServers, *errors.ServiceError)
}
