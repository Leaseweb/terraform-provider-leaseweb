package instances

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &instancesDataSource{}
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{}
}

type instancesDataSource struct {
	client client.Client
}
