package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
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
	MarketAppId         types.String `tfsdk:"market_app_id"`
	AutoScalingGroup    types.Object `tfsdk:"auto_scaling_group"`
	Iso                 types.Object `tfsdk:"iso"`
	PrivateNetwork      types.Object `tfsdk:"private_network"`
	SshKey              types.String `tfsdk:"ssh_key"`
}

func (i *Instance) Populate(
	instance *publicCloud.InstanceDetails,
	autoScalingGroupDetails *publicCloud.AutoScalingGroupDetails,
	loadBalancerDetails *publicCloud.LoadBalancerDetails,
	ctx context.Context,
) diag.Diagnostics {
	i.Id = basetypes.NewStringValue(instance.GetId())
	i.Region = basetypes.NewStringValue(instance.GetRegion())
	i.Reference = basetypes.NewStringValue(instance.GetReference())
	i.State = basetypes.NewStringValue(string(instance.GetState()))
	i.ProductType = basetypes.NewStringValue(instance.GetProductType())
	i.HasPublicIpv4 = basetypes.NewBoolValue(instance.GetHasPublicIpV4())
	i.HasPrivateNetwork = basetypes.NewBoolValue(instance.GetIncludesPrivateNetwork())
	i.Type = basetypes.NewStringValue(string(instance.GetType()))
	i.RootDiskSize = basetypes.NewInt64Value(int64(instance.GetRootDiskSize()))
	i.RootDiskStorageType = basetypes.NewStringValue(string(instance.GetRootDiskStorageType()))
	i.StartedAt = basetypes.NewStringValue(instance.GetStartedAt().String())
	i.MarketAppId = basetypes.NewStringValue(instance.GetMarketAppId())

	operatingSystemObject, diags := utils.ConvertSdkModelToResourceObject(
		instance.GetOperatingSystem(),
		OperatingSystem{}.AttributeTypes(),
		ctx,
		newOperatingSystem,
	)
	if diags.HasError() {
		return diags
	}
	i.OperatingSystem = operatingSystemObject

	contractObject, diags := utils.ConvertSdkModelToResourceObject(
		instance.GetContract(),
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if diags.HasError() {
		return diags
	}
	i.Contract = contractObject

	isoObject, diags := utils.ConvertSdkModelToResourceObject(
		instance.GetIso(),
		Iso{}.AttributeTypes(),
		ctx,
		newIso,
	)
	if diags.HasError() {
		return diags
	}
	i.Iso = isoObject

	privateNetworkObject, diags := utils.ConvertSdkModelToResourceObject(
		instance.GetPrivateNetwork(),
		PrivateNetwork{}.AttributeTypes(),
		ctx,
		newPrivateNetwork,
	)
	if diags.HasError() {
		return diags
	}
	i.PrivateNetwork = privateNetworkObject

	resourcesObject, diags := utils.ConvertSdkModelToResourceObject(
		instance.GetResources(),
		Resources{}.AttributeTypes(),
		ctx,
		newResources,
	)
	if diags.HasError() {
		return diags
	}
	i.Resources = resourcesObject

	autoScalingGroupObject, diags := generateAutoScalingGroup(ctx, autoScalingGroupDetails, loadBalancerDetails)
	if diags.HasError() {
		return diags
	}
	i.AutoScalingGroup = autoScalingGroupObject

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
		types.ObjectType{AttrTypes: Ip{}.AttributeTypes()},
		ips,
	)
	if diags != nil {
		return diags
	}
	i.Ips = ipsObject

	return nil
}

func generateAutoScalingGroup(
	ctx context.Context,
	sdkAutoScalingGroupDetails *publicCloud.AutoScalingGroupDetails,
	sdkLoadBalancerDetails *publicCloud.LoadBalancerDetails,
) (basetypes.ObjectValue, diag.Diagnostics) {
	if sdkAutoScalingGroupDetails == nil {
		return types.ObjectNull(AutoScalingGroup{}.AttributeTypes()), nil
	}

	autoScalingGroup, diags := newAutoScalingGroup(
		ctx,
		*sdkAutoScalingGroupDetails,
		sdkLoadBalancerDetails,
	)
	if diags.HasError() {
		return types.ObjectNull(AutoScalingGroup{}.AttributeTypes()), diags
	}

	autoScalingGroupObject, diags := types.ObjectValueFrom(
		ctx,
		autoScalingGroup.AttributeTypes(),
		autoScalingGroup,
	)
	if diags.HasError() {
		return types.ObjectNull(AutoScalingGroup{}.AttributeTypes()), diags
	}

	return autoScalingGroupObject, nil
}
