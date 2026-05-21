package vectorsearch

import (
	"regexp"
	"strconv"
	"strings"
)

// specTokenPattern 用于识别 query 中的规格表达，如 550ml、0.55L、5kg 等。
var specTokenPattern = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(ml|l|g|kg|mg|斤|两|箱|盒|包|袋|瓶|支|桶|听|罐|片)`)

// extractSpecTerms 从 query 和 token 中抽取并标准化规格词。
func extractSpecTerms(query string, tokens []string) []string {
	result := make([]string, 0, len(tokens)+2)
	if query != "" {
		result = append(result, findNormalizedSpecs(query)...)
	}
	for _, token := range tokens {
		result = append(result, findNormalizedSpecs(token)...)
	}
	return dedupeKeywords(result)
}

// extractCodeTerms 从 query 和 token 中抽取编码类词。
func extractCodeTerms(query string, tokens []string) []string {
	result := make([]string, 0, len(tokens)+1)
	if isCodeLike(query) {
		result = append(result, compactCode(query))
	}
	for _, token := range tokens {
		if isCodeLike(token) {
			result = append(result, compactCode(token))
		}
	}
	return dedupeKeywords(result)
}

// findNormalizedSpecs 从文本中提取所有规格，并归一化到统一单位。
func findNormalizedSpecs(value string) []string {
	if value == "" {
		return nil
	}
	matches := specTokenPattern.FindAllStringSubmatch(strings.ToLower(value), -1)
	if len(matches) == 0 {
		return nil
	}
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		normalized := normalizeSpec(match[1], match[2])
		if normalized != "" {
			result = append(result, normalized)
		}
	}
	return dedupeKeywords(result)
}

// normalizeSpec 将不同单位表达映射为统一格式，便于精确比较规格。
func normalizeSpec(numberText, unit string) string {
	numberText = strings.TrimSpace(numberText)
	unit = strings.ToLower(strings.TrimSpace(unit))
	if numberText == "" || unit == "" {
		return ""
	}

	value, err := strconv.ParseFloat(numberText, 64)
	if err != nil {
		return strings.ToLower(numberText + unit)
	}

	switch unit {
	case "l":
		value = value * 1000
		unit = "ml"
	case "kg":
		value = value * 1000
		unit = "g"
	case "斤":
		value = value * 500
		unit = "g"
	case "两":
		value = value * 50
		unit = "g"
	}

	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10) + unit
	}
	return strconv.FormatFloat(value, 'f', -1, 64) + unit
}

// compactCode 统一编码表达，去掉空格并转为小写。
func compactCode(value string) string {
	value = NormalizeWhitespace(strings.ToLower(value))
	return strings.ReplaceAll(value, " ", "")
}

// scoreSpecTerms 计算 query 规格词与候选字段规格词的匹配度。
func scoreSpecTerms(query *ProcessedQuery, fields ...string) float64 {
	if query == nil || len(query.SpecTerms) == 0 {
		return 0
	}
	fieldSpecs := make([]string, 0, len(fields)*2)
	for _, field := range fields {
		fieldSpecs = append(fieldSpecs, findNormalizedSpecs(strings.ToLower(field))...)
	}
	fieldSpecs = dedupeKeywords(fieldSpecs)
	if len(fieldSpecs) == 0 {
		return 0
	}

	matched := 0
	for _, spec := range query.SpecTerms {
		for _, fieldSpec := range fieldSpecs {
			if spec == fieldSpec {
				matched++
				break
			}
		}
	}
	return clampFloat(float64(matched)/float64(len(query.SpecTerms)), 0, 1)
}
