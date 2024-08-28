package dedicated_server_repository

import (
	"context"

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
	dedicatedServerApi    sdk.DedicatedServerApi
	token                 string
	adaptDedicatedServer  func(sdkDedicatedServer dedicatedServer.Server) domain.DedicatedServer
	adaptOperatingSystems func(sdkOperatingSystem []dedicatedServer.OperatingSystem) domain.OperatingSystems
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
		dedicatedServerApi:    client.DedicatedServerAPI,
		token:                 token,
		adaptDedicatedServer:  to_domain_entity.AdaptDedicatedServer,
		adaptOperatingSystems: to_domain_entity.AdaptOperatingSystems,
	}
}
