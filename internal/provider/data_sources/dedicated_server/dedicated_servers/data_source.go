package dedicated_servers

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &dedicatedServerDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServerDataSource{}
)

func NewDedicatedServerDataSource() datasource.DataSource {
	return &dedicatedServerDataSource{}
}

type dedicatedServerDataSource struct {
	client client.Client
}
