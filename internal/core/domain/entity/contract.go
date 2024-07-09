package entity

import (
	"fmt"
	"time"

	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var ErrContractTermCannotBeZero = fmt.Errorf(
	"contract.term cannot be 0 when contract.type is %q",
	enum.ContractTypeMonthly,
)

var ErrContractTermMustBeZero = fmt.Errorf(
	"contract.term must be 0 when contract.type is %q",
	enum.ContractTypeHourly,
)

type Contract struct {
	BillingFrequency enum.ContractBillingFrequency
	Term             enum.ContractTerm
	Type             enum.ContractType
	EndsAt           *time.Time
	RenewalsAt       time.Time
	CreatedAt        time.Time
	State            enum.ContractState
}

func NewContract(
	billingFrequency enum.ContractBillingFrequency,
	term enum.ContractTerm,
	contractType enum.ContractType,
	renewalsAt time.Time,
	createdAt time.Time,
	state enum.ContractState,
	endsAt *time.Time,
) (*Contract, error) {
	contract := Contract{
		BillingFrequency: billingFrequency,
		Term:             term,
		Type:             contractType,
		RenewalsAt:       renewalsAt,
		CreatedAt:        createdAt,
		State:            state,
		EndsAt:           endsAt,
	}

	if contractType == enum.ContractTypeMonthly && term == enum.ContractTermZero {
		return nil, ErrContractTermCannotBeZero
	}

	if contractType == enum.ContractTypeHourly && term != enum.ContractTermZero {
		return nil, ErrContractTermMustBeZero
	}

	return &contract, nil
}
