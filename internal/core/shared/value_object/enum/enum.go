package enum

import (
	"errors"
	"fmt"
)

var ErrCannotFindEnumForValue = errors.New("cannot find enum for value")

func FindEnumForString[T fmt.Stringer](
	value string,
	enumValues []T,
	defaultEnum T,
) (T, error) {
	for _, enum := range enumValues {
		if enum.String() == value {
			return enum, nil
		}
	}
	return defaultEnum, ErrCannotFindEnumForValue
}

func FindEnumForInt[T intEnum](
	value int64,
	enumValues []T,
	defaultEnum T,
) (T, error) {
	for _, enum := range enumValues {
		if enum.Value() == value {
			return enum, nil
		}
	}
	return defaultEnum, ErrCannotFindEnumForValue
}
