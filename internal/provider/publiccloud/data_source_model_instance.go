package publiccloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/model"
)

type DataSourceModelInstance struct {
	Id                  types.String            `tfsdk:"id"`
	Region              types.String            `tfsdk:"region"`
	Reference           types.String            `tfsdk:"reference"`
	Image               DataSourceModelImage    `tfsdk:"image"`
	State               types.String            `tfsdk:"state"`
	Type                types.String            `tfsdk:"type"`
	RootDiskSize        types.Int64             `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String            `tfsdk:"root_disk_storage_type"`
	Ips                 []DataSourceModelIp     `tfsdk:"ips"`
	Contract            DataSourceModelContract `tfsdk:"contract"`
	MarketAppId         types.String            `tfsdk:"market_app_id"`
}

func newDataSourceModelInstance(sdkInstance publicCloud.Instance) DataSourceModelInstance {
	var ips []DataSourceModelIp
	for _, ip := range sdkInstance.Ips {
		ips = append(ips, newDataSourceModelIp(ip))
	}

	return DataSourceModelInstance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           model.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		Image:               newDataSourceModelImage(sdkInstance.Image),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		Ips:                 ips,
		Contract:            newDataSourceModelContract(sdkInstance.Contract),
		MarketAppId:         model.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}
}
