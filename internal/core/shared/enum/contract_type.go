package enum

import (
	"terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type ContractType string

func (t ContractType) String() string {
	return string(t)
}

func (t ContractType) Values() []string {
	return enum_utils.ConvertStringEnumToValues(contractTypes)
}

const (
	ContractTypeHourly  ContractType = "HOURLY"
	ContractTypeMonthly ContractType = "MONTHLY"
)

var contractTypes = []ContractType{
	ContractTypeHourly, ContractTypeMonthly,
}

func NewContractType(s string) (ContractType, error) {
	return enum_utils.FindEnumForString(s, contractTypes, ContractTypeHourly)
}
