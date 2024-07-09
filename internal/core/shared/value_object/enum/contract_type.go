package enum

type ContractType string

func (t ContractType) String() string {
	return string(t)
}

type ContractTypes []ContractType

const (
	ContractTypeHourly  ContractType = "HOURLY"
	ContractTypeMonthly ContractType = "MONTHLY"
)

var ContractTypeValues = ContractTypes{
	ContractTypeHourly, ContractTypeMonthly,
}
