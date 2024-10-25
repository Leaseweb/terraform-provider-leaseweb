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

// ToInt64 converts all slice values to int64.
func (i IntMarkdownList) ToInt64() []int64 {
	var returnValues []int64

	for _, i := range i {
		returnValues = append(returnValues, int64(i))
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

// StringTypeArrayToMarkdown converts any slice of custom types that are underlying string types to markdown string.
func StringTypeArrayToMarkdown[T ~string](enumValues []T) string {
	markdown := "\n"
	for _, v := range enumValues {
		markdown += "  - *" + string(v) + "*\n"
	}
	return markdown
}
