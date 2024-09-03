package dedicated_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (d *dedicatedServerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"asset_id": schema.StringAttribute{
				Computed: true,
			},
			"serial_number": schema.StringAttribute{
				Computed: true,
			},
			"rack": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"capacity": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"location": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"rack": schema.StringAttribute{
						Computed: true,
					},
					"site": schema.StringAttribute{
						Computed: true,
					},
					"suite": schema.StringAttribute{
						Computed: true,
					},
					"unit": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"feature_availability": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"automation": schema.BoolAttribute{
						Computed: true,
					},
					"ipmi_reboot": schema.BoolAttribute{
						Computed: true,
					},
					"power_cycle": schema.BoolAttribute{
						Computed: true,
					},
					"private_network": schema.BoolAttribute{
						Computed: true,
					},
					"remote_management": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"contract": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"customer_id": schema.StringAttribute{
						Computed: true,
					},
					"delivery_status": schema.StringAttribute{
						Computed: true,
					},
					"reference": schema.StringAttribute{
						Computed: true,
					},
					"sales_org_id": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"power_ports": portSchema(),
			"private_networks": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"link_speed": schema.Int32Attribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"subnet": schema.StringAttribute{
							Computed: true,
						},
						"vlan_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"network_interfaces": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"public":            networkInterfaceSchema(),
					"internal":          networkInterfaceSchema(),
					"remote_management": networkInterfaceSchema(),
				},
			},
			"specs": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ram": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"size": schema.Int32Attribute{
								Computed: true,
							},
							"unit": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"cpu": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"quantity": schema.Int32Attribute{
								Computed: true,
							},
							"type": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"hardware_raid_capable": schema.BoolAttribute{
						Computed: true,
					},
					"chassis": schema.StringAttribute{
						Computed: true,
					},
					"hdd": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
								"amount": schema.Int32Attribute{
									Computed: true,
								},
								"size": schema.Float32Attribute{
									Computed: true,
								},
								"type": schema.StringAttribute{
									Computed: true,
								},
								"unit": schema.StringAttribute{
									Computed: true,
								},
								"performance_type": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
					"pci_cards": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func networkInterfaceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"mac": schema.StringAttribute{
				Computed: true,
			},
			"ip": schema.StringAttribute{
				Computed: true,
			},
			"gateway": schema.StringAttribute{
				Computed: true,
			},
			"location_id": schema.StringAttribute{
				Computed: true,
			},
			"null_routed": schema.BoolAttribute{
				Computed: true,
			},
			"ports": portSchema(),
		},
	}
}

func portSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Computed: true,
				},
				"port": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
