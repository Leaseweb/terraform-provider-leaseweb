package entity

import (
	"time"

	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type Instance struct {
	Id                  uuid.UUID
	Region              string
	Reference           *string
	StartedAt           *time.Time
	Resources           Resources
	Image               Image
	State               enum.State
	ProductType         string
	HasPublicIpv4       bool
	HasPrivateNetwork   bool
	Type                string
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

type OptionalInstanceValues struct {
	Reference        *string
	Iso              *Iso
	MarketAppId      *string
	SshKey           *value_object.SshKey
	StartedAt        *time.Time
	PrivateNetwork   *PrivateNetwork
	AutoScalingGroup *AutoScalingGroup
}

type OptionalCreateInstanceValues struct {
	MarketAppId *string
	Reference   *string
	SshKey      *value_object.SshKey
}

func NewInstance(
	id uuid.UUID,
	region string,
	resources Resources,
	image Image,
	state enum.State,
	productType string,
	hasPublicIpv4 bool,
	hasPrivateNetwork bool,
	rootDiskSize value_object.RootDiskSize,
	instanceType string,
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

func NewCreateInstance(
	region string,
	instanceType string,
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

	return instance
}
