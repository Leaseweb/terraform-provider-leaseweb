package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               Image        `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 []Ip         `tfsdk:"ips"`
	Contract            Contract     `tfsdk:"contract"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
}
