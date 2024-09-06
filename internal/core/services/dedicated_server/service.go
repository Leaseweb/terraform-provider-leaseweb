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

func (srv Service) GetDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) (
	*domain.DataTrafficNotificationSetting,
	*errors.ServiceError,
) {

	dataTrafficNotificationSetting, err := srv.dedicatedServerRepository.GetDataTrafficNotificationSetting(ctx, serverId, dataTrafficNotificationSettingId)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetDataTrafficNotificationSetting",
			*err,
		)
	}

	return dataTrafficNotificationSetting, nil
}

func (srv Service) CreateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (
	*domain.DataTrafficNotificationSetting,
	*errors.ServiceError,
) {

	createdDataTrafficNotificationSetting, err := srv.dedicatedServerRepository.CreateDataTrafficNotificationSetting(ctx, serverId, dataTrafficNotificationSetting)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"CreateDataTrafficNotificationSetting",
			*err,
		)
	}

	return createdDataTrafficNotificationSetting, nil
}

func (srv Service) UpdateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (
	*domain.DataTrafficNotificationSetting,
	*errors.ServiceError,
) {
	updatedDataTrafficNotificationSetting, err := srv.dedicatedServerRepository.UpdateDataTrafficNotificationSetting(ctx, serverId, dataTrafficNotificationSettingId, dataTrafficNotificationSetting)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"UpdateDataTrafficNotificationSetting",
			*err,
		)
	}
	return updatedDataTrafficNotificationSetting, nil
}

func (srv Service) DeleteDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) *errors.ServiceError {

	err := srv.dedicatedServerRepository.DeleteDataTrafficNotificationSetting(ctx, serverId, dataTrafficNotificationSettingId)
	if err != nil {
		return errors.NewFromRepositoryError(
			"DeleteDataTrafficNotificationSetting",
			*err,
		)
	}
	return nil
}

func New(dedicatedServerRepository ports.DedicatedServerRepository) Service {
	return Service{dedicatedServerRepository: dedicatedServerRepository}
}
