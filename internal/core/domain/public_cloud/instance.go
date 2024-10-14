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
	Region              string
	Reference           *string
	Image               Image
	State               enum.State
	Type                string
	RootDiskStorageType enum.StorageType
	RootDiskSize        value_object.RootDiskSize
	Ips                 Ips
	Contract            Contract
	MarketAppId         *string
	SshKey              *value_object.SshKey
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
	Reference   *string
	MarketAppId *string
	SshKey      *value_object.SshKey
	StartedAt   *time.Time
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
	region string,
	image Image,
	state enum.State,
	rootDiskSize value_object.RootDiskSize,
	instanceType string,
	rootDiskStorageType enum.StorageType,
	ips Ips,
	contract Contract,
	optional OptionalInstanceValues,
) Instance {
	instance := Instance{
		Id:                  id,
		Region:              region,
		Image:               image,
		State:               state,
		Type:                instanceType,
		RootDiskStorageType: rootDiskStorageType,
		RootDiskSize:        rootDiskSize,
		Ips:                 ips,
		Contract:            contract,
	}

	instance.Reference = optional.Reference
	instance.MarketAppId = optional.MarketAppId
	instance.SshKey = optional.SshKey

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

	return &instance, nil
}
