package enum

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
	return findEnumForString(s, contractStateValues, ContractStateActive)
}
