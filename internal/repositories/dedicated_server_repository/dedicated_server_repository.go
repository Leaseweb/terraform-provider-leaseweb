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
	dedicatedServerApi   sdk.DedicatedServerApi
	token                string
	adaptDedicatedServer func(sdkDedicatedServer dedicatedServer.Server) domain.DedicatedServer
}

// Injects the authentication token into the context for the sdk.
func (p DedicatedServerRepository) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: p.token, Prefix: ""},
		},
	)
}

func (p DedicatedServerRepository) GetAllDedicatedServers(ctx context.Context) (
	domain.DedicatedServers,
	*shared.RepositoryError,
) {
	var dedicatedServers domain.DedicatedServers

	request := p.dedicatedServerApi.GetServerList(p.authContext(ctx))

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
			dedicatedServers = append(dedicatedServers, p.adaptDedicatedServer(sdkDedicatedServer))
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
		dedicatedServerApi:   client.DedicatedServerAPI,
		token:                token,
		adaptDedicatedServer: to_domain_entity.AdaptDedicatedServer,
	}
}
