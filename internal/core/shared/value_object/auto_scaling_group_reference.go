package value_object

import (
	"fmt"
)

var ErrReferenceIsTooLong = fmt.Errorf("reference can only be %d characters long", 255)

type AutoScalingGroupReference struct {
	value string
}

func NewAutoScalingGroupReference(value string) (*AutoScalingGroupReference, error) {
	if len(value) > 255 {
		return nil, ErrReferenceIsTooLong
	}

	return &AutoScalingGroupReference{value: value}, nil
}

func (a AutoScalingGroupReference) String() string {
	return a.value
}
