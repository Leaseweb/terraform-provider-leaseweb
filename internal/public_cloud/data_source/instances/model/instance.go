package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type instance struct {
	Id                  types.String      `tfsdk:"id"`
	Region              types.String      `tfsdk:"region"`
	Reference           types.String      `tfsdk:"reference"`
	Resources           resources         `tfsdk:"resources"`
	OperatingSystem     operatingSystem   `tfsdk:"operating_system"`
	State               types.String      `tfsdk:"state"`
	ProductType         types.String      `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool        `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool        `tfsdk:"has_private_network"`
	Type                types.String      `tfsdk:"type"`
	RootDiskSize        types.Int64       `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String      `tfsdk:"root_disk_storage_type"`
	Ips                 []ip              `tfsdk:"ips"`
	StartedAt           types.String      `tfsdk:"started_at"`
	Contract            contract          `tfsdk:"contract"`
	MarketAppId         types.String      `tfsdk:"market_app_id"`
	AutoScalingGroup    *autoScalingGroup `tfsdk:"auto_scaling_group"`
	Iso                 *iso              `tfsdk:"iso"`
	PrivateNetwork      privateNetwork    `tfsdk:"private_network"`
}

func newInstance(sdkInstanceDetails publicCloud.InstanceDetails) instance {
	instanceIso, instanceIsoOk := sdkInstanceDetails.GetIsoOk()
	instanceAutoScalingGroup, instanceAutoScalingGroupOk := sdkInstanceDetails.GetAutoScalingGroupOk()

	instance := instance{
		Id:                  basetypes.NewStringValue(sdkInstanceDetails.GetId()),
		Region:              basetypes.NewStringValue(sdkInstanceDetails.GetRegion()),
		Reference:           basetypes.NewStringValue(sdkInstanceDetails.GetReference()),
		Resources:           newResources(sdkInstanceDetails.GetResources()),
		OperatingSystem:     newOperatingSystem(sdkInstanceDetails.GetOperatingSystem()),
		State:               basetypes.NewStringValue(string(sdkInstanceDetails.GetState())),
		ProductType:         basetypes.NewStringValue(sdkInstanceDetails.GetProductType()),
		HasPublicIpv4:       basetypes.NewBoolValue(sdkInstanceDetails.GetHasPublicIpV4()),
		HasPrivateNetwork:   basetypes.NewBoolValue(sdkInstanceDetails.GetIncludesPrivateNetwork()),
		Type:                basetypes.NewStringValue(string(sdkInstanceDetails.GetType())),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstanceDetails.GetRootDiskSize())),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstanceDetails.GetRootDiskStorageType())),
		StartedAt:           basetypes.NewStringValue(sdkInstanceDetails.GetStartedAt().String()),
		Contract:            newContract(sdkInstanceDetails.GetContract()),
		MarketAppId:         basetypes.NewStringValue(sdkInstanceDetails.GetMarketAppId()),
		AutoScalingGroup: utils.ConvertNullableSdkModelToDatasourceModel(
			instanceAutoScalingGroup,
			instanceAutoScalingGroupOk,
			newAutoScalingGroup,
		),
		Iso: utils.ConvertNullableSdkModelToDatasourceModel(
			instanceIso,
			instanceIsoOk,
			newIso,
		),
		PrivateNetwork: newPrivateNetwork(sdkInstanceDetails.GetPrivateNetwork()),
	}

	for _, sdkIpDetails := range sdkInstanceDetails.Ips {
		ip := newIp(sdkIpDetails)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}
