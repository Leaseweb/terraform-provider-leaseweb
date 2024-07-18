package enum

type ContractTerm int64

func (t ContractTerm) Value() int {
	return int(t)
}

func (t ContractTerm) Values() []int {
	return convertIntEnumToValues(contractTerms)
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
	return findEnumForInt(value, contractTerms, ContractTermZero)
}
