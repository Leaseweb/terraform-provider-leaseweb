package control_panels

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &controlPanelDataSource{}
	_ datasource.DataSourceWithConfigure = &controlPanelDataSource{}
)

func New() datasource.DataSource {
	return &controlPanelDataSource{}
}

type controlPanelDataSource struct {
	client client.Client
}
