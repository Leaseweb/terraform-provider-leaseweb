package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyStringEnum string

func (d dummyStringEnum) String() string {
	return string(d)
}

type dummyIntEnum int64

func (i dummyIntEnum) Value() int64 {
	return int64(i)
}

const (
	defaultStringEnumValue dummyStringEnum = "default"
	setStringEnumValue     dummyStringEnum = "tralala"

	defaultIntEnumValue dummyIntEnum = 5
	setIntEnumValue     dummyIntEnum = 6
)

func Test_newEnumFromString(t *testing.T) {
	t.Run("returns value when found", func(t *testing.T) {

		enum, err := FindEnumForString(
			"tralala",
			[]dummyStringEnum{setStringEnumValue},
			defaultStringEnumValue,
		)

		assert.NoError(t, err)
		assert.Equal(t, setStringEnumValue, enum)
	})

	t.Run("returns error when value is not found", func(t *testing.T) {

		_, err := FindEnumForString(
			"abc",
			[]dummyStringEnum{setStringEnumValue},
			defaultStringEnumValue,
		)

		assert.Error(t, err)
	})
}

func TestFindEnumForInt(t *testing.T) {
	t.Run("returns value when found", func(t *testing.T) {

		enum, err := FindEnumForInt(
			6,
			[]dummyIntEnum{setIntEnumValue},
			defaultIntEnumValue,
		)

		assert.NoError(t, err)
		assert.Equal(t, setIntEnumValue, enum)
	})

	t.Run("returns error when value is not found", func(t *testing.T) {

		_, err := FindEnumForInt(
			7,
			[]dummyIntEnum{setIntEnumValue},
			defaultIntEnumValue,
		)

		assert.Error(t, err)
	})
}
