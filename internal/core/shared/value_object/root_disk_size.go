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

const minRootDiskSize int = 5
const maxRootDiskSize int = 1000

type RootDiskSize struct {
	Value           int
	MinRootDiskSize int
	MaxRootDiskSize int
}

func NewRootDiskSize(value int) (*RootDiskSize, error) {
	if value < minRootDiskSize {
		return nil, InvalidRootDiskSize{
			msg: fmt.Sprintf(
				"value %d is too small, minimum rootDiskSize is %d",
				value,
				minRootDiskSize,
			),
			value: value,
		}
	}
	if value > maxRootDiskSize {
		return nil, InvalidRootDiskSize{
			msg: fmt.Sprintf(
				"value %d is too large, maximum rootDiskSize is %d",
				value,
				maxRootDiskSize,
			),
			value: value,
		}
	}

	return &RootDiskSize{Value: value}, nil
}
