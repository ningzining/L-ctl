package caseutil

import "strings"

const (
	UnderLine = '_'
	A         = 'A'
	Z         = 'Z'
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

// ToUnderLineCase 转换下划线命名
// source 待转换的字符串
func ToUnderLineCase(source string) string {
	if len(source) == 0 {
		return source
	}
	var sb strings.Builder
	for _, s := range source {
		if s > A && s < Z {
			sb.WriteString(string(UnderLine))
			toLower := strings.ToLower(string(s))
			sb.WriteString(toLower)
		} else {
			sb.WriteRune(s)
		}
	}
	return sb.String()
}
