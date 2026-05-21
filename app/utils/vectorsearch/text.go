package vectorsearch

import (
	"strings"
	"unicode"
)

// BuildLabeledContent 将结构化字段拼装为带标签的文本内容。
func BuildLabeledContent(values []LabeledValue, max int) string {
	lines := make([]string, 0, len(values))
	for _, item := range values {
		value := NormalizeWhitespace(item.Value)
		if value == "" {
			continue
		}
		label := NormalizeWhitespace(item.Label)
		if label == "" {
			lines = append(lines, value)
			continue
		}
		lines = append(lines, label+": "+value)
	}
	return TrimRunes(strings.Join(lines, "\n"), max)
}

// TrimRunes 按 rune 长度裁剪文本，避免中文被截断成乱码。
func TrimRunes(value string, max int) string {
	value = NormalizeWhitespace(value)
	if max <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}

// NormalizeWhitespace 统一处理多余空格、换行与首尾空白。
func NormalizeWhitespace(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	var builder strings.Builder
	lastSpace := false
	for _, r := range value {
		if unicode.IsSpace(r) {
			if lastSpace {
				continue
			}
			builder.WriteByte(' ')
			lastSpace = true
			continue
		}
		builder.WriteRune(r)
		lastSpace = false
	}
	return strings.TrimSpace(builder.String())
}

// containsCJK 判断字符串中是否包含中文字符。
func containsCJK(value string) bool {
	for _, r := range value {
		if isCJK(r) {
			return true
		}
	}
	return false
}

// isCodeLike 判断 query 是否更像编码、SKU、条码等精确标识。
func isCodeLike(value string) bool {
	if value == "" || containsCJK(value) {
		return false
	}
	compact := strings.ReplaceAll(value, " ", "")
	if len(compact) < 4 {
		return false
	}
	hasDigit := false
	for _, r := range compact {
		if unicode.IsDigit(r) {
			hasDigit = true
			continue
		}
		if unicode.IsLetter(r) || r == '-' || r == '_' || r == '/' {
			continue
		}
		return false
	}
	return hasDigit
}

// isCJK 判断 rune 是否属于中文汉字。
func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r)
}
