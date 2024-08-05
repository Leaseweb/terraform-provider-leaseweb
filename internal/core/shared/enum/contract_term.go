package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type ContractTerm int64

func (t ContractTerm) Value() int {
	return int(t)
}

func (t ContractTerm) Values() []int {
	return enum_utils.ConvertIntEnumToValues(contractTerms)
}

const (
	ContractTermZero ContractTerm = iota
	ContractTermOne
	ContractTermThree  ContractTerm = iota + 1
	ContractTermSix    ContractTerm = iota + 3
	ContractTermTwelve ContractTerm = iota + 8
)

var contractTerms = []ContractTerm{
	ContractTermZero,
	ContractTermOne,
	ContractTermThree,
	ContractTermSix,
	ContractTermTwelve,
}

func NewContractTerm(value int) (ContractTerm, error) {
	return enum_utils.FindEnumForInt(value, contractTerms, ContractTermZero)
}
