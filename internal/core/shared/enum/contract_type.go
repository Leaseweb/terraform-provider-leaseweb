package enum

type ContractType string

func (t ContractType) String() string {
	return string(t)
}

func (t ContractType) Values() []string {
	return convertStringEnumToValues(contractTypes)
}

const (
	ContractTypeHourly  ContractType = "HOURLY"
	ContractTypeMonthly ContractType = "MONTHLY"
)

var contractTypes = []ContractType{
	ContractTypeHourly, ContractTypeMonthly,
}

func NewContractType(s string) (ContractType, error) {
	return findEnumForString(s, contractTypes, ContractTypeHourly)
}
