package controlpanels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *dedicatedServerControlPanelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"control_panels": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Id of the control panel.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the control panel.",
						},
					},
				},
			},
			"operating_system_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter control panels by operating system id.",
			},
		},
	}
}
