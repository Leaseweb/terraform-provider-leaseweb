package enum_utils

import (
	"fmt"
)

type errCannotFindEnumForValue[T string | int] struct {
	msg string
}

func (c errCannotFindEnumForValue[T]) Error() string {
	return c.msg
}

// FindEnumForString Return enum for passed string or return an error if the enum is not found.
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
	return defaultEnum, &errCannotFindEnumForValue[string]{msg: fmt.Sprintf(
		"cannot find %T for %q",
		defaultEnum,
		value,
	)}
}

// FindEnumForInt Return enum for passed into or return an error if the enum is not found.
func FindEnumForInt[T intEnum](
	value int,
	enumValues []T,
	defaultEnum T,
) (T, error) {
	for _, enum := range enumValues {
		if enum.Value() == value {
			return enum, nil
		}
	}
	return defaultEnum, &errCannotFindEnumForValue[int]{msg: fmt.Sprintf(
		"cannot find %T for %d",
		defaultEnum,
		value,
	)}
}

// ConvertIntEnumToValues Return an int array for all enum values.
func ConvertIntEnumToValues[T intEnum](enumValues []T) []int {
	var values []int

	for _, v := range enumValues {
		values = append(values, v.Value())
	}

	return values
}

// ConvertStringEnumToValues Return a string array for all enum values.
func ConvertStringEnumToValues[T fmt.Stringer](enumValues []T) []string {
	var values []string

	for _, v := range enumValues {
		values = append(values, v.String())
	}

	return values
}
