package utils

// EnumToSlice converts any slice of custom types (enums) that are underlying string types to a slice of strings.
func EnumToSlice[T ~string](enumValues []T) []string {
	strValues := make([]string, len(enumValues))
	for i, v := range enumValues {
		strValues[i] = string(v)
	}
	return strValues
}

// EnumToMarkdown converts any slice of custom types (enums) that are underlying string types to markdown string.
func EnumToMarkdown[T ~string](enumValues []T) string {
	markdown := "\n"
	for _, v := range enumValues {
		markdown += "  - *" + string(v) + "*\n"
	}
	return markdown
}
