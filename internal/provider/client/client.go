package client

import (
	publiccloudservice "terraform-provider-leaseweb/internal/core/services/public_cloud"
	"terraform-provider-leaseweb/internal/handlers/public_cloud"
	"terraform-provider-leaseweb/internal/repositories/public_cloud_repository"
)

type Client struct {
	PublicCloudHandler public_cloud.PublicCloudHandler
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	publicCloudRepository := public_cloud_repository.NewPublicCloudRepository(
		token,
		public_cloud_repository.Optional{Host: optional.Host, Scheme: optional.Scheme},
	)
	publicCloudService := publiccloudservice.New(publicCloudRepository)

	return Client{PublicCloudHandler: public_cloud.NewPublicCloudHandler(publicCloudService)}
}
