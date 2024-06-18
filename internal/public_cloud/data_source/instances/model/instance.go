package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type instance struct {
	Id                  types.String    `tfsdk:"id"`
	EquipmentId         types.String    `tfsdk:"equipment_id"`
	SalesOrgId          types.String    `tfsdk:"sales_org_id"`
	CustomerId          types.String    `tfsdk:"customer_id"`
	Region              types.String    `tfsdk:"region"`
	Reference           types.String    `tfsdk:"reference"`
	Resources           resources       `tfsdk:"resource"`
	OperatingSystem     operatingSystem `tfsdk:"operating_system"`
	State               types.String    `tfsdk:"state"`
	ProductType         types.String    `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool      `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool      `tfsdk:"has_private_network"`
	Type                types.String    `tfsdk:"type"`
	RootDiskSize        types.Int64     `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String    `tfsdk:"root_disk_storage_type"`
	Ips                 []ip            `tfsdk:"ips"`
	StartedAt           types.String    `tfsdk:"started_at"`
	Contract            contract        `tfsdk:"contract"`
	Iso                 iso             `tfsdk:"iso"`
	MarketAppId         types.String    `tfsdk:"market_app_id"`
	PrivateNetwork      privateNetwork  `tfsdk:"private_network"`
}

func newInstance(sdkInstance publicCloud.Instance) instance {
	instance := instance{
		Id:                  basetypes.NewStringValue(sdkInstance.GetId()),
		EquipmentId:         basetypes.NewStringValue(sdkInstance.GetEquipmentId()),
		SalesOrgId:          basetypes.NewStringValue(sdkInstance.GetSalesOrgId()),
		CustomerId:          basetypes.NewStringValue(sdkInstance.GetCustomerId()),
		Region:              basetypes.NewStringValue(sdkInstance.GetRegion()),
		Reference:           basetypes.NewStringValue(sdkInstance.GetReference()),
		Resources:           newResources(*sdkInstance.Resources),
		OperatingSystem:     newOperatingSystem(*sdkInstance.OperatingSystem),
		State:               basetypes.NewStringValue(string(sdkInstance.GetState())),
		ProductType:         basetypes.NewStringValue(sdkInstance.GetProductType()),
		HasPublicIpv4:       basetypes.NewBoolValue(sdkInstance.GetHasPublicIpV4()),
		HasPrivateNetwork:   basetypes.NewBoolValue(sdkInstance.GetincludesPrivateNetwork()),
		Type:                basetypes.NewStringValue(string(sdkInstance.GetType())),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.GetRootDiskSize())),
		RootDiskStorageType: basetypes.NewStringValue(sdkInstance.GetRootDiskStorageType()),
		StartedAt:           basetypes.NewStringValue(sdkInstance.GetStartedAt().String()),
		Contract:            newContract(sdkInstance.GetContract()),
		Iso:                 newIso(sdkInstance.GetIso()),
		MarketAppId:         basetypes.NewStringValue(sdkInstance.GetMarketAppId()),
		PrivateNetwork:      newPrivateNetwork(sdkInstance.GetPrivateNetwork()),
	}

	for _, sdkIp := range sdkInstance.Ips {
		ip := newIp(sdkIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}
