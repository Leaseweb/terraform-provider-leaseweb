package domain

import (
	"fmt"
)

type ErrInstanceTypeNotFound struct {
	msg string
}

func (e ErrInstanceTypeNotFound) Error() string {
	return e.msg
}

type InstanceTypes []InstanceType

func (i InstanceTypes) ContainsName(name string) bool {
	for _, instanceType := range i {
		if name == instanceType.Name {
			return true
		}
	}

	return false
}

func (i InstanceTypes) GetByName(name string) (*InstanceType, error) {
	for _, instanceType := range i {
		if name == instanceType.Name {
			return &instanceType, nil
		}
	}

	return nil, ErrInstanceTypeNotFound{fmt.Sprintf(
		"instance type with name %q not found",
		name,
	)}
}

func (i InstanceTypes) ToArray() []string {
	var values []string
	for _, instanceType := range i {
		values = append(values, instanceType.String())
	}

	return values
}
