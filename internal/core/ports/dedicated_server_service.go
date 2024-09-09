package ports

import (
	"context"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
)

// DedicatedServerService gets data associated with dedicated_server.
type DedicatedServerService interface {
	// GetAllDedicatedServers gets dedicated_server.DedicatedServers.
	GetAllDedicatedServers(ctx context.Context) (domain.DedicatedServers, *errors.ServiceError)

	// GetAllOperatingSystems gets dedicated_server.OperatingSystems.
	GetAllOperatingSystems(ctx context.Context) (domain.OperatingSystems, *errors.ServiceError)

	// GetAllControlPanels gets dedicated_server.ControlPanels.
	GetAllControlPanels(ctx context.Context) (domain.ControlPanels, *errors.ServiceError)

	// CreateNotificationSettingBandwidth creates dedicated_server.NotificationSettingBandwidth.
	CreateNotificationSettingBandwidth(
		notificationSettingBandwidth domain.NotificationSettingBandwidth,
		ctx context.Context,
	) (*domain.NotificationSettingBandwidth, *errors.ServiceError)
}
