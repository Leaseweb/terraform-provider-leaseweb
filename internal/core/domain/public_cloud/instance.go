package public_cloud

import (
	"fmt"
	"slices"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
)

type ErrInvalidInstanceTypePassed struct {
	msg string
}

func (e ErrInvalidInstanceTypePassed) Error() string {
	return e.msg
}

type ReasonInstanceCannotBeTerminated string

type Instance struct {
	Id                  string
	Region              Region
	Reference           *string
	StartedAt           *time.Time
	Resources           Resources
	Image               Image
	State               enum.State
	ProductType         string
	HasPublicIpv4       bool
	HasPrivateNetwork   bool
	HasUserData         bool
	Type                InstanceType
	RootDiskStorageType enum.StorageType
	RootDiskSize        value_object.RootDiskSize
	Ips                 Ips
	Contract            Contract
	Iso                 *Iso
	MarketAppId         *string
	PrivateNetwork      *PrivateNetwork
	SshKey              *value_object.SshKey
	AutoScalingGroup    *AutoScalingGroup
}

// CanBeTerminated determines if an Instance can be terminated.
// This is not the case when:
//   - Instance.State is enum.StateCreating
//   - Instance.State is enum.StateDestroying
//   - Instance.State is enum.StateDestroyed
//   - Contract.EndsAt is not null
func (i Instance) CanBeTerminated() (bool, *ReasonInstanceCannotBeTerminated) {
	if i.State == enum.StateCreating || i.State == enum.StateDestroying || i.State == enum.StateDestroyed {
		reason := ReasonInstanceCannotBeTerminated(
			fmt.Sprintf("state is %q", i.State.String()),
		)
		return false, &reason
	}

	if i.Contract.EndsAt != nil {
		reason := ReasonInstanceCannotBeTerminated(
			fmt.Sprintf("contract.endsAt is %q", i.Contract.EndsAt.String()),
		)
		return false, &reason
	}

	return true, nil
}

// OptionalInstanceValues contains optional supported instance fields.
type OptionalInstanceValues struct {
	Reference        *string
	Iso              *Iso
	MarketAppId      *string
	SshKey           *value_object.SshKey
	StartedAt        *time.Time
	PrivateNetwork   *PrivateNetwork
	AutoScalingGroup *AutoScalingGroup
}

// OptionalCreateInstanceValues contains optional supported fields for instance creation.
type OptionalCreateInstanceValues struct {
	MarketAppId  *string
	Reference    *string
	SshKey       *value_object.SshKey
	RootDiskSize *value_object.RootDiskSize
}

type OptionalUpdateInstanceValues struct {
	Type             *string
	Reference        *string
	ContractType     *enum.ContractType
	Term             *enum.ContractTerm
	BillingFrequency *enum.ContractBillingFrequency
	RootDiskSize     *value_object.RootDiskSize
}

// NewInstance creates a new instance with all supported options.
func NewInstance(
	id string,
	region Region,
	resources Resources,
	image Image,
	state enum.State,
	productType string,
	hasPublicIpv4 bool,
	hasPrivateNetwork bool,
	hasUserData bool,
	rootDiskSize value_object.RootDiskSize,
	instanceType InstanceType,
	rootDiskStorageType enum.StorageType,
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
		HasUserData:         hasUserData,
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

// NewCreateInstance creates a new instance with only all the supported fields for instance creation.
func NewCreateInstance(
	region string,
	instanceType string,
	rootDiskStorageType enum.StorageType,
	imageId string,
	contractType enum.ContractType,
	contractTerm enum.ContractTerm,
	billingFrequency enum.ContractBillingFrequency,
	optional OptionalCreateInstanceValues,
	allowedInstanceTypes []string,
) (*Instance, error) {
	if !slices.Contains(allowedInstanceTypes, instanceType) {
		return nil, ErrInvalidInstanceTypePassed{
			msg: fmt.Sprintf("instance type %q is not allowed", instanceType),
		}
	}

	instance := Instance{
		Region:              Region{Name: region},
		Type:                InstanceType{Name: instanceType},
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

	return &instance, nil
}

// NewUpdateInstance creates a new instance with only all the supported fields for instance updates.
func NewUpdateInstance(
	id string,
	options OptionalUpdateInstanceValues,
	allowedInstanceTypes []string,
	currentInstanceType string,
) (*Instance, error) {
	instance := Instance{Id: id}

	instance.Reference = options.Reference

	allowedInstanceTypes = append(allowedInstanceTypes, currentInstanceType)
	if options.Type != nil {
		if !slices.Contains(allowedInstanceTypes, *options.Type) {
			return nil, ErrInvalidInstanceTypePassed{
				msg: fmt.Sprintf("instance type %q is not allowed", *options.Type),
			}
		}
		instance.Type = InstanceType{Name: *options.Type}
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

	return &instance, nil
}
