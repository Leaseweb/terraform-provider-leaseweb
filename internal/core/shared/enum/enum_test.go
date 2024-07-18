package enum

import (
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

func Test_newEnumFromString(t *testing.T) {
	t.Run("returns value when found", func(t *testing.T) {

		enum, err := findEnumForString(
			"tralala",
			[]dummyStringEnum{setStringEnumValue},
			defaultStringEnumValue,
		)

		assert.NoError(t, err)
		assert.Equal(t, setStringEnumValue, enum)
	})

	t.Run("returns error when value is not found", func(t *testing.T) {

		_, err := findEnumForString(
			"abc",
			[]dummyStringEnum{setStringEnumValue},
			defaultStringEnumValue,
		)

		assert.ErrorContains(t, err, "abc")
		assert.ErrorContains(t, err, "enum.dummyStringEnum")
	})
}

func TestFindEnumForInt(t *testing.T) {
	t.Run("returns value when found", func(t *testing.T) {

		enum, err := findEnumForInt(
			6,
			[]dummyIntEnum{setIntEnumValue},
			defaultIntEnumValue,
		)

		assert.NoError(t, err)
		assert.Equal(t, setIntEnumValue, enum)
	})

	t.Run("returns error when value is not found", func(t *testing.T) {

		_, err := findEnumForInt(
			7,
			[]dummyIntEnum{setIntEnumValue},
			defaultIntEnumValue,
		)

		assert.ErrorContains(t, err, "7")
		assert.ErrorContains(t, err, "enum.dummyIntEnum")
	})
}

func Test_convertIntEnumToValues(t *testing.T) {
	got := convertIntEnumToValues([]dummyIntEnum{setIntEnumValue})
	want := []int{6}

	assert.Equal(t, want, got)
}

func Test_convertStringEnumToValues(t *testing.T) {
	got := convertStringEnumToValues([]dummyStringEnum{setStringEnumValue})
	want := []string{"tralala"}

	assert.Equal(t, want, got)
}
