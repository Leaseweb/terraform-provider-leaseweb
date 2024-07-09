package value_object

import (
	"errors"
)

var ErrRootDiskSizeIsTooSmall = errors.New("rootDiskSize is too small")
var ErrRootDiskSizeIsTooLarge = errors.New("rootDiskSize is too large")

const minRootDiskSize int64 = 5
const maxRootDiskSize int64 = 1000

type RootDiskSize struct {
	Value           int64
	MinRootDiskSize int64
	MaxRootDiskSize int64
}

func NewRootDiskSize(value int64) (*RootDiskSize, error) {
	if value < minRootDiskSize {
		return nil, ErrRootDiskSizeIsTooSmall
	}
	if value > maxRootDiskSize {
		return nil, ErrRootDiskSizeIsTooLarge
	}

	return &RootDiskSize{Value: value}, nil
}
