package masterretrieval

import (
	"context"
	"errors"
	"strconv"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/utils/retrieval"
	"nova-factory-server/app/utils/retrieval/schema"
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
	return retrieval.RetrieveSingleQueryWithEmbedding(ctx, query, defaultProductRetrieverTopK, maxProductRetrieverTopK,
		func(payload retrieval.QueryPayload, topK int, vector []float32) ([]*schema.Document, error) {
			data, err := r.dao.BatchSearch(nil, &mastermodels.ProductVectorBatchSearchReq{
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
				documents = append(documents, toProductVectorDocument(row))
			}
			return documents, nil
		},
		opts...)
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
