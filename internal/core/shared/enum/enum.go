package enum

import (
	"fmt"
)

type errCannotFindEnumForValue[T string | int] struct {
	msg string
}

func (c errCannotFindEnumForValue[T]) Error() string {
	return c.msg
}

func findEnumForString[T fmt.Stringer](
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

func findEnumForInt[T intEnum](
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

func convertIntEnumToValues[T intEnum](enumValues []T) []int {
	var values []int

	for _, v := range enumValues {
		values = append(values, v.Value())
	}

	return values
}
