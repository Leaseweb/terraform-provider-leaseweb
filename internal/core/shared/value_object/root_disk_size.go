package value_object

import (
	"fmt"
)

type InvalidRootDiskSize struct {
	msg   string
	value int
}

func (e InvalidRootDiskSize) Error() string {
	return e.msg
}

const MinRootDiskSize int = 5
const MaxRootDiskSize int = 1000

// RootDiskSize ensures that rootDiskSize is between 5 & 1000 gigabytes.
type RootDiskSize struct {
	Value           int
	MinRootDiskSize int
	MaxRootDiskSize int
}

func NewRootDiskSize(value int) (*RootDiskSize, error) {
	if value < MinRootDiskSize {
		return nil, InvalidRootDiskSize{
			msg: fmt.Sprintf(
				"value %d is too small, minimum rootDiskSize is %d",
				value,
				MinRootDiskSize,
			),
			value: value,
		}
	}
	if value > MaxRootDiskSize {
		return nil, InvalidRootDiskSize{
			msg: fmt.Sprintf(
				"value %d is too large, maximum rootDiskSize is %d",
				value,
				MaxRootDiskSize,
			),
			value: value,
		}
	}

	return &RootDiskSize{Value: value}, nil
}
