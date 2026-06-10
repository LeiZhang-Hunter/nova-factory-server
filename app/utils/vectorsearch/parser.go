package vectorsearch

import "strings"

// tokenize 对 query 做基础切词，并补充中文 n-gram 片段。
func tokenize(query string) []string {
	replaced := strings.Map(func(r rune) rune {
		switch {
		case isCJK(r), isAlphaNumeric(r):
			return r
		default:
			return ' '
		}
	}, query)

	parts := strings.Fields(replaced)
	tokens := make([]string, 0, len(parts)*2)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := defaultStopWords[part]; ok {
			continue
		}
		tokens = append(tokens, part)
		if containsCJK(part) {
			tokens = append(tokens, buildCJKNGrams(part)...)
		}
	}
	return dedupeKeywords(tokens)
}

// expandTokens 根据内置同义词词典补充扩展词。
func expandTokens(tokens []string) []string {
	expanded := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if synonyms, ok := defaultSynonyms[token]; ok {
			expanded = append(expanded, synonyms...)
		}
	}
	return dedupeKeywords(expanded)
}

// dedupeKeywords 负责统一大小写、过滤停用词并去重。
func dedupeKeywords(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	seen := make(map[string]struct{}, len(tokens))
	for _, token := range tokens {
		token = NormalizeWhitespace(strings.ToLower(token))
		if token == "" {
			continue
		}
		if _, ok := defaultStopWords[token]; ok {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		result = append(result, token)
	}
	return result
}

// buildQueryText 组合生成给 embedding 与 hybrid 检索使用的查询文本。
func buildQueryText(original string, tokens, expanded, categoryTerms []string, duplicateWeight int) string {
	parts := make([]string, 0, 1+len(tokens)+len(expanded)*duplicateWeight+len(categoryTerms))
	if original != "" {
		parts = append(parts, original)
	}
	parts = append(parts, tokens...)
	parts = append(parts, categoryTerms...)
	for _, token := range expanded {
		for i := 0; i < duplicateWeight; i++ {
			parts = append(parts, token)
		}
	}
	return strings.Join(dedupeKeywords(parts), " ")
}

// buildCJKNGrams 为中文短词构造 2-gram / 3-gram 片段，增强召回。
func buildCJKNGrams(value string) []string {
	runes := []rune(value)
	if len(runes) <= 1 || len(runes) > 8 {
		return nil
	}
	result := make([]string, 0, len(runes)*2)
	for size := 2; size <= 3; size++ {
		if len(runes) < size {
			break
		}
		for i := 0; i+size <= len(runes); i++ {
			result = append(result, string(runes[i:i+size]))
		}
	}
	return result
}

// isAlphaNumeric 判断 rune 是否为英文字母或数字。
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
