package shopretrieval

import (
	"context"
	"errors"
	"strconv"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/utils/retrieval"
	"nova-factory-server/app/utils/retrieval/schema"
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
	return retrieval.RetrieveSingleQueryWithEmbedding(ctx, query, defaultGoodsRetrieverTopK, maxGoodsRetrieverTopK,
		func(payload retrieval.QueryPayload, topK int, vector []float32) ([]*schema.Document, error) {
			data, err := r.dao.BatchSearch(nil, &shopmodels.GoodsVectorBatchSearchReq{
				Queries:     []string{payload.Original},
				SearchTexts: []string{payload.HybridText},
				Limit:       topK,
			}, [][]float32{vector})
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
				documents = append(documents, toGoodsVectorDocument(row))
			}
			return documents, nil
		},
		opts...)
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
