package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DedicatedServer struct {
	Id                  types.String        `tfsdk:"id"`
	AssetId             types.String        `tfsdk:"asset_id"`
	SerialNumber        types.String        `tfsdk:"serial_number"`
	Rack                Rack                `tfsdk:"rack"`
	Location            Location            `tfsdk:"location"`
	FeatureAvailability FeatureAvailability `tfsdk:"feature_availability"`
	Contract            Contract            `tfsdk:"contract"`
	PowerPorts          Ports               `tfsdk:"power_ports"`
	PrivateNetworks     PrivateNetworks     `tfsdk:"private_networks"`
	NetworkInterfaces   NetworkInterfaces   `tfsdk:"network_interfaces"`
	Specs               Specs               `tfsdk:"specs"`
}
