package caseutil

import "strings"

const (
	UnderLine = '_'
)

// ToCamelCase 转换驼峰命名
// source 待转换的字符串
// isTitleCase 首字母是否大写
func ToCamelCase(source string, isTitleCase bool) string {
	if len(source) == 0 {
		return source
	}
	var sb strings.Builder
	upper := isTitleCase
	for _, s := range source {
		if s == UnderLine {
			upper = true
		} else if upper {
			toUpper := strings.ToUpper(string(s))
			sb.WriteString(toUpper)
			upper = false
		} else {
			sb.WriteRune(s)
		}
	}
	return sb.String()
}
