package shared

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

func NewIntMarkdownList(values []int) IntMarkdownList {
	i := IntMarkdownList{}
	for _, value := range values {
		i = append(i, value)
	}

	return i
}
