// Package client implements access to facades.
package client

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/contracts"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/repository"
)

// ProviderData TODO: Refactor this part, data can be managed directly, not within client.
type ProviderData struct {
	ApiKey string
	Host   *string
	Scheme *string
}

// The Client handles instantiation of the facades.
type Client struct {
	ProviderData          ProviderData
	PublicCloudRepository contracts.PublicCloudRepository
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

	return Client{
		ProviderData: ProviderData{
			ApiKey: token,
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
		PublicCloudRepository: &publicCloudRepository,
	}
}
