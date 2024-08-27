package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Instance struct {
	Id                  types.String      `tfsdk:"id"`
	Region              Region            `tfsdk:"region"`
	Reference           types.String      `tfsdk:"reference"`
	Resources           Resources         `tfsdk:"resources"`
	Image               Image             `tfsdk:"image"`
	State               types.String      `tfsdk:"state"`
	ProductType         types.String      `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool        `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool        `tfsdk:"has_private_network"`
	Type                InstanceType      `tfsdk:"type"`
	RootDiskSize        types.Int64       `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String      `tfsdk:"root_disk_storage_type"`
	Ips                 []Ip              `tfsdk:"ips"`
	StartedAt           types.String      `tfsdk:"started_at"`
	Contract            Contract          `tfsdk:"contract"`
	MarketAppId         types.String      `tfsdk:"market_app_id"`
	AutoScalingGroup    *AutoScalingGroup `tfsdk:"auto_scaling_group"`
	Iso                 *Iso              `tfsdk:"iso"`
	PrivateNetwork      *PrivateNetwork   `tfsdk:"private_network"`
	Volume              *Volume           `tfsdk:"volume"`
}
