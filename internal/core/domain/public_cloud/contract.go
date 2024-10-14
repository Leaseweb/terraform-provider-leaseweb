package public_cloud

import (
	"fmt"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
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
	State            enum.ContractState
}

// NewContract creates a new contract. Also ensures that contractType / contractTerm combination is valid.
func NewContract(
	billingFrequency enum.ContractBillingFrequency,
	term enum.ContractTerm,
	contractType enum.ContractType,
	state enum.ContractState,
	endsAt *time.Time,
) (*Contract, error) {
	contract := Contract{
		BillingFrequency: billingFrequency,
		Term:             term,
		Type:             contractType,
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
