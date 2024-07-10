package opts

import "fmt"

type OptsError struct {
	msg string
}

func (e *OptsError) Error() string {
	return e.msg
}

func cannotSetInstanceType(instanceType string) *OptsError {
	return &OptsError{
		fmt.Sprintf("Cannot set instanceType %q", instanceType),
	}
}

func cannotSetOperatingSystemId(operatingSystemId string) *OptsError {
	return &OptsError{
		fmt.Sprintf("Cannot set operatingSystemId %q", operatingSystemId),
	}
}

func cannotSetRootDiskStorageType(rootDiskStorageType string) *OptsError {
	return &OptsError{
		fmt.Sprintf(
			"Cannot set rootDiskStorageType %q",
			rootDiskStorageType,
		),
	}
}

func cannotSetContractTerm(term int64) *OptsError {
	return &OptsError{
		fmt.Sprintf("Cannot set contract.term %d", term),
	}
}

func cannotSetContractBillingFrequency(billingFrequency int64) *OptsError {
	return &OptsError{
		fmt.Sprintf(
			"Cannot set contract.billingFrequency %d",
			billingFrequency,
		),
	}
}

func cannotSetContractType(contractType string) *OptsError {
	return &OptsError{
		fmt.Sprintf(
			"Cannot set contract.type %q",
			contractType,
		),
	}
}
