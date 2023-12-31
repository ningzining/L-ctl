package caseutil

import "strings"

const (
	UnderLine = '_'
	A         = 'A'
	Z         = 'Z'
)

// UpperCamelCase 转换大驼峰命名
// source 待转换的字符串
func UpperCamelCase(source string) string {
	if len(source) == 0 {
		return source
	}
	var sb strings.Builder
	change := true
	for _, s := range source {
		if s == UnderLine {
			change = true
		} else if change {
			toUpper := strings.ToUpper(string(s))
			sb.WriteString(toUpper)
			change = false
		} else {
			sb.WriteRune(s)
		}
	}
	return sb.String()
}

// LowerCamelCase 转换小驼峰命名
// source 待转换的字符串
func LowerCamelCase(source string) string {
	if len(source) == 0 {
		return source
	}
	var sb strings.Builder
	change := false
	for _, s := range source {
		if s == UnderLine {
			change = true
		} else if change {
			upper := strings.ToUpper(string(s))
			sb.WriteString(upper)
			change = false
		} else {
			sb.WriteRune(s)
		}
	}
	return sb.String()
}

// UnderLineCase 转换下划线命名
// source 待转换的字符串
func UnderLineCase(source string) string {
	if len(source) == 0 {
		return source
	}
	var sb strings.Builder
	for _, s := range source {
		if s >= A && s <= Z {
			if sb.Len() > 0 {
				sb.WriteString(string(UnderLine))
			}
			lower := strings.ToLower(string(s))
			sb.WriteString(lower)
		} else {
			sb.WriteRune(s)
		}
	}
	return sb.String()
}
