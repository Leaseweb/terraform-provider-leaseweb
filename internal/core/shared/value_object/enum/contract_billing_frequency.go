package enum

type ContractBillingFrequency int64

type ContractBillingFrequencies []ContractBillingFrequency

func (c ContractBillingFrequency) Value() int64 {
	return int64(c)
}

const (
	ContractBillingFrequencyZero ContractBillingFrequency = iota
	ContractBillingFrequencyOne
	ContractBillingFrequencyThree  ContractBillingFrequency = iota + 1
	ContractBillingFrequencySix    ContractBillingFrequency = iota + 3
	ContractBillingFrequencyTwelve ContractBillingFrequency = iota + 8
)

var ContractBillingFrequencyValues = ContractBillingFrequencies{
	ContractBillingFrequencyZero,
	ContractBillingFrequencyOne,
	ContractBillingFrequencyThree,
	ContractBillingFrequencySix,
	ContractBillingFrequencyTwelve,
}
