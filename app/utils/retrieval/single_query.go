package retrieval

import (
	"context"
	"errors"
	"fmt"

	"nova-factory-server/app/utils/retrieval/schema"
	searchutil "nova-factory-server/app/utils/vectorsearch"
)

// QueryPayload 表示单条检索请求在预处理后的关键信息。
type QueryPayload struct {
	Original   string
	HybridText string
}

// SingleQueryRetrieveFunc 定义单条检索执行函数。
//
// 调用方只需要关心：
// 1. 如何使用处理后的 query 与向量访问底层存储
// 2. 如何把底层结果映射为统一文档
type SingleQueryRetrieveFunc func(payload QueryPayload, topK int, vector []float32) ([]*schema.Document, error)

// RetrieveSingleQueryWithEmbedding 统一处理单条检索中的公共步骤。
//
// 公共步骤包括：
// 1. query 预处理
// 2. embedding 组件校验
// 3. 生成查询向量
// 4. 规范化 TopK
// 5. 结果阈值过滤
func RetrieveSingleQueryWithEmbedding(
	ctx context.Context,
	query string,
	defaultTopK int,
	maxTopK int,
	run SingleQueryRetrieveFunc,
	opts ...Option,
) ([]*schema.Document, error) {
	if run == nil {
		return nil, errors.New("检索执行函数不能为空")
	}

	processed := searchutil.ProcessQuery(query)
	if processed == nil || processed.Original == "" {
		return nil, errors.New("搜索内容不能为空")
	}

	options := ApplyOptions(opts...)
	if options.Embedding == nil {
		return nil, errors.New("缺少 embedding 组件")
	}

	topK := resolveTopK(options.TopK, defaultTopK, maxTopK)
	vectors, err := options.Embedding.EmbedStrings(ctx, []string{processed.EmbeddingText})
	if err != nil {
		return nil, fmt.Errorf("生成搜索向量失败: %w", err)
	}
	if len(vectors) != 1 || len(vectors[0]) == 0 {
		return nil, errors.New("生成搜索向量失败: 返回结果为空")
	}

	documents, err := run(QueryPayload{
		Original:   processed.Original,
		HybridText: processed.HybridText,
	}, topK, float64SliceToFloat32(vectors[0]))
	if err != nil {
		return nil, err
	}
	return filterDocumentsByScore(documents, options.ScoreThreshold), nil
}

func resolveTopK(topK *int, defaultTopK int, maxTopK int) int {
	if topK == nil || *topK <= 0 {
		return defaultTopK
	}
	if maxTopK > 0 && *topK > maxTopK {
		return maxTopK
	}
	return *topK
}

func float64SliceToFloat32(values []float64) []float32 {
	result := make([]float32, len(values))
	for idx, value := range values {
		result[idx] = float32(value)
	}
	return result
}

func filterDocumentsByScore(documents []*schema.Document, threshold *float64) []*schema.Document {
	if len(documents) == 0 {
		return make([]*schema.Document, 0)
	}
	if threshold == nil {
		return documents
	}
	filtered := make([]*schema.Document, 0, len(documents))
	for _, document := range documents {
		if document == nil || document.Score < *threshold {
			continue
		}
		filtered = append(filtered, document)
	}
	return filtered
}
