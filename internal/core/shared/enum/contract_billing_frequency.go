package enum

type ContractBillingFrequency int

func (c ContractBillingFrequency) Value() int {
	return int(c)
}

func (c ContractBillingFrequency) Values() []int {
	var values []int

	for _, v := range contractBillingFrequencies {
		values = append(values, int(v))
	}

	return values
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
	return findEnumForInt(
		value,
		contractBillingFrequencies,
		ContractBillingFrequencyZero,
	)
}
