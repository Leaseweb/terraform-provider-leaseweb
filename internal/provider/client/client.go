package client

import (
	dedicatedserverservice "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/dedicated_server"
	publiccloudservice "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/dedicated_server_repository"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/public_cloud_repository"
)

// The Client handles instantiation of the facades.
type Client struct {
	PublicCloudFacade     public_cloud.PublicCloudFacade
	DedicatedServerFacade dedicated_server.DedicatedServerFacade
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	publicCloudRepository := public_cloud_repository.NewPublicCloudRepository(
		token,
		public_cloud_repository.Optional{
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
	)
	publicCloudService := publiccloudservice.New(publicCloudRepository)

	dedicatedServerRepository := dedicated_server_repository.NewDedicatedServerRepository(
		token,
		dedicated_server_repository.Optional{
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
	)
	dedicatedServerService := dedicatedserverservice.New(dedicatedServerRepository)

	return Client{
		PublicCloudFacade:     public_cloud.NewPublicCloudFacade(&publicCloudService),
		DedicatedServerFacade: dedicated_server.New(dedicatedServerService),
	}
}
