package domain

import (
	"time"

	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type Instance struct {
	Id                  value_object.Uuid
	Region              string
	Reference           *string
	StartedAt           *time.Time
	Resources           Resources
	Image               Image
	State               enum.State
	ProductType         string
	HasPublicIpv4       bool
	HasPrivateNetwork   bool
	Type                enum.InstanceType
	RootDiskStorageType enum.RootDiskStorageType
	RootDiskSize        value_object.RootDiskSize
	Ips                 Ips
	Contract            Contract
	Iso                 *Iso
	MarketAppId         *string
	PrivateNetwork      *PrivateNetwork
	SshKey              *value_object.SshKey
	AutoScalingGroup    *AutoScalingGroup
}

// OptionalInstanceValues Optional supported instance fields.
type OptionalInstanceValues struct {
	Reference        *string
	Iso              *Iso
	MarketAppId      *string
	SshKey           *value_object.SshKey
	StartedAt        *time.Time
	PrivateNetwork   *PrivateNetwork
	AutoScalingGroup *AutoScalingGroup
}

// OptionalCreateInstanceValues Optional supported fields for instance creation.
type OptionalCreateInstanceValues struct {
	MarketAppId  *string
	Reference    *string
	SshKey       *value_object.SshKey
	RootDiskSize *value_object.RootDiskSize
}

type OptionalUpdateInstanceValues struct {
	Type             *enum.InstanceType
	Reference        *string
	ContractType     *enum.ContractType
	Term             *enum.ContractTerm
	BillingFrequency *enum.ContractBillingFrequency
	RootDiskSize     *value_object.RootDiskSize
}

// NewInstance Create a new instance with all supported options.
func NewInstance(
	id value_object.Uuid,
	region string,
	resources Resources,
	image Image,
	state enum.State,
	productType string,
	hasPublicIpv4 bool,
	hasPrivateNetwork bool,
	rootDiskSize value_object.RootDiskSize,
	instanceType enum.InstanceType,
	rootDiskStorageType enum.RootDiskStorageType,
	ips Ips,
	contract Contract,
	optional OptionalInstanceValues,
) Instance {
	instance := Instance{
		Id:                  id,
		Region:              region,
		Resources:           resources,
		Image:               image,
		State:               state,
		ProductType:         productType,
		HasPublicIpv4:       hasPublicIpv4,
		HasPrivateNetwork:   hasPrivateNetwork,
		Type:                instanceType,
		RootDiskStorageType: rootDiskStorageType,
		RootDiskSize:        rootDiskSize,
		Ips:                 ips,
		Contract:            contract,
	}

	instance.Iso = optional.Iso
	instance.Reference = optional.Reference
	instance.MarketAppId = optional.MarketAppId
	instance.SshKey = optional.SshKey
	instance.StartedAt = optional.StartedAt
	instance.PrivateNetwork = optional.PrivateNetwork
	instance.AutoScalingGroup = optional.AutoScalingGroup

	return instance
}

// NewCreateInstance All the supported fields for instance creation.
func NewCreateInstance(
	region string,
	instanceType enum.InstanceType,
	rootDiskStorageType enum.RootDiskStorageType,
	imageId enum.ImageId,
	contractType enum.ContractType,
	contractTerm enum.ContractTerm,
	billingFrequency enum.ContractBillingFrequency,
	optional OptionalCreateInstanceValues,
) Instance {
	instance := Instance{
		Region:              region,
		Type:                instanceType,
		RootDiskStorageType: rootDiskStorageType,
		Image:               Image{Id: imageId},
		Contract: Contract{
			Type:             contractType,
			Term:             contractTerm,
			BillingFrequency: billingFrequency,
		},
	}

	instance.MarketAppId = optional.MarketAppId
	instance.Reference = optional.Reference
	instance.SshKey = optional.SshKey

	if optional.RootDiskSize != nil {
		instance.RootDiskSize = *optional.RootDiskSize
	}

	return instance
}

// NewUpdateInstance All the supported fields for instance updates.
func NewUpdateInstance(
	id value_object.Uuid,
	options OptionalUpdateInstanceValues,
) Instance {
	instance := Instance{Id: id}

	instance.Reference = options.Reference

	if options.Type != nil {
		instance.Type = *options.Type
	}

	if options.ContractType != nil {
		instance.Contract.Type = *options.ContractType
	}
	if options.Term != nil {
		instance.Contract.Term = *options.Term
	}
	if options.BillingFrequency != nil {
		instance.Contract.BillingFrequency = *options.BillingFrequency
	}
	if options.RootDiskSize != nil {
		instance.RootDiskSize = *options.RootDiskSize
	}

	return instance
}
