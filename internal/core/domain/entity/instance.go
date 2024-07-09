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
	options OptionalInstanceValues,
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

	instance.Iso = options.Iso
	instance.Reference = options.Reference
	instance.MarketAppId = options.MarketAppId
	instance.SshKey = options.SshKey
	instance.StartedAt = options.StartedAt
	instance.PrivateNetwork = options.PrivateNetwork
	instance.AutoScalingGroup = options.AutoScalingGroup

	return instance
}
