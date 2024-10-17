package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/model"
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

func newInstance(sdkInstance publicCloud.Instance) Instance {
	var ips []Ip
	for _, ip := range sdkInstance.Ips {
		ips = append(ips, newIp(ip))
	}

	return Instance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           model.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		Image:               newImage(sdkInstance.Image),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		Ips:                 ips,
		Contract:            newContract(sdkInstance.Contract),
		MarketAppId:         model.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}
}
