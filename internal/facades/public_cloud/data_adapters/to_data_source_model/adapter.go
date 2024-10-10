// Package to_data_source_model implements adapters to convert domain entities to sdk models.
package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
)

func AdaptInstances(domainInstances public_cloud.Instances) model.Instances {
	var instances model.Instances

	for _, domainInstance := range domainInstances {
		instance := adaptInstance(domainInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func adaptInstance(domainInstance public_cloud.Instance) model.Instance {
	instance := model.Instance{
		Id:     basetypes.NewStringValue(domainInstance.Id),
		Region: basetypes.NewStringValue(domainInstance.Region),
		Reference: shared.AdaptNullableStringToStringValue(
			domainInstance.Reference,
		),
		Image: adaptImage(domainInstance.Image),
		State: basetypes.NewStringValue(string(domainInstance.State)),
		Type:  basetypes.NewStringValue(domainInstance.Type),
		RootDiskSize: basetypes.NewInt64Value(
			int64(domainInstance.RootDiskSize.Value),
		),
		RootDiskStorageType: basetypes.NewStringValue(
			string(domainInstance.RootDiskStorageType),
		),
		Contract: adaptContract(
			domainInstance.Contract,
		),
		MarketAppId: shared.AdaptNullableStringToStringValue(
			domainInstance.MarketAppId,
		),
	}

	for _, autoScalingGroupIp := range domainInstance.Ips {
		ip := adaptIp(autoScalingGroupIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}

func adaptImage(domainImage public_cloud.Image) model.Image {
	image := model.Image{
		Id: basetypes.NewStringValue(domainImage.Id),
	}

	return image
}

func adaptContract(contract public_cloud.Contract) model.Contract {
	return model.Contract{
		BillingFrequency: basetypes.NewInt64Value(
			int64(contract.BillingFrequency),
		),
		Term:   basetypes.NewInt64Value(int64(contract.Term)),
		Type:   basetypes.NewStringValue(string(contract.Type)),
		EndsAt: shared.AdaptNullableTimeToStringValue(contract.EndsAt),
		State:  basetypes.NewStringValue(string(contract.State)),
	}
}

func adaptIp(ip public_cloud.Ip) model.Ip {
	return model.Ip{
		Ip: basetypes.NewStringValue(ip.Ip),
	}
}
