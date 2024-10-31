package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOffset(t *testing.T) {
	t.Run(
		"can not increment when offset is equal to totalCount",
		func(t *testing.T) {
			got := NewOffset(10, 5, 5)

			assert.Nil(t, got)
		},
	)

	t.Run(
		"can not increment when offset + limit passes totalCount",
		func(t *testing.T) {
			got := NewOffset(7, 3, 8)

			assert.Nil(t, got)
		},
	)

	t.Run(
		"can increment when offset + limit is less than totalCount",
		func(t *testing.T) {
			got := NewOffset(2, 3, 10)

			assert.Equal(t, int32(5), *got)
		},
	)
}

func ExampleNewOffset() {
	offset := NewOffset(0, 5, 12)
	fmt.Println(*offset)
	// Output:
	// 5
}

func ExampleNewOffset_second() {
	offset := NewOffset(10, 12, 5)
	fmt.Println(offset)
	// Output:
	// <nil>
}
