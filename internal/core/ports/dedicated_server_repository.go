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

	// GetDedicatedServer returns data for a singular dedicated_server from the dedicated server api.
	GetDedicatedServer(ctx context.Context, id string) (
		*domain.DedicatedServer,
		*shared.RepositoryError,
	)

	// GetDataTrafficNotificationSetting returns data for a singular data traffic notification setting from the dedicated server api.
	GetDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) (
		*domain.DataTrafficNotificationSetting,
		*shared.RepositoryError,
	)

	// CreateDataTrafficNotificationSetting creates a new data traffic notification service in the dedicated server api.
	CreateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (
		*domain.DataTrafficNotificationSetting,
		*shared.RepositoryError,
	)

	// UpdateDataTrafficNotificationSetting updates a data traffic notification service in the dedicated server api.
	UpdateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (
		*domain.DataTrafficNotificationSetting,
		*shared.RepositoryError,
	)

	// DeleteDataTrafficNotificationSetting deletes a data traffic notification setting in the dedicated server api.
	DeleteDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) *shared.RepositoryError
}
