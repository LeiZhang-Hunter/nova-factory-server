package vectorsearch

import "strings"

// extractCategoryTerms 从 query 和 token 中提取商品分类词。
func extractCategoryTerms(query string, tokens []string) []string {
	result := make([]string, 0, len(tokens))
	if query != "" {
		result = append(result, matchCategoryTerms(query)...)
	}
	for _, token := range tokens {
		result = append(result, matchCategoryTerms(token)...)
	}
	return dedupeKeywords(result)
}

func matchCategoryTerms(value string) []string {
	value = strings.ToLower(NormalizeWhitespace(value))
	if value == "" {
		return nil
	}
	result := make([]string, 0, 2)
	for term := range defaultCategoryTerms {
		if term == value || strings.Contains(value, term) {
			result = append(result, term)
		}
	}
	return dedupeKeywords(result)
}
