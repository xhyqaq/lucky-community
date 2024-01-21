package utils

import (
	"html/template"
	"strings"
	"unicode/utf8"
)

func GlobalFunc() template.FuncMap {
	return template.FuncMap{
		"DiffForHumans":    DiffForHumans,
		"ToDateTimeString": ToDateTimeString,
		"Html":             Html,
		"RemindName":       RemindName,
		"StrLimit":         Limit,
		"StrJoin":          strings.Join,
	}
}

// Limit 将字符串以指定长度进行截断
func Limit(s string, start, length int, append string) string {
	strLen := utf8.RuneCountInString(s)
	if strLen <= 0 {
		return ""
	}

	if length >= strLen {
		length = strLen
		append = ""
	}

	runes := []rune(s)

	return string(runes[start:length]) + append
}

// Html 解析 HTML
func Html(s string) template.HTML {
	return template.HTML(s)
}

// RemindName 提醒名
func RemindName(a string) string {
	m := map[string]string{
		"comment:topic": "评论了你的话题",
		"reply:comment": "回复了你的评论",
		"like:topic":    "赞了你的话题",
		"like:comment":  "赞了你的评论",
		"follow:user":   "关注了你",
	}
	if v, ok := m[a]; !ok {
		return ""
	} else {
		return v
	}
}
