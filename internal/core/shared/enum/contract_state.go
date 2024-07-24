package enum

import (
	"terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type ContractState string

func (s ContractState) String() string {
	return string(s)
}

const (
	ContractStateActive          ContractState = "ACTIVE"
	ContractStateDeleteScheduled ContractState = "DELETE_SCHEDULED"
)

var contractStateValues = []ContractState{
	ContractStateActive,
	ContractStateDeleteScheduled,
}

func NewContractState(s string) (ContractState, error) {
	return enum_utils.FindEnumForString(s, contractStateValues, ContractStateActive)
}
