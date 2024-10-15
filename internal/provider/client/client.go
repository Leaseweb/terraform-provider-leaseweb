// Package client implements access to facades.
package client

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/contracts"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/repository"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/service/public_cloud"
)

// ProviderData TODO: Refactor this part, data can be managed directly, not within client.
type ProviderData struct {
	ApiKey string
	Host   *string
	Scheme *string
}

// The Client handles instantiation of the facades.
type Client struct {
	ProviderData       ProviderData
	PublicCloudService contracts.PublicCloudService
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	publicCloudRepository := repository.NewPublicCloudRepository(
		token,
		repository.Optional{
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
	)
	publicCloudService := public_cloud.New(publicCloudRepository)

	return Client{
		ProviderData: ProviderData{
			ApiKey: token,
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
		PublicCloudService: &publicCloudService,
	}
}
