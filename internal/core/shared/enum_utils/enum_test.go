package enum_utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyStringEnum string

func (d dummyStringEnum) String() string {
	return string(d)
}

type dummyIntEnum int

func (i dummyIntEnum) Value() int {
	return int(i)
}

const (
	defaultStringEnumValue dummyStringEnum = "default"
	setStringEnumValue     dummyStringEnum = "tralala"

	defaultIntEnumValue dummyIntEnum = 5
	setIntEnumValue     dummyIntEnum = 6
)

type StringEnum string

func (e StringEnum) String() string {
	return string(e)
}

type IntEnum int

func (i IntEnum) Value() int {
	return int(i)
}

func TestNewEnumFromString(t *testing.T) {
	t.Run("returns value when found", func(t *testing.T) {

		foundEnum, err := FindEnumForString(
			"tralala",
			[]dummyStringEnum{setStringEnumValue},
			defaultStringEnumValue,
		)

		assert.NoError(t, err)
		assert.Equal(t, setStringEnumValue, foundEnum)
	})

	t.Run("returns error when value is not found", func(t *testing.T) {

		_, err := FindEnumForString(
			"abc",
			[]dummyStringEnum{setStringEnumValue},
			defaultStringEnumValue,
		)

		assert.ErrorContains(t, err, "abc")
		assert.ErrorContains(t, err, "dummyStringEnum")
	})
}

func TestFindEnumForInt(t *testing.T) {
	t.Run("returns value when found", func(t *testing.T) {

		foundEnum, err := FindEnumForInt(
			6,
			[]dummyIntEnum{setIntEnumValue},
			defaultIntEnumValue,
		)

		assert.NoError(t, err)
		assert.Equal(t, setIntEnumValue, foundEnum)
	})

	t.Run("returns error when value is not found", func(t *testing.T) {

		_, err := FindEnumForInt(
			7,
			[]dummyIntEnum{setIntEnumValue},
			defaultIntEnumValue,
		)

		assert.ErrorContains(t, err, "7")
		assert.ErrorContains(t, err, "dummyIntEnum")
	})
}

func TestConvertIntEnumToValues(t *testing.T) {
	got := ConvertIntEnumToValues([]dummyIntEnum{setIntEnumValue})
	want := []int{6}

	assert.Equal(t, want, got)
}

func TestConvertStringEnumToValues(t *testing.T) {
	got := ConvertStringEnumToValues([]dummyStringEnum{setStringEnumValue})
	want := []string{"tralala"}

	assert.Equal(t, want, got)
}

func ExampleFindEnumForString() {
	const (
		FirstValue  StringEnum = "firstValue"
		SecondValue StringEnum = "secondValue"
	)
	enumValues := []StringEnum{FirstValue, SecondValue}

	foundEnum, _ := FindEnumForString("secondValue", enumValues, FirstValue)

	fmt.Println(foundEnum)
	// Output: secondValue
}

func ExampleFindEnumForInt() {
	const (
		FirstValue  IntEnum = 1
		SecondValue IntEnum = 2
	)
	enumValues := []IntEnum{FirstValue, SecondValue}

	foundEnum, _ := FindEnumForInt(2, enumValues, FirstValue)

	fmt.Println(foundEnum)
	// Output: 2
}

func ExampleConvertIntEnumToValues() {
	const (
		FirstValue  IntEnum = 1
		SecondValue IntEnum = 2
	)
	enumValues := []IntEnum{FirstValue, SecondValue}

	values := ConvertIntEnumToValues(enumValues)

	fmt.Println(values)
	// Output: [1 2]
}

func ExampleConvertStringEnumToValues() {
	const (
		FirstValue  StringEnum = "firstValue"
		SecondValue StringEnum = "secondValue"
	)
	enumValues := []StringEnum{FirstValue, SecondValue}

	values := ConvertStringEnumToValues(enumValues)

	fmt.Println(values)
	// Output: [firstValue secondValue]
}
