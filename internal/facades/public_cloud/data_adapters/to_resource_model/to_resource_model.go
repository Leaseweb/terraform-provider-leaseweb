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
	plan.RootDiskSize = basetypes.NewInt64Value(
		int64(instance.RootDiskSize.Value),
	)
	plan.RootDiskStorageType = basetypes.NewStringValue(
		string(instance.RootDiskStorageType),
	)
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
		State:            basetypes.NewStringValue(string(contract.State)),
	}, nil
}

func adaptIp(
	ctx context.Context,
	ip public_cloud.Ip,
) (*model.Ip, error) {
	return &model.Ip{
		Ip: basetypes.NewStringValue(ip.Ip),
	}, nil
}
