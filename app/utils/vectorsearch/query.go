package vectorsearch

import (
	"errors"
	"strings"
)

// NormalizeQueries 对批量查询做基础清洗，并过滤空查询。
func NormalizeQueries(queries []string) ([]string, error) {
	result := make([]string, 0, len(queries))
	for _, query := range queries {
		query = NormalizeWhitespace(query)
		if query == "" {
			continue
		}
		result = append(result, query)
	}
	if len(result) == 0 {
		return nil, errors.New("搜索内容不能为空")
	}
	return result, nil
}

// ProcessQueries 批量处理查询，输出结构化 query 结果。
func ProcessQueries(queries []string) ([]*ProcessedQuery, error) {
	normalized, err := NormalizeQueries(queries)
	if err != nil {
		return nil, err
	}
	result := make([]*ProcessedQuery, 0, len(normalized))
	for _, query := range normalized {
		result = append(result, ProcessQuery(query))
	}
	return result, nil
}

// ProcessQuery 将单条原始查询转换为结构化检索表达。
func ProcessQuery(query string) *ProcessedQuery {
	original := NormalizeWhitespace(query)
	normalized := strings.ToLower(original)
	tokens := tokenize(normalized)
	expanded := expandTokens(tokens)
	categoryTerms := extractCategoryTerms(normalized, tokens)
	specTerms := extractSpecTerms(normalized, tokens)
	codeTerms := extractCodeTerms(normalized, tokens)
	keywords := dedupeKeywords(append(append(append(append([]string{normalized}, tokens...), expanded...), categoryTerms...), specTerms...))
	return &ProcessedQuery{
		Original:       original,
		Normalized:     normalized,
		Tokens:         tokens,
		ExpandedTokens: expanded,
		Keywords:       keywords,
		CategoryTerms:  categoryTerms,
		SpecTerms:      specTerms,
		CodeTerms:      codeTerms,
		EmbeddingText:  buildQueryText(original, tokens, expanded, categoryTerms, specTerms, codeTerms, 2),
		HybridText:     buildQueryText(original, tokens, expanded, categoryTerms, specTerms, codeTerms, 1),
		IsCodeLike:     isCodeLike(normalized),
		IsShortQuery:   len([]rune(original)) <= 4 || len(tokens) <= 2,
	}
}
