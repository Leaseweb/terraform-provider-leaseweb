package enum

type ContractBillingFrequency int64

const (
	ContractBillingFrequencyZero ContractBillingFrequency = iota
	ContractBillingFrequencyOne
	ContractBillingFrequencyThree  ContractBillingFrequency = iota + 1
	ContractBillingFrequencySix    ContractBillingFrequency = iota + 3
	ContractBillingFrequencyTwelve ContractBillingFrequency = iota + 8
)
