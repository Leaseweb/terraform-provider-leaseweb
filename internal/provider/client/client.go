// Package client implements access to facades.
package client

import (
	publiccloudservice "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/public_cloud_repository"
)

// TODO: Refactor this part, ProviderData can be managed directly, not within client.
type ProviderData struct {
	ApiKey string
	Host   *string
	Scheme *string
}

// The Client handles instantiation of the facades.
type Client struct {
	ProviderData      ProviderData
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
		ProviderData: ProviderData{
			ApiKey: token,
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
		PublicCloudFacade: public_cloud.NewPublicCloudFacade(&publicCloudService),
	}
}
