package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
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
		Id: utils.GenerateString(
			sdkInstance.HasId(),
			sdkInstance.GetId(),
		),
		EquipmentId: utils.GenerateString(
			sdkInstance.HasEquipmentId(),
			sdkInstance.GetEquipmentId(),
		),
		SalesOrgId: utils.GenerateString(
			sdkInstance.HasSalesOrgId(),
			sdkInstance.GetSalesOrgId(),
		),
		CustomerId: utils.GenerateString(
			sdkInstance.HasCustomerId(),
			sdkInstance.GetCustomerId(),
		),
		Region: utils.GenerateString(
			sdkInstance.HasRegion(),
			sdkInstance.GetRegion(),
		),
		Reference: utils.GenerateString(
			sdkInstance.HasReference(),
			sdkInstance.GetReference(),
		),
		Resources:       newResources(*sdkInstance.Resources),
		OperatingSystem: newOperatingSystem(*sdkInstance.OperatingSystem),
		State: utils.GenerateString(
			sdkInstance.HasState(),
			string(sdkInstance.GetState()),
		),
		ProductType: utils.GenerateString(
			sdkInstance.HasProductType(),
			sdkInstance.GetProductType(),
		),
		HasPublicIpv4: utils.GenerateBool(
			sdkInstance.HasHasPublicIpV4(),
			sdkInstance.GetHasPublicIpV4(),
		),
		HasPrivateNetwork: utils.GenerateBool(
			sdkInstance.HasincludesPrivateNetwork(),
			sdkInstance.GetincludesPrivateNetwork(),
		),
		Type: utils.GenerateString(
			sdkInstance.HasType(),
			sdkInstance.GetType(),
		),
		RootDiskSize: utils.GenerateInt(
			sdkInstance.HasRootDiskSize(),
			sdkInstance.GetRootDiskSize(),
		),
		RootDiskStorageType: utils.GenerateString(
			sdkInstance.HasRootDiskStorageType(),
			sdkInstance.GetRootDiskStorageType(),
		),
		StartedAt: utils.GenerateDateTime(sdkInstance.GetStartedAt()),
		Contract:  newContract(sdkInstance.GetContract()),
		Iso:       newIso(sdkInstance.GetIso()),
		MarketAppId: utils.GenerateString(
			sdkInstance.HasMarketAppId(),
			sdkInstance.GetMarketAppId(),
		),
		PrivateNetwork: newPrivateNetwork(sdkInstance.GetPrivateNetwork()),
	}

	for _, sdkIp := range sdkInstance.Ips {
		ip := newIp(sdkIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}
