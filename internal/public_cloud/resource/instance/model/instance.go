package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	EquipmentId         types.String `tfsdk:"equipment_id"`
	SalesOrgId          types.String `tfsdk:"sales_org_id"`
	CustomerId          types.String `tfsdk:"customer_id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Resources           types.Object `tfsdk:"resource"`
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
	SshKey              types.String `tfsdk:"ssh_key"`
}

func (i *Instance) Populate(instance *publicCloud.Instance, ctx context.Context) diag.Diagnostics {
	operatingSystem := newOperatingSystem(instance.OperatingSystem)
	contract := newContract(instance.Contract)
	iso := newIso(instance.GetIso())
	privateNetwork := newPrivateNetwork(instance.GetPrivateNetwork())

	resourcesModel, diags := newResources(ctx, instance.Resources)
	if diags != nil {
		return diags
	}

	i.Id = basetypes.NewStringValue(instance.GetId())
	i.EquipmentId = basetypes.NewStringValue(instance.GetEquipmentId())
	i.SalesOrgId = basetypes.NewStringValue(instance.GetSalesOrgId())
	i.CustomerId = basetypes.NewStringValue(instance.GetCustomerId())
	i.Region = basetypes.NewStringValue(instance.GetRegion())
	i.Reference = basetypes.NewStringValue(instance.GetReference())
	i.State = basetypes.NewStringValue(string(instance.GetState()))
	i.ProductType = basetypes.NewStringValue(instance.GetProductType())
	i.HasPublicIpv4 = basetypes.NewBoolValue(instance.GetHasPublicIpV4())
	i.HasPrivateNetwork = basetypes.NewBoolValue(instance.GetincludesPrivateNetwork())
	i.Type = basetypes.NewStringValue(string(instance.GetType()))
	i.RootDiskSize = basetypes.NewInt64Value(int64(instance.GetRootDiskSize()))
	i.RootDiskStorageType = basetypes.NewStringValue(instance.GetRootDiskStorageType())
	i.StartedAt = basetypes.NewStringValue(instance.GetStartedAt().String())
	i.MarketAppId = basetypes.NewStringValue(instance.GetMarketAppId())

	operatingSystemObject, diags := types.ObjectValueFrom(
		ctx,
		operatingSystem.attributeTypes(),
		operatingSystem,
	)
	if diags != nil {
		return diags
	}
	i.OperatingSystem = operatingSystemObject

	contractObject, diags := types.ObjectValueFrom(
		ctx,
		contract.attributeTypes(),
		contract,
	)
	if diags != nil {
		return diags
	}
	i.Contract = contractObject

	isoObject, diags := types.ObjectValueFrom(ctx, iso.attributeTypes(), iso)
	if diags != nil {
		return diags
	}
	i.Iso = isoObject

	privateNetworkObject, diags := types.ObjectValueFrom(
		ctx,
		privateNetwork.attributeTypes(),
		privateNetwork,
	)
	if diags != nil {
		return diags
	}
	i.PrivateNetwork = privateNetworkObject

	resourcesObject, diags := types.ObjectValueFrom(
		ctx,
		resourcesModel.attributeTypes(),
		resourcesModel,
	)
	if diags != nil {
		return diags
	}
	i.Resources = resourcesObject

	var ips []Ip
	for _, ip := range instance.Ips {
		ipObject, diags := newIp(ctx, &ip)
		if diags != nil {
			return diags
		}
		ips = append(ips, ipObject)
	}
	ipsObject, diags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: Ip{}.attributeTypes()},
		ips,
	)
	if diags != nil {
		return diags
	}
	i.Ips = ipsObject

	return nil
}
