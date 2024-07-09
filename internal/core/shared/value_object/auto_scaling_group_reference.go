package value_object

import (
	"fmt"
)

const maxAutoScalingGroupReferenceLength = 255

var ErrReferenceIsTooLong = fmt.Errorf(
	"reference can only be %d characters long",
	maxAutoScalingGroupReferenceLength,
)

type AutoScalingGroupReference struct {
	value                              string
	MaxAutoScalingGroupReferenceLength int
}

func NewAutoScalingGroupReference(value string) (*AutoScalingGroupReference, error) {
	if len(value) > maxAutoScalingGroupReferenceLength {
		return nil, ErrReferenceIsTooLong
	}

	return &AutoScalingGroupReference{
		value:                              value,
		MaxAutoScalingGroupReferenceLength: maxAutoScalingGroupReferenceLength,
	}, nil
}

func (a AutoScalingGroupReference) String() string {
	return a.value
}
