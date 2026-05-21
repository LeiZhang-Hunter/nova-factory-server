package vectorsearch

import (
	"errors"
	"strings"
	"unicode"
)

type LabeledValue struct {
	Label string
	Value string
}

type ProcessedQuery struct {
	Original       string
	Normalized     string
	Tokens         []string
	ExpandedTokens []string
	Keywords       []string
	EmbeddingText  string
	HybridText     string
	IsCodeLike     bool
	IsShortQuery   bool
}

var defaultStopWords = map[string]struct{}{
	"的": {}, "了": {}, "和": {}, "及": {}, "与": {}, "或": {}, "及其": {},
	"产品": {}, "商品": {}, "一下": {}, "一个": {}, "一些": {}, "查询": {},
}

var defaultSynonyms = map[string][]string{
	"矿泉水": {"纯净水", "饮用水"},
	"纯净水": {"矿泉水", "饮用水"},
	"饮料":  {"饮品", "饮用"},
	"可乐":  {"汽水"},
	"牛奶":  {"鲜奶", "纯牛奶"},
	"酸奶":  {"发酵乳", "乳饮料"},
	"方便面": {"泡面"},
	"纸巾":  {"抽纸", "面纸", "纸抽"},
	"湿巾":  {"湿纸巾"},
	"洗衣液": {"洗衣剂"},
	"洗洁精": {"洗涤剂"},
	"香皂":  {"肥皂"},
	"大米":  {"稻米", "米"},
	"食用油": {"调和油", "菜籽油"},
}

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

func ProcessQuery(query string) *ProcessedQuery {
	original := NormalizeWhitespace(query)
	normalized := strings.ToLower(original)
	tokens := tokenize(normalized)
	expanded := expandTokens(tokens)
	keywords := dedupeKeywords(append(append([]string{normalized}, tokens...), expanded...))
	return &ProcessedQuery{
		Original:       original,
		Normalized:     normalized,
		Tokens:         tokens,
		ExpandedTokens: expanded,
		Keywords:       keywords,
		EmbeddingText:  buildQueryText(original, tokens, expanded, 2),
		HybridText:     buildQueryText(original, tokens, expanded, 1),
		IsCodeLike:     isCodeLike(normalized),
		IsShortQuery:   len([]rune(original)) <= 4 || len(tokens) <= 2,
	}
}

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

func tokenize(query string) []string {
	replaced := strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), isCJK(r):
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

func expandTokens(tokens []string) []string {
	expanded := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if synonyms, ok := defaultSynonyms[token]; ok {
			expanded = append(expanded, synonyms...)
		}
	}
	return dedupeKeywords(expanded)
}

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

func buildQueryText(original string, tokens, expanded []string, duplicateWeight int) string {
	parts := make([]string, 0, 1+len(tokens)+len(expanded)*duplicateWeight)
	if original != "" {
		parts = append(parts, original)
	}
	parts = append(parts, tokens...)
	for _, token := range expanded {
		for i := 0; i < duplicateWeight; i++ {
			parts = append(parts, token)
		}
	}
	return strings.Join(dedupeKeywords(parts), " ")
}

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

func containsCJK(value string) bool {
	for _, r := range value {
		if isCJK(r) {
			return true
		}
	}
	return false
}

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

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r)
}
