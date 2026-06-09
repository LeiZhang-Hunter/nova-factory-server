package vectorsearch

import (
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/store"
	"strings"
)

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

const maxCategoryTreeDepth = 8 // 分类树最大遍历深度，避免无限递归

func matchCategoryTerms(value string) []string {
	value = strings.ToLower(NormalizeWhitespace(value))
	if value == "" {
		return nil
	}
	result := make([]string, 0, 2)

	categoryStore := store.GetStore(shop.ShopCategoryStoreName)

	defaultCategoryTerms, exist := categoryStore.Get()
	if !exist {
		return result
	}

	for _, term := range defaultCategoryTerms {
		matchCategoryTermsRecursive(term, value, &result, 0)
	}
	return dedupeKeywords(result)
}

// matchCategoryTermsRecursive 递归遍历分类树节点，将命中分类名追加到 result。
// depth 从 0 开始计数，超过 maxCategoryTreeDepth 时停止向下递归。
func matchCategoryTermsRecursive(term store.ShopCategoryData, value string, result *[]string, depth int) {
	if term == nil || depth > maxCategoryTreeDepth {
		return
	}
	name := strings.ToLower(NormalizeWhitespace(term.Name()))
	if name != "" && (name == value || strings.Contains(value, name)) {
		*result = append(*result, name)
	}
	for _, child := range term.ChildrenData() {
		matchCategoryTermsRecursive(child, value, result, depth+1)
	}
}
