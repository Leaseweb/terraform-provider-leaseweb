// Package to_resource_model implements adapters to convert domain entities to resource models.
package to_resource_model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func AdaptInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*model.Instance, error) {
	plan := model.Instance{}

	plan.Id = basetypes.NewStringValue(instance.Id)
	plan.Region = basetypes.NewStringValue(instance.Region)
	plan.Type = basetypes.NewStringValue(instance.Type)
	plan.Reference = shared.AdaptNullableStringToStringValue(instance.Reference)
	plan.State = basetypes.NewStringValue(string(instance.State))
	plan.ProductType = basetypes.NewStringValue(instance.ProductType)
	plan.HasPublicIpv4 = basetypes.NewBoolValue(instance.HasPublicIpv4)
	plan.HasPrivateNetwork = basetypes.NewBoolValue(instance.HasPrivateNetwork)
	plan.RootDiskSize = basetypes.NewInt64Value(
		int64(instance.RootDiskSize.Value),
	)
	plan.RootDiskStorageType = basetypes.NewStringValue(
		string(instance.RootDiskStorageType),
	)
	plan.StartedAt = shared.AdaptNullableTimeToStringValue(instance.StartedAt)
	plan.MarketAppId = shared.AdaptNullableStringToStringValue(
		instance.MarketAppId,
	)

	// TODO Enable SSH key support
	/**
	  if instance.SshKey != nil {
	  	plan.SshKey = basetypes.NewStringValue(instance.SshKey.String())
	  }
	*/

	image, err := shared.AdaptDomainEntityToResourceObject(
		instance.Image,
		model.Image{}.AttributeTypes(),
		ctx,
		adaptImage,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Image = image

	contract, err := shared.AdaptDomainEntityToResourceObject(
		instance.Contract,
		model.Contract{}.AttributeTypes(),
		ctx,
		adaptContract,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Contract = contract

	iso, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.Iso,
		model.Iso{}.AttributeTypes(),
		ctx,
		adaptIso,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Iso = iso

	privateNetwork, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.PrivateNetwork,
		model.PrivateNetwork{}.AttributeTypes(),
		ctx,
		adaptPrivateNetwork,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.PrivateNetwork = privateNetwork

	resources, err := shared.AdaptDomainEntityToResourceObject(
		instance.Resources,
		model.Resources{}.AttributeTypes(),
		ctx,
		adaptResources,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Resources = resources

	autoScalingGroup, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.AutoScalingGroup,
		model.AutoScalingGroup{}.AttributeTypes(),
		ctx,
		adaptAutoScalingGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.AutoScalingGroup = autoScalingGroup

	ips, err := shared.AdaptEntitiesToListValue(
		instance.Ips,
		model.Ip{}.AttributeTypes(),
		ctx,
		adaptIp,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Ips = ips

	return &plan, nil
}

func adaptImage(
	ctx context.Context,
	image public_cloud.Image,
) (*model.Image, error) {
	plan := &model.Image{}

	plan.Id = basetypes.NewStringValue(image.Id)
	plan.Name = basetypes.NewStringValue(image.Name)
	plan.Family = basetypes.NewStringValue(image.Family)
	plan.Flavour = basetypes.NewStringValue(image.Flavour)
	plan.Custom = shared.AdaptBoolToBoolValue(image.Custom)

	return plan, nil
}

func adaptContract(
	ctx context.Context,
	contract public_cloud.Contract,
) (*model.Contract, error) {

	return &model.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(contract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(contract.Term)),
		Type:             basetypes.NewStringValue(string(contract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(contract.EndsAt),
		RenewalsAt:       basetypes.NewStringValue(contract.RenewalsAt.String()),
		CreatedAt:        basetypes.NewStringValue(contract.CreatedAt.String()),
		State:            basetypes.NewStringValue(string(contract.State)),
	}, nil
}

func adaptIso(
	ctx context.Context,
	iso public_cloud.Iso,
) (*model.Iso, error) {
	return &model.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}, nil
}

func adaptPrivateNetwork(
	ctx context.Context,
	privateNetwork public_cloud.PrivateNetwork,
) (*model.PrivateNetwork, error) {
	return &model.PrivateNetwork{
		PrivateNetworkId: basetypes.NewStringValue(privateNetwork.Id),
		Status:           basetypes.NewStringValue(privateNetwork.Status),
		Subnet:           basetypes.NewStringValue(privateNetwork.Subnet),
	}, nil
}

func adaptResources(
	ctx context.Context,
	domainResources public_cloud.Resources,
) (*model.Resources, error) {
	var resources model.Resources

	cpu, cpuDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.Cpu,
		model.Cpu{}.AttributeTypes(),
		ctx,
		adaptCpu,
	)
	if cpuDiags != nil {
		return &resources, cpuDiags
	}
	resources.Cpu = cpu

	memory, memoryDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.Memory,
		model.Memory{}.AttributeTypes(),
		ctx,
		adaptMemory,
	)
	if memoryDiags != nil {
		return &resources, memoryDiags
	}
	resources.Memory = memory

	publicNetworkSpeed, publicNetworkSpeedDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.PublicNetworkSpeed,
		model.NetworkSpeed{}.AttributeTypes(),
		ctx,
		adaptNetworkSpeed,
	)
	if publicNetworkSpeedDiags != nil {
		return &resources, publicNetworkSpeedDiags
	}
	resources.PublicNetworkSpeed = publicNetworkSpeed

	privateNetworkSpeed, privateNetworkSpeedDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.PrivateNetworkSpeed,
		model.NetworkSpeed{}.AttributeTypes(),
		ctx,
		adaptNetworkSpeed,
	)
	if privateNetworkSpeedDiags != nil {
		return &resources, privateNetworkSpeedDiags
	}
	resources.PrivateNetworkSpeed = privateNetworkSpeed

	return &resources, nil
}

func adaptCpu(
	ctx context.Context,
	cpu public_cloud.Cpu,
) (*model.Cpu, error) {
	return &model.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}, nil
}

func adaptMemory(
	ctx context.Context,
	memory public_cloud.Memory,
) (*model.Memory, error) {
	return &model.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}, nil
}

func adaptNetworkSpeed(
	ctx context.Context,
	networkSpeed public_cloud.NetworkSpeed,
) (*model.NetworkSpeed, error) {

	return &model.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}, nil
}

func adaptAutoScalingGroup(
	ctx context.Context,
	autoScalingGroup public_cloud.AutoScalingGroup,
) (*model.AutoScalingGroup, error) {
	return &model.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.DesiredAmount,
		),
		Region: basetypes.NewStringValue(autoScalingGroup.Region),
		Reference: basetypes.NewStringValue(
			autoScalingGroup.Reference.String(),
		),
		CreatedAt: basetypes.NewStringValue(
			autoScalingGroup.CreatedAt.String(),
		),
		UpdatedAt: basetypes.NewStringValue(
			autoScalingGroup.UpdatedAt.String(),
		),
		StartsAt: shared.AdaptNullableTimeToStringValue(
			autoScalingGroup.StartsAt,
		),
		EndsAt: shared.AdaptNullableTimeToStringValue(
			autoScalingGroup.EndsAt,
		),
		MinimumAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.MinimumAmount,
		),
		MaximumAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.MaximumAmount,
		),
		CpuThreshold: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.CpuThreshold,
		),
		WarmupTime: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.WarmupTime,
		),
		CooldownTime: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.CooldownTime,
		),
	}, nil
}

func adaptIp(
	ctx context.Context,
	ip public_cloud.Ip,
) (*model.Ip, error) {

	ddos, diags := shared.AdaptNullableDomainEntityToResourceObject(
		ip.Ddos,
		model.Ddos{}.AttributeTypes(),
		ctx,
		adaptDdos,
	)

	if diags != nil {
		return nil, diags
	}

	return &model.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.AdaptNullableStringToStringValue(ip.ReverseLookup),
		Ddos:          ddos,
	}, nil
}

func adaptDdos(ctx context.Context, ddos public_cloud.Ddos) (*model.Ddos, error) {
	return &model.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}, nil
}
