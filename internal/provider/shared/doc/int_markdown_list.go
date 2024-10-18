package doc

import (
	"strconv"
)

type IntMarkdownList []int

func (i IntMarkdownList) Markdown() string {
	markdown := "\n"
	for _, i := range i {
		markdown += "  - *" + strconv.Itoa(i) + "*\n"
	}

	return markdown
}

func (i IntMarkdownList) ToInt64() []int64 {
	var returnValues []int64

	for _, i := range i {
		returnValues = append(returnValues, int64(i))
	}

	return returnValues
}

func NewIntMarkdownList[T ~int32](values []T) IntMarkdownList {
	i := IntMarkdownList{}
	for _, value := range values {
		i = append(i, int(value))
	}

	return i
}
