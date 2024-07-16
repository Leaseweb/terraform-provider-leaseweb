package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/utils"
)

type instance struct {
	Id                  types.String      `tfsdk:"id"`
	Region              types.String      `tfsdk:"region"`
	Reference           types.String      `tfsdk:"reference"`
	Resources           resources         `tfsdk:"resources"`
	Image               image             `tfsdk:"image"`
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
	PrivateNetwork      *privateNetwork   `tfsdk:"private_network"`
}

func newInstance(entityInstance domain.Instance) instance {

	instance := instance{
		Id:                  basetypes.NewStringValue(entityInstance.Id.String()),
		Region:              basetypes.NewStringValue(entityInstance.Region),
		Reference:           utils.ConvertNullableStringToStringValue(entityInstance.Reference),
		Resources:           newResources(entityInstance.Resources),
		Image:               newImage(entityInstance.Image),
		State:               basetypes.NewStringValue(string(entityInstance.State)),
		ProductType:         basetypes.NewStringValue(entityInstance.ProductType),
		HasPublicIpv4:       basetypes.NewBoolValue(entityInstance.HasPublicIpv4),
		HasPrivateNetwork:   basetypes.NewBoolValue(entityInstance.HasPrivateNetwork),
		Type:                basetypes.NewStringValue(entityInstance.Type),
		RootDiskSize:        basetypes.NewInt64Value(int64(entityInstance.RootDiskSize.Value)),
		RootDiskStorageType: basetypes.NewStringValue(string(entityInstance.RootDiskStorageType)),
		StartedAt:           utils.ConvertNullableTimeToStringValue(entityInstance.StartedAt),
		Contract:            newContract(entityInstance.Contract),
		MarketAppId:         utils.ConvertNullableStringToStringValue(entityInstance.MarketAppId),
		AutoScalingGroup: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityInstance.AutoScalingGroup,
			newAutoScalingGroup,
		),
		Iso: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityInstance.Iso,
			newIso,
		),
		PrivateNetwork: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityInstance.PrivateNetwork,
			newPrivateNetwork,
		),
	}

	for _, entityIp := range entityInstance.Ips {
		ip := newIp(entityIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}
