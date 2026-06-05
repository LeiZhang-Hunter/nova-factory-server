package vectorsearch

import (
	"errors"
	"strings"

	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization"
	normalizationapi "nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/regex"
	replacepkg "nova-factory-server/app/utils/vectorsearch/normalization/replace"
	"nova-factory-server/app/utils/vectorsearch/normalization/shopcategory"
	"nova-factory-server/app/utils/vectorsearch/normalization/whitespace"
)

// GoodsSearchQueryPayload 保存原始查询词和送入向量检索的增强查询词。
type GoodsSearchQueryPayload struct {
	Original   string
	SearchText string
}

// BuildGoodsSearchQueryPayloads 将原始商品名做清洗，并提取分类、规格等结构化信息，
// 组装为更适合向量检索理解的增强查询串。
func BuildGoodsSearchQueryPayloads(goodsNames []string) ([]GoodsSearchQueryPayload, error) {
	pipeline := normalization.NewPipeline(normalizationapi.Config{
		Interceptors: []*normalizationapi.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: shopcategory.Type},
		},
	})
	specPipeline, err := newGoodsSearchSpecPipeline()
	if err != nil {
		return nil, err
	}

	payloads := make([]GoodsSearchQueryPayload, 0, len(goodsNames))
	for _, goodsName := range goodsNames {
		original := NormalizeWhitespace(goodsName)
		if original == "" {
			continue
		}

		result, err := pipeline.Normalize(original)
		if err != nil {
			return nil, err
		}
		specResult, err := specPipeline.Normalize(original)
		if err != nil {
			return nil, err
		}

		payloads = append(payloads, GoodsSearchQueryPayload{
			Original:   original,
			SearchText: buildGoodsSearchQueryText(result, specResult),
		})
	}
	if len(payloads) == 0 {
		return nil, errors.New("商品名称不能为空")
	}
	return payloads, nil
}

// GoodsSearchQueryTexts 提取最终送入向量检索的查询串。
func GoodsSearchQueryTexts(payloads []GoodsSearchQueryPayload) []string {
	result := make([]string, 0, len(payloads))
	for _, payload := range payloads {
		if payload.SearchText == "" {
			continue
		}
		result = append(result, payload.SearchText)
	}
	return result
}

// newGoodsSearchSpecPipeline 构造规格提取 pipeline。
// 先统一常见的规格分隔符写法，再用正则截取规格主体。
func newGoodsSearchSpecPipeline() (*normalization.Pipeline, error) {
	replaceProperties, err := cfg.Pack(replacepkg.Config{
		Old: "x",
		New: "*",
	})
	if err != nil {
		return nil, err
	}
	regexProperties, err := cfg.Pack(regex.Config{
		Pattern:     goodsSearchSpecPattern,
		Replacement: "$1",
	})
	if err != nil {
		return nil, err
	}
	return normalization.NewPipeline(normalizationapi.Config{
		Interceptors: []*normalizationapi.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: replacepkg.Type, Properties: replaceProperties},
			{Type: regex.Type, Properties: regexProperties},
		},
	}), nil
}

// buildGoodsSearchQueryText 将归一化主文本、分类名和规格词拼接为增强检索串。
func buildGoodsSearchQueryText(result normalization.Result, specResult normalization.Result) string {
	parts := make([]string, 0, 1+len(result.Categories)+len(specResult.Matches))
	parts = append(parts, NormalizeWhitespace(result.Value))
	for _, category := range result.Categories {
		if category.Name == "" {
			continue
		}
		parts = append(parts, NormalizeWhitespace(category.Name))
	}
	parts = append(parts, extractGoodsSearchSpecTerms(specResult)...)
	return strings.Join(dedupeGoodsSearchParts(parts), " ")
}

// extractGoodsSearchSpecTerms 优先复用通用 query 处理逻辑中的规格提取结果；
// 如果没有抽到规格，再回退使用 regex step 的输出值。
func extractGoodsSearchSpecTerms(result normalization.Result) []string {
	if processed := ProcessQuery(result.Value); processed != nil && len(processed.SpecTerms) > 0 {
		return processed.SpecTerms
	}
	for _, match := range result.Matches {
		if match.Kind == "regex_replace" {
			value := NormalizeWhitespace(result.Value)
			if value != "" {
				return []string{value}
			}
			break
		}
	}
	return nil
}

// dedupeGoodsSearchParts 对增强查询串片段去重，避免重复词过度放大权重。
func dedupeGoodsSearchParts(parts []string) []string {
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = NormalizeWhitespace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		result = append(result, part)
	}
	return result
}
