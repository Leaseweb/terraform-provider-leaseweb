package instances

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type instancesDataSourceModel struct {
	Instances []instancesModel `tfsdk:"instances"`
}

type instancesModel struct {
	Id                  types.String         `tfsdk:"id"`
	EquipmentId         types.String         `tfsdk:"equipment_id"`
	SalesOrgId          types.String         `tfsdk:"sales_org_id"`
	CustomerId          types.String         `tfsdk:"customer_id"`
	Region              types.String         `tfsdk:"region"`
	Reference           types.String         `tfsdk:"reference"`
	Resources           resourcesModel       `tfsdk:"resources"`
	OperatingSystem     operatingSystemModel `tfsdk:"operating_system"`
	State               types.String         `tfsdk:"state"`
	ProductType         types.String         `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool           `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool           `tfsdk:"has_private_network"`
	Type                types.String         `tfsdk:"type"`
	RootDiskSize        types.Int64          `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String         `tfsdk:"root_disk_storage_type"`
	Ips                 []ipModel            `tfsdk:"ips"`
	StartedAt           types.String         `tfsdk:"started_at"`
	Contract            contractModel        `tfsdk:"contract"`
	Iso                 isoModel             `tfsdk:"iso"`
	MarketAppId         types.String         `tfsdk:"market_app_id"`
	PrivateNetwork      privateNetworkModel  `tfsdk:"private_network"`
}

type resourcesModel struct {
	Cpu                 cpuModel          `tfsdk:"cpu"`
	Memory              memoryModel       `tfsdk:"memory"`
	PublicNetworkSpeed  networkSpeedModel `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed networkSpeedModel `tfsdk:"private_network_speed"`
}

type cpuModel struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

type memoryModel struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}

type networkSpeedModel struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

type operatingSystemModel struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Version      types.String   `tfsdk:"version"`
	Family       types.String   `tfsdk:"family"`
	Flavour      types.String   `tfsdk:"flavour"`
	Architecture types.String   `tfsdk:"architecture"`
	MarketApps   []types.String `tfsdk:"market_apps"`
	StorageTypes []types.String `tfsdk:"storage_types"`
}

type ipModel struct {
	Ip            types.String `tfsdk:"ip"`
	PrefixLength  types.String `tfsdk:"prefix_length"`
	Version       types.Int64  `tfsdk:"version"`
	NullRouted    types.Bool   `tfsdk:"null_routed"`
	MainIp        types.Bool   `tfsdk:"main_ip"`
	NetworkType   types.String `tfsdk:"network_type"`
	ReverseLookup types.String `tfsdk:"reverse_lookup"`
	Ddos          dDosModel    `tfsdk:"ddos"`
}

type dDosModel struct {
	DetectionProfile types.String `tfsdk:"detection_profile"`
	ProtectionType   types.String `tfsdk:"protection_type"`
}
type contractModel struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	RenewalsAt       types.String `tfsdk:"renewals_at"`
	CreatedAt        types.String `tfsdk:"created_at"`
	State            types.String `tfsdk:"state"`
}

type isoModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type privateNetworkModel struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}
