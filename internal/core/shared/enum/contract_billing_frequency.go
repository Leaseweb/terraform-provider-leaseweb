package enum

import (
	"terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type ContractBillingFrequency int

func (c ContractBillingFrequency) Value() int {
	return int(c)
}

func (c ContractBillingFrequency) Values() []int {
	return enum_utils.ConvertIntEnumToValues(contractBillingFrequencies)
}

const (
	ContractBillingFrequencyZero ContractBillingFrequency = iota
	ContractBillingFrequencyOne
	ContractBillingFrequencyThree  ContractBillingFrequency = iota + 1
	ContractBillingFrequencySix    ContractBillingFrequency = iota + 3
	ContractBillingFrequencyTwelve ContractBillingFrequency = iota + 8
)

var contractBillingFrequencies = []ContractBillingFrequency{
	ContractBillingFrequencyZero,
	ContractBillingFrequencyOne,
	ContractBillingFrequencyThree,
	ContractBillingFrequencySix,
	ContractBillingFrequencyTwelve,
}

func NewContractBillingFrequency(value int) (ContractBillingFrequency, error) {
	return enum_utils.FindEnumForInt(
		value,
		contractBillingFrequencies,
		ContractBillingFrequencyZero,
	)
}
