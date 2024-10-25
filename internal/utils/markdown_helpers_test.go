package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntMarkdownList_Markdown(t *testing.T) {
	t.Run("returns valid markdown", func(t *testing.T) {
		got := IntMarkdownList{1, 2, 3}.Markdown()
		want := "\n  - *1*\n  - *2*\n  - *3*\n"
		assert.Equal(t, want, got)
	})
}

func ExampleIntMarkdownList_Markdown() {
	list := IntMarkdownList{1, 2, 3}

	fmt.Println(list.Markdown())
	/**
	  Output:
	  - *1*
	  - *2*
	  - *3*
	*/
}

func TestIntMarkdownList_ToInt64(t *testing.T) {
	got := IntMarkdownList{1, 2, 3}.ToInt64()
	want := []int64{1, 2, 3}

	assert.Equal(t, want, got)
}

func ExampleIntMarkdownList_ToInt64() {
	list := IntMarkdownList{1, 2, 3}.ToInt64()

	fmt.Println(list)
	// Output []{1, 2, 3}
}

func TestNewIntMarkdownList(t *testing.T) {
	got := NewIntMarkdownList([]int32{1})

	assert.Equal(t, []int64{1}, got.ToInt64())
}

func TestStringTypeArrayToMarkdown(t *testing.T) {
	type underlyingString string

	enumValues := []underlyingString{
		"TEST_ONE",
		"TEST_TWO",
	}

	want := "\n  - *TEST_ONE*\n  - *TEST_TWO*\n"
	got := StringTypeArrayToMarkdown(enumValues)
	assert.Equal(t, want, got)
}

func ExampleStringTypeArrayToMarkdown() {
	type underlyingString string

	enumValues := []underlyingString{
		"TEST_ONE",
		"TEST_TWO",
	}

	markdown := StringTypeArrayToMarkdown(enumValues)

	fmt.Println(markdown)
	// Output "\n  - *TEST_ONE*\n  - *TEST_TWO*\n"
}
