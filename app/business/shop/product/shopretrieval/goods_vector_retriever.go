package shopretrieval

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/utils/retrieval"
	"nova-factory-server/app/utils/retrieval/schema"
	searchutil "nova-factory-server/app/utils/vectorsearch"
)

const (
	defaultGoodsRetrieverTopK = 10
	maxGoodsRetrieverTopK     = 50
)

// GoodsVectorRetriever 将商城商品向量检索适配为通用 Retriever 接口。
type GoodsVectorRetriever struct {
	dao shopdao.IShopGoodsVectorDao
}

// NewGoodsVectorRetriever 创建商城商品检索器。
func NewGoodsVectorRetriever(dao shopdao.IShopGoodsVectorDao) retrieval.Retriever {
	return &GoodsVectorRetriever{dao: dao}
}

// Retrieve 执行单条商品检索，并输出统一文档结构。
func (r *GoodsVectorRetriever) Retrieve(ctx context.Context, query string, opts ...retrieval.Option) ([]*schema.Document, error) {
	if r == nil || r.dao == nil {
		return nil, errors.New("商城商品检索器未初始化")
	}

	processed := searchutil.ProcessQuery(query)
	if processed == nil || processed.Original == "" {
		return nil, errors.New("搜索内容不能为空")
	}

	options := retrieval.ApplyOptions(opts...)
	if options.Embedding == nil {
		return nil, errors.New("缺少 embedding 组件")
	}

	topK := resolveGoodsRetrieverTopK(options.TopK)
	vectors, err := options.Embedding.EmbedStrings(ctx, []string{processed.EmbeddingText})
	if err != nil {
		return nil, fmt.Errorf("生成搜索向量失败: %w", err)
	}
	if len(vectors) != 1 || len(vectors[0]) == 0 {
		return nil, errors.New("生成搜索向量失败: 返回结果为空")
	}

	data, err := r.dao.BatchSearch(nil, &shopmodels.GoodsVectorBatchSearchReq{
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
		documents = append(documents, toGoodsVectorDocument(row))
	}
	return documents, nil
}

func toGoodsVectorDocument(item *shopmodels.GoodsVectorSearchItem) *schema.Document {
	if item == nil {
		return nil
	}
	return &schema.Document{
		ID:      strconv.FormatInt(item.GoodsDBID, 10),
		Content: item.Content,
		Score:   float64(item.Score),
		Metadata: map[string]any{
			"goods":     item,
			"goodsName": item.GoodsName,
			"goodsCode": item.GoodsCode,
			"skuId":     item.SkuID,
			"skuName":   item.SkuName,
		},
	}
}

func resolveGoodsRetrieverTopK(topK *int) int {
	if topK == nil || *topK <= 0 {
		return defaultGoodsRetrieverTopK
	}
	if *topK > maxGoodsRetrieverTopK {
		return maxGoodsRetrieverTopK
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
