package masterdaoimpl

import (
	"context"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// executeProductVectorSearch 根据检索计划执行混合搜索或降级搜索。
func executeProductVectorSearch(
	ctx context.Context,
	client *milvusclient.Client,
	collectionName string,
	limit int,
	outputFields []string,
	plans []productVectorSearchPlan,
) ([]milvusclient.ResultSet, error) {
	supportHybrid, err := supportsProductVectorHybridSearch(ctx, client, collectionName)
	if err != nil {
		return nil, err
	}

	// 将策略层生成的检索计划统一转换为 Milvus 所需的向量与文本输入。
	searchVectors := make([]entity.Vector, 0, len(plans))
	textQueries := make([]entity.Vector, 0, len(plans))
	for _, plan := range plans {
		searchVectors = append(searchVectors, plan.Vector)
		textQueries = append(textQueries, entity.Text(plan.SearchText))
	}

	if supportHybrid {
		// HybridSearch 同时走稠密向量和 BM25 稀疏召回，通过 RRF 进行融合。
		sparseAnnParam := index.NewSparseAnnParam()
		sparseAnnParam.WithDropRatio(0.2)

		denseRequest := milvusclient.NewAnnRequest(productVectorEmbeddingField, limit, searchVectors...)
		sparseRequest := milvusclient.NewAnnRequest(productVectorContextSparseField, limit, textQueries...).
			WithAnnParam(sparseAnnParam)

		//ranker := entity.NewFunction().
		//	WithName("rrf").
		//	WithType(entity.FunctionTypeRerank).
		//	WithParam("reranker", "rrf").
		//	WithParam("k", "100")

		return client.HybridSearch(ctx, milvusclient.NewHybridSearchOption(
			collectionName,
			limit,
			denseRequest,
			sparseRequest,
		).WithReranker(milvusclient.NewWeightedReranker([]float64{
			0.3, 0.8,
		})).
			WithOutputFields(outputFields...))
	}

	// 不支持 hybrid 的旧 collection 降级为单路向量检索，兼容历史数据。
	return client.Search(ctx, milvusclient.NewSearchOption(collectionName, limit,
		searchVectors).
		WithANNSField(productVectorEmbeddingField).
		WithOutputFields(outputFields...))
}

// buildProductVectorOutputFields 返回检索结果需要回填的字段列表。
func buildProductVectorOutputFields() []string {
	return []string{
		productVectorPKField,
		productVectorNameField,
		productVectorBarCodeField,
		productVectorCategoryIDField,
		productVectorCategoryNameField,
		productVectorUnitIDField,
		productVectorUnitNameField,
		productVectorStandardField,
		productVectorRemarkField,
		productVectorExpiryDayField,
		productVectorWeightField,
		productVectorPurchasePriceField,
		productVectorSalePriceField,
		productVectorMinPriceField,
		productVectorContentField,
	}
}
