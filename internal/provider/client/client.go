package client

import (
	public_cloud_service "terraform-provider-leaseweb/internal/core/services/public_cloud"
	"terraform-provider-leaseweb/internal/handlers/public_cloud"
	"terraform-provider-leaseweb/internal/repositories/instance_repository"
)

type Client struct {
	PublicCloudHandler public_cloud.PublicCloudHandler
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	publicCloudRepository := instance_repository.NewPublicCloudRepository(
		token,
		instance_repository.Optional{Host: optional.Host, Scheme: optional.Scheme},
	)
	publicCloudService := public_cloud_service.New(publicCloudRepository)

	return Client{PublicCloudHandler: public_cloud.NewPublicCloudHandler(publicCloudService)}
}
