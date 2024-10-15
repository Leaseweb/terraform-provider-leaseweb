package instances

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	publicCloud "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared_schemas/public_cloud"
)

func (d *instancesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	publicCloudService := publicCloud.Service{}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The instance unique identifier",
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"reference": schema.StringAttribute{
							Computed:    true,
							Description: "The identifying name set to the instance",
						},
						"image": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "Image ID",
								},
							},
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The instance's current state",
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"root_disk_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances",
						},
						"root_disk_storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "The root disk's storage type",
						},
						"ips": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{Computed: true},
								},
							},
						},
						"contract": public_cloud.Contract(false, &publicCloudService),
						"market_app_id": schema.StringAttribute{
							Computed:    true,
							Description: "Market App ID",
						},
					},
				},
			},
		},
	}
}
