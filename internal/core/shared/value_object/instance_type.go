package value_object

import (
	"fmt"
	"slices"
)

type ErrInvalidInstanceType struct {
	msg string
}

func (e ErrInvalidInstanceType) Error() string {
	return e.msg
}

// InstanceType validates that the passed instance appears in allowedTypes.
type InstanceType struct {
	Type string
}

func (i InstanceType) String() string {
	return i.Type
}

// NewInstanceType generated an InstanceType that is in allowedInstanceTypes.
func NewInstanceType(
	instanceType string,
	allowedInstanceTypes []string,
) (*InstanceType, error) {
	if !slices.Contains(allowedInstanceTypes, instanceType) {
		return nil, ErrInvalidInstanceType{msg: fmt.Sprintf(
			"instance type %q is not allowed",
			instanceType,
		)}
	}

	return &InstanceType{Type: instanceType}, nil
}

// NewUnvalidatedInstanceType allows a new InstanceType to be generated without validation.
func NewUnvalidatedInstanceType(instanceType string) InstanceType {
	return InstanceType{Type: instanceType}
}
