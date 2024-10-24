package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumToSlice(t *testing.T) {
	type underlyingString string

	enumValues := []underlyingString{
		"TEST_ONE",
		"TEST_TWO",
	}

	want := []string{
		"TEST_ONE",
		"TEST_TWO",
	}

	got := EnumToSlice(enumValues)
	assert.Equal(t, want, got)
}

func TestEnumToMarkdown(t *testing.T) {
	type underlyingString string

	enumValues := []underlyingString{
		"TEST_ONE",
		"TEST_TWO",
	}

	want := "\n  - *TEST_ONE*\n  - *TEST_TWO*\n"
	got := EnumToMarkdown(enumValues)
	assert.Equal(t, want, got)
}
