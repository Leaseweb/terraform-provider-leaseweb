package operating_systems

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &operatingSystemsDataSource{}
	_ datasource.DataSourceWithConfigure = &operatingSystemsDataSource{}
)

func New() datasource.DataSource {
	return &operatingSystemsDataSource{}
}

type operatingSystemsDataSource struct {
	client client.Client
}
