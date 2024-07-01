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
		fmt.Sprintf("Cannot set instanceType \"%v\" ", instanceType),
	}
}

func cannotSetOperatingSystemId(operatingSystemId string) *OptsError {
	return &OptsError{
		fmt.Sprintf("Cannot set operatingSystemId \"%v\" ", operatingSystemId),
	}
}

func cannotSetRootDiskStorageType(rootDiskStorageType string) *OptsError {
	return &OptsError{
		fmt.Sprintf("Cannot set rootDiskStorageType \"%v\" ", rootDiskStorageType),
	}
}
