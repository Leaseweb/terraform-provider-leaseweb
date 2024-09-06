// Package dedicated_server_repository implements repository logic
// to access the dedicated_server sdk.
package dedicated_server_repository

import (
	"context"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/dedicated_server_repository/data_adapters/to_sdk_model"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/dedicated_server_repository/data_adapters/to_domain_entity"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/sdk"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// Optional contains optional values that can be passed to NewDedicatedServerRepository.
type Optional struct {
	Host   *string
	Scheme *string
}

// DedicatedServerRepository fulfills contract for ports.DedicatedServerRepository.
type DedicatedServerRepository struct {
	dedicatedServerApi                      sdk.DedicatedServerApi
	token                                   string
	adaptDedicatedServer                    func(sdkDedicatedServer dedicatedServer.Server) domain.DedicatedServer
	adaptOperatingSystems                   func(sdkOperatingSystem []dedicatedServer.OperatingSystem) domain.OperatingSystems
	adaptControlPanels                      func(sdkControlPanel []dedicatedServer.ControlPanel) domain.ControlPanels
	adaptDataTrafficNotificationSetting     func(sdkDataTrafficNotificationSetting dedicatedServer.DataTrafficNotificationSetting) domain.DataTrafficNotificationSetting
	adaptDataTrafficNotificationSettingOpts func(dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) dedicatedServer.DataTrafficNotificationSettingOpts
}

// Injects the authentication token into the context for the sdk.
func (r DedicatedServerRepository) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: r.token, Prefix: ""},
		},
	)
}

func (r DedicatedServerRepository) GetAllDedicatedServers(ctx context.Context) (
	domain.DedicatedServers,
	*shared.RepositoryError,
) {
	var dedicatedServers domain.DedicatedServers

	request := r.dedicatedServerApi.GetServerList(r.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetAllDedicatedServers", err, response)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		50, // metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError("GetAllDedicatedServers", err, response)
		}

		for _, sdkDedicatedServer := range result.GetServers() {
			dedicatedServers = append(dedicatedServers, r.adaptDedicatedServer(sdkDedicatedServer))
		}

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllDedicatedServers", err, response)
		}
	}

	return dedicatedServers, nil
}

func (r DedicatedServerRepository) GetAllOperatingSystems(ctx context.Context) (
	domain.OperatingSystems,
	*shared.RepositoryError,
) {
	var operatingSystems domain.OperatingSystems

	request := r.dedicatedServerApi.GetOperatingSystemList(r.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetAllOperatingSystems", err, response)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError("GetAllOperatingSystems", err, response)
		}

		operatingSystems = append(operatingSystems, r.adaptOperatingSystems(result.GetOperatingSystems())...)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllOperatingSystems", err, response)
		}
	}

	return operatingSystems, nil
}

func (d DedicatedServerRepository) GetAllControlPanels(ctx context.Context) (
	domain.ControlPanels,
	*shared.RepositoryError,
) {
	var controlPanels domain.ControlPanels

	request := d.dedicatedServerApi.GetControlPanelList(d.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetAllControlPanels", err, response)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError("GetAllControlPanels", err, response)
		}

		controlPanels = append(controlPanels, d.adaptControlPanels(result.GetControlPanels())...)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllControlPanels", err, response)
		}
	}

	return controlPanels, nil
}

func (p DedicatedServerRepository) GetDedicatedServer(ctx context.Context, id string) (
	*domain.DedicatedServer,
	*shared.RepositoryError,
) {
	request := p.dedicatedServerApi.GetServer(p.authContext(ctx), id)
	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetDedicatedServer", err, response)
	}
	dedicatedServer := p.adaptDedicatedServer(*result)
	return &dedicatedServer, nil
}

func (p DedicatedServerRepository) GetDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) (
	*domain.DataTrafficNotificationSetting,
	*shared.RepositoryError,
) {
	request := p.dedicatedServerApi.GetServerDataTrafficNotificationSetting(p.authContext(ctx), serverId, dataTrafficNotificationSettingId)
	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetServerDataTrafficNotificationSetting", err, response)
	}
	dataTrafficNotificationSetting := p.adaptDataTrafficNotificationSetting(*result)
	return &dataTrafficNotificationSetting, nil
}

func (p DedicatedServerRepository) CreateDataTrafficNotificationSetting(
	ctx context.Context,
	serverId string,
	dataTrafficNotificationSetting domain.DataTrafficNotificationSetting,
) (
	*domain.DataTrafficNotificationSetting,
	*shared.RepositoryError,
) {
	dataTrafficNotificationSettingOpts := p.adaptDataTrafficNotificationSettingOpts(dataTrafficNotificationSetting)
	request := p.dedicatedServerApi.CreateServerDataTrafficNotificationSetting(p.authContext(ctx), serverId).DataTrafficNotificationSettingOpts(dataTrafficNotificationSettingOpts)
	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("CreateDataTrafficNotificationSetting", err, response)
	}
	createdDataTrafficNotificationSetting := p.adaptDataTrafficNotificationSetting(*result)
	return &createdDataTrafficNotificationSetting, nil
}

func (p DedicatedServerRepository) UpdateDataTrafficNotificationSetting(
	ctx context.Context,
	serverId string,
	dataTrafficNotificationSettingId string,
	dataTrafficNotificationSetting domain.DataTrafficNotificationSetting,
) (
	*domain.DataTrafficNotificationSetting,
	*shared.RepositoryError,
) {
	dataTrafficNotificationSettingOpts := p.adaptDataTrafficNotificationSettingOpts(dataTrafficNotificationSetting)
	request := p.dedicatedServerApi.UpdateServerDataTrafficNotificationSetting(p.authContext(ctx), serverId, dataTrafficNotificationSettingId).DataTrafficNotificationSettingOpts(dataTrafficNotificationSettingOpts)
	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("UpdateDataTrafficNotificationSetting", err, response)
	}
	// TODO: check if we improve
	updatedDataTrafficNotificationSetting := domain.NewUpdateDataTrafficNotificationSetting(
		result.GetFrequency(),
		result.GetThreshold(),
		result.GetUnit(),
	)
	return &updatedDataTrafficNotificationSetting, nil
}

func (p DedicatedServerRepository) DeleteDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) *shared.RepositoryError {
	request := p.dedicatedServerApi.DeleteServerDataTrafficNotificationSetting(p.authContext(ctx), serverId, dataTrafficNotificationSettingId)
	response, err := request.Execute()

	if err != nil {
		return shared.NewSdkError("DeleteDataTrafficNotificationSetting", err, response)
	}
	return nil
}

func NewDedicatedServerRepository(
	token string,
	optional Optional,
) DedicatedServerRepository {
	configuration := dedicatedServer.NewConfiguration()

	if optional.Host != nil {
		configuration.Host = *optional.Host
	}
	if optional.Scheme != nil {
		configuration.Scheme = *optional.Scheme
	}

	client := *dedicatedServer.NewAPIClient(configuration)

	return DedicatedServerRepository{
		dedicatedServerApi:                      client.DedicatedServerAPI,
		token:                                   token,
		adaptDedicatedServer:                    to_domain_entity.AdaptDedicatedServer,
		adaptOperatingSystems:                   to_domain_entity.AdaptOperatingSystems,
		adaptControlPanels:                      to_domain_entity.AdaptControlPanels,
		adaptDataTrafficNotificationSetting:     to_domain_entity.AdaptDataTrafficNotificationSetting,
		adaptDataTrafficNotificationSettingOpts: to_sdk_model.AdaptToDataTrafficNotificationSettingOpts,
	}
}
