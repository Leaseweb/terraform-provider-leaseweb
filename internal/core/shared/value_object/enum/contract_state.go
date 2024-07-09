package enum

type ContractState string

func (s ContractState) String() string {
	return string(s)
}

type ContractStates []ContractState

const (
	ContractStateActive          ContractState = "ACTIVE"
	ContractStateDeleteScheduled ContractState = "DELETE_SCHEDULED"
)

var ContractStateValues = ContractStates{
	ContractStateActive,
	ContractStateDeleteScheduled,
}
