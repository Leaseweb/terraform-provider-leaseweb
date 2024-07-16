package client

import (
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/services/public_cloud_service"
	"terraform-provider-leaseweb/internal/repositories/instance_repository"
)

type Client struct {
	PublicCloud ports.PublicCloudService
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

	return Client{PublicCloud: public_cloud_service.New(publicCloudRepository)}
}
