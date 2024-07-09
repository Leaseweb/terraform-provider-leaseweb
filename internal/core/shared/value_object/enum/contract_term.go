package enum

type ContractTerm int64

func (t ContractTerm) Value() int64 {
	return int64(t)
}

type ContractTerms []ContractTerm

const (
	ContractTermZero ContractTerm = iota
	ContractTermOne
	ContractTermThree  ContractTerm = iota + 1
	ContractTermSix    ContractTerm = iota + 3
	ContractTermTwelve ContractTerm = iota + 8
)

var ContractTermValues = ContractTerms{
	ContractTermZero,
	ContractTermOne,
	ContractTermThree,
	ContractTermSix,
	ContractTermTwelve,
}
