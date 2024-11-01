package utils

import (
	"strconv"
)

// IntMarkdownList implements helpers to use int64 sets in validation & documentation.
type IntMarkdownList []int

// Markdown returns a string with all the values in Markdown list format.
func (i IntMarkdownList) Markdown() string {
	markdown := "\n"
	for _, i := range i {
		markdown += "  - *" + strconv.Itoa(i) + "*\n"
	}

	return markdown
}

// ToInt32 converts all slice values to int32.
func (i IntMarkdownList) ToInt32() []int32 {
	var returnValues []int32

	for _, i := range i {
		returnValues = append(returnValues, int32(i))
	}

	return returnValues
}

// NewIntMarkdownList instantiates a new IntMarkdownList.
func NewIntMarkdownList[T ~int32](values []T) IntMarkdownList {
	i := IntMarkdownList{}
	for _, value := range values {
		i = append(i, int(value))
	}

	return i
}

// StringTypeArrayToMarkdown converts any slice of string enums
// that are underlying string types to a Markdown list string.
func StringTypeArrayToMarkdown[T ~string](enumValues []T) string {
	markdown := "\n"
	for _, v := range enumValues {
		markdown += "  - *" + string(v) + "*\n"
	}
	return markdown
}
