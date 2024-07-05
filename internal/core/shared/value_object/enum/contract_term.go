package enum

type ContractTerm int64

const (
	ContractTermZero ContractTerm = iota
	ContractTermOne
	ContractTermThree  ContractTerm = iota + 1
	ContractTermSix    ContractTerm = iota + 3
	ContractTermTwelve ContractTerm = iota + 8
)
