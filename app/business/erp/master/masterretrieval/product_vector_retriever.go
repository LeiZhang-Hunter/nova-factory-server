package masterretrieval

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/utils/retrieval"
	"nova-factory-server/app/utils/retrieval/schema"
	searchutil "nova-factory-server/app/utils/vectorsearch"
)

const (
	defaultProductRetrieverTopK = 10
	maxProductRetrieverTopK     = 50
)

// ProductVectorRetriever 将 ERP 商品向量检索适配为通用 Retriever 接口。
type ProductVectorRetriever struct {
	dao masterdao.IProductVectorDao
}

// NewProductVectorRetriever 创建商品向量检索器。
func NewProductVectorRetriever(dao masterdao.IProductVectorDao) retrieval.Retriever {
	return &ProductVectorRetriever{dao: dao}
}

// Retrieve 执行单条商品检索，并将业务结果映射为统一文档结构。
func (r *ProductVectorRetriever) Retrieve(ctx context.Context, query string, opts ...retrieval.Option) ([]*schema.Document, error) {
	if r == nil || r.dao == nil {
		return nil, errors.New("商品检索器未初始化")
	}

	processed := searchutil.ProcessQuery(query)
	if processed == nil || processed.Original == "" {
		return nil, errors.New("搜索内容不能为空")
	}

	options := retrieval.ApplyOptions(opts...)
	if options.Embedding == nil {
		return nil, errors.New("缺少 embedding 组件")
	}

	topK := resolveProductRetrieverTopK(options.TopK)
	vectors, err := options.Embedding.EmbedStrings(ctx, []string{processed.EmbeddingText})
	if err != nil {
		return nil, fmt.Errorf("生成搜索向量失败: %w", err)
	}
	if len(vectors) != 1 || len(vectors[0]) == 0 {
		return nil, errors.New("生成搜索向量失败: 返回结果为空")
	}

	data, err := r.dao.BatchSearch(nil, &mastermodels.ProductVectorBatchSearchReq{
		Queries:     []string{processed.Original},
		SearchTexts: []string{processed.HybridText},
		Limit:       topK,
	}, [][]float32{float64SliceToFloat32(vectors[0])})
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Rows) == 0 || data.Rows[0] == nil || len(data.Rows[0].Rows) == 0 {
		return make([]*schema.Document, 0), nil
	}

	rows := data.Rows[0].Rows
	documents := make([]*schema.Document, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		score := float64(row.Score)
		if options.ScoreThreshold != nil && score < *options.ScoreThreshold {
			continue
		}
		documents = append(documents, toProductVectorDocument(row))
	}
	return documents, nil
}

func toProductVectorDocument(item *mastermodels.ProductVectorSearchItem) *schema.Document {
	if item == nil {
		return nil
	}
	return &schema.Document{
		ID:      strconv.FormatInt(item.ProductID, 10),
		Content: item.Content,
		Score:   float64(item.Score),
		Metadata: map[string]any{
			"product": item,
			"name":    item.Name,
			"barCode": item.BarCode,
		},
	}
}

func resolveProductRetrieverTopK(topK *int) int {
	if topK == nil || *topK <= 0 {
		return defaultProductRetrieverTopK
	}
	if *topK > maxProductRetrieverTopK {
		return maxProductRetrieverTopK
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
