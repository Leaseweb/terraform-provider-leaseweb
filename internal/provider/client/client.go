package client

import (
	publiccloudservice "terraform-provider-leaseweb/internal/core/services/public_cloud"
	"terraform-provider-leaseweb/internal/facades/public_cloud"
	"terraform-provider-leaseweb/internal/repositories/public_cloud_repository"
)

// The Client handles instantiation of the facades.
type Client struct {
	PublicCloudFacade public_cloud.PublicCloudFacade
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

	return Client{
		PublicCloudFacade: public_cloud.NewPublicCloudFacade(publicCloudService),
	}
}
