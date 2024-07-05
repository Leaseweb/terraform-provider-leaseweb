package enum

type ContractState string

const (
	Active          ContractState = "ACTIVE"
	DeleteScheduled ContractState = "DELETE_SCHEDULED"
)
