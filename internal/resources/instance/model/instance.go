package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	EquipmentId         types.String `tfsdk:"equipment_id"`
	SalesOrgId          types.String `tfsdk:"sales_org_id"`
	CustomerId          types.String `tfsdk:"customer_id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Resources           types.Object `tfsdk:"resources"`
	OperatingSystem     types.Object `tfsdk:"operating_system"`
	State               types.String `tfsdk:"state"`
	ProductType         types.String `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool   `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool   `tfsdk:"has_private_network"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 types.List   `tfsdk:"ips"`
	StartedAt           types.String `tfsdk:"started_at"`
	Contract            types.Object `tfsdk:"contract"`
	Iso                 types.Object `tfsdk:"iso"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
	PrivateNetwork      types.Object `tfsdk:"private_network"`
}

func (i *Instance) Populate(instance *publicCloud.Instance, ctx context.Context) {
	operatingSystem := newOperatingSystem(instance.OperatingSystem)
	contract := newContract(instance.Contract)
	iso := newIso(instance.GetIso())
	privateNetwork := newPrivateNetwork(instance.GetPrivateNetwork())

	i.Id = resources.GetStringValue(instance.HasId(), instance.GetId())
	i.EquipmentId = resources.GetStringValue(instance.HasEquipmentId(), instance.GetEquipmentId())
	i.SalesOrgId = resources.GetStringValue(instance.HasSalesOrgId(), instance.GetSalesOrgId())
	i.CustomerId = resources.GetStringValue(instance.HasCustomerId(), instance.GetCustomerId())
	i.Region = resources.GetStringValue(instance.HasRegion(), instance.GetRegion())
	i.Reference = resources.GetStringValue(instance.HasReference(), instance.GetReference())
	i.State = resources.GetStringValue(instance.HasState(), string(instance.GetState()))
	i.ProductType = resources.GetStringValue(instance.HasProductType(), instance.GetProductType())
	i.HasPublicIpv4 = resources.GetBoolValue(instance.HasHasPublicIpV4(), instance.GetHasPublicIpV4())
	i.HasPrivateNetwork = resources.GetBoolValue(instance.HasincludesPrivateNetwork(), instance.GetincludesPrivateNetwork())
	i.Type = resources.GetStringValue(instance.HasType(), instance.GetType())
	i.RootDiskSize = resources.GetIntValue(instance.HasRootDiskSize(), instance.GetRootDiskSize())
	i.RootDiskStorageType = resources.GetStringValue(instance.HasRootDiskStorageType(), instance.GetRootDiskStorageType())
	i.StartedAt = resources.GetDateTime(instance.GetStartedAt())
	i.MarketAppId = resources.GetStringValue(instance.HasMarketAppId(), instance.GetMarketAppId())
	i.OperatingSystem, _ = types.ObjectValueFrom(ctx, i.OperatingSystem.AttributeTypes(ctx), operatingSystem)
	i.Contract, _ = types.ObjectValueFrom(ctx, i.Contract.AttributeTypes(ctx), contract)
	i.Iso, _ = types.ObjectValueFrom(ctx, i.Iso.AttributeTypes(ctx), iso)
	i.PrivateNetwork, _ = types.ObjectValueFrom(ctx, i.PrivateNetwork.AttributeTypes(ctx), privateNetwork)

	i.populateResources(ctx, instance.Resources)

	var ips []Ip

	for _, ip := range instance.Ips {
		ips = append(ips, newIp(&ip))
	}

	i.Ips, _ = types.ListValueFrom(ctx, i.Ips.ElementType(ctx), ips)
}

func (i *Instance) populateResources(ctx context.Context, sdkResources *publicCloud.InstanceResources) {
	resourcesModel, _ := types.ObjectValueFrom(ctx, i.Resources.AttributeTypes(ctx), newResources(sdkResources))

	i.Resources = resourcesModel
}
