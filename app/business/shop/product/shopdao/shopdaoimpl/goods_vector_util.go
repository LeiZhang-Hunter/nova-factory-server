package shopdaoimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/product/shopmodels"
	searchutil "nova-factory-server/app/utils/vectorsearch"
	"nova-factory-server/app/utils/vectorsearch/goods"
	"nova-factory-server/app/utils/vectorsearch/normalization"
	"sort"
	"strconv"
	"strings"
)

const (
	// 默认 collection 名及 BM25 function 名称。
	defaultGoodsVectorCollection = "shop_goods_vectors"
	goodsTextBm25EmbFunctionName = "text_bm25_emb"

	// 商品主数据字段。
	goodsVectorGoodsIDField     = "goods_id"
	goodsVectorDBIDField        = "goods_db_id"
	goodsVectorNameField        = "goods_name"
	goodsVectorCodeField        = "goods_code"
	goodsVectorCategoryField    = "category_name"
	goodsVectorDescriptionField = "description"

	// SKU 粒度字段。当前向量库按 SKU 落行，便于把规格差异保留下来。
	goodsVectorSkuIdField          = "sku_id"
	goodsVectorSkuNameField        = "sku_name"
	goodsVectorSkuDescriptionField = "sku_description"

	// 价格、重量、库存等可直接回传给业务侧的数值字段。
	goodsVectorRetailPriceField   = "retail_price"
	goodsVectorWeightPriceField   = "weight"
	goodsVectorQuantityPriceField = "quantity"
	goodsVectorIsSaleField        = "is_sale"

	// 检索相关字段。
	// metadata 存放规格、分类等结构化信息，用于过滤；
	// content 是送入 BM25 与 embedding 的原始语义文本；
	// vector / text_sparse_vector 分别用于 dense 与 sparse 检索。
	goodsVectorMetadataField      = "metadata"
	goodsVectorContentField       = "content"
	goodsVectorEmbeddingField     = "vector"
	goodsVectorContextSparseField = "text_sparse_vector"

	// 主要索引名。
	goodsIdxVectorField      = "idx_goods_vector"
	goodsIdxTextSparseVector = "idx_goods_text_sparse_vector"
	goodsIdxCategoryField    = "idx_goods_category_name"
	goodsIdxIsSaleField      = "idx_goods_is_sale"

	// 字符串字段长度上限，避免写入 Milvus 时超出 schema 约束。
	goodsVectorPKMaxLength       = 128
	goodsVectorNameMaxLength     = 512
	goodsVectorCodeMaxLength     = 128
	goodsVectorCategoryMaxLength = 256
	goodsVectorSkuNameMaxLength  = 512
	goodsVectorDescMaxLength     = 4096
	goodsVectorContentMaxLength  = 16384

	// 检索候选集控制参数。
	// 先扩大召回，再在应用层重排，通常比直接取 limit 个候选更稳。
	goodsVectorSearchCandidateMultiplier = 3
	goodsVectorSearchMinCandidates       = 20
	goodsVectorSearchMaxCandidates       = 100
	goodsVectorSearchDefaultLimit        = 10
	goodsVectorSearchMaxLimit            = 50
)

// ShopGoodsVectorDaoImpl 负责商品向量数据的落库、检索以及 Milvus collection 自检。
// 当前还额外挂了一个 metadataExtractor，用于把商品文本中的规格、分类等结构化信息
// 提前提取到 metadata 字段，便于后续过滤检索。
type ShopGoodsVectorDaoImpl struct {
	metadataExtractor *goods.MetadataExtractor
}

// goodsVectorConfig 封装商品向量表的最小配置，目前只暴露 collection 名。
type goodsVectorConfig struct {
	Collection string `mapstructure:"collection"`
}

// goodsVectorRows 是 Milvus 列式写入所需的中间结构。
// DAO 先把商品与 SKU 业务对象展开成同长度切片，再一次性调用 Upsert。
type goodsVectorRows struct {
	goodsDBIDs      []int64
	goodsNames      []string
	goodsCodes      []string
	categorys       []string
	descriptions    []string
	contents        []string
	skuIDs          []int64
	skuNames        []string
	skuDescriptions []string
	retailPrices    []float64
	weights         []float64
	quantities      []int64
	metadatas       [][]byte
	vectors         [][]float32
	saleFlags       []bool
}

// goodsSearchRuntimeQuery 表示一次批量检索中单条 query 的运行时上下文。
// 除了原始 query 外，还会缓存：
// 1. 向量检索用的 dense vector
// 2. hybrid 检索用的 text
// 3. metadata 过滤表达式
// 4. 原始顺序下标，便于分组搜索后再回填结果
type goodsSearchRuntimeQuery struct {
	index      int
	query      string
	vector     entity.Vector
	text       entity.Vector
	filterExpr string
}

// supportsGoodsVectorHybridSearch 检查当前 collection 是否具备稀疏向量字段。
func supportsGoodsVectorHybridSearch(ctx context.Context, client *milvusclient.Client, collectionName string) (bool, error) {
	collection, err := client.DescribeCollection(ctx, milvusclient.NewDescribeCollectionOption(collectionName))
	if err != nil {
		return false, fmt.Errorf("读取 Milvus collection 结构失败: %w", err)
	}
	if collection == nil || collection.Schema == nil {
		return false, fmt.Errorf("Milvus collection %s 缺少 schema", collectionName)
	}

	for _, field := range collection.Schema.Fields {
		if field == nil {
			continue
		}
		if field.Name == goodsVectorContextSparseField {
			return true, nil
		}
	}
	return false, nil
}

// buildEmptyGoodsVectorBatchSearchData 为每条 query 构造空结果，保持返回结构稳定。
func buildEmptyGoodsVectorBatchSearchData(queries []string) *shopmodels.GoodsVectorBatchSearchData {
	rows := make([]*shopmodels.GoodsVectorBatchSearchItem, 0, len(queries))
	for _, query := range queries {
		// 即使完全无结果，也保留每条 query 的壳结构，便于上层按输入顺序消费。
		rows = append(rows, &shopmodels.GoodsVectorBatchSearchItem{
			Query: query,
			Rows:  make([]*shopmodels.GoodsVectorSearchItem, 0),
			Total: 0,
		})
	}
	return &shopmodels.GoodsVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}
}

// parseGoodsVectorSearchResultSet 将 Milvus 返回结果映射为业务搜索结果结构。
func parseGoodsVectorSearchResultSet(resultSet milvusclient.ResultSet) (*shopmodels.GoodsVectorSearchData, error) {
	if resultSet.Err != nil {
		return nil, resultSet.Err
	}

	rows := make([]*shopmodels.GoodsVectorSearchItem, 0, resultSet.ResultCount)
	for i := 0; i < resultSet.ResultCount; i++ {
		// Milvus 返回的是按列组织的数据，这里逐列回填成业务层更容易使用的 item 结构。
		item := &shopmodels.GoodsVectorSearchItem{
			Score: resultSet.Scores[i],
		}
		var err error
		if item.GoodsDBID, err = getInt64ColumnValue(&resultSet, goodsVectorDBIDField, i); err != nil {
			return nil, err
		}
		if item.GoodsName, err = getStringColumnValue(&resultSet, goodsVectorNameField, i); err != nil {
			return nil, err
		}
		if item.GoodsCode, err = getStringColumnValue(&resultSet, goodsVectorCodeField, i); err != nil {
			return nil, err
		}
		if item.CategoryName, err = getStringColumnValue(&resultSet, goodsVectorCategoryField, i); err != nil {
			return nil, err
		}
		if item.Description, err = getStringColumnValue(&resultSet, goodsVectorDescriptionField, i); err != nil {
			return nil, err
		}
		if item.SkuID, err = getInt64ColumnValue(&resultSet, goodsVectorSkuIdField, i); err != nil {
			return nil, err
		}
		if item.SkuName, err = getStringColumnValue(&resultSet, goodsVectorSkuNameField, i); err != nil {
			return nil, err
		}
		if item.SkuDescription, err = getStringColumnValue(&resultSet, goodsVectorSkuDescriptionField, i); err != nil {
			return nil, err
		}
		if item.RetailPrice, err = getFloat64ColumnValue(&resultSet, goodsVectorRetailPriceField, i); err != nil {
			return nil, err
		}
		if item.Weight, err = getFloat64ColumnValue(&resultSet, goodsVectorWeightPriceField, i); err != nil {
			return nil, err
		}
		if item.Quantity, err = getInt64ColumnValue(&resultSet, goodsVectorQuantityPriceField, i); err != nil {
			return nil, err
		}
		if item.Content, err = getStringColumnValue(&resultSet, goodsVectorContentField, i); err != nil {
			return nil, err
		}
		rows = append(rows, item)
	}

	return &shopmodels.GoodsVectorSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

// normalizeGoodsVectorSearchLimit 将调用方 limit 归一化到 DAO 允许的安全范围内。
func normalizeGoodsVectorSearchLimit(limit int) int {
	if limit <= 0 {
		return goodsVectorSearchDefaultLimit
	}
	if limit > goodsVectorSearchMaxLimit {
		return goodsVectorSearchMaxLimit
	}
	return limit
}

// resolveGoodsVectorSearchCandidateLimit 计算 Milvus 初召候选数。
// 这里不会直接使用最终 limit，而是适当放大候选集，留给应用层业务重排更多空间。
func resolveGoodsVectorSearchCandidateLimit(limit int) int {
	limit = normalizeGoodsVectorSearchLimit(limit)
	candidateLimit := limit * goodsVectorSearchCandidateMultiplier
	if candidateLimit < goodsVectorSearchMinCandidates {
		candidateLimit = goodsVectorSearchMinCandidates
	}
	if candidateLimit > goodsVectorSearchMaxCandidates {
		candidateLimit = goodsVectorSearchMaxCandidates
	}
	if candidateLimit < limit {
		return limit
	}
	return candidateLimit
}

// rerankGoodsVectorSearchRows 将 Milvus 初召结果映射为通用候选，再结合 query 做业务重排。
func rerankGoodsVectorSearchRows(query string, rows []*shopmodels.GoodsVectorSearchItem, limit int) []*shopmodels.GoodsVectorSearchItem {
	if len(rows) == 0 {
		return make([]*shopmodels.GoodsVectorSearchItem, 0)
	}

	processedQuery := searchutil.ProcessQuery(query)
	candidates := make([]searchutil.RankCandidate, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		candidates = append(candidates, searchutil.RankCandidate{
			ID:        row.GoodsDBID,
			Title:     row.GoodsName,
			Code:      row.GoodsCode,
			Category:  row.CategoryName,
			Unit:      "",
			Standard:  row.SkuName,
			Remark:    row.SkuDescription,
			Content:   strings.TrimSpace(row.Content + " " + row.Description),
			BaseScore: row.Score,
		})
	}
	if len(candidates) == 0 {
		return make([]*shopmodels.GoodsVectorSearchItem, 0)
	}

	ranked := searchutil.RerankCandidates(processedQuery, candidates, normalizeGoodsVectorSearchLimit(limit))
	result := make([]*shopmodels.GoodsVectorSearchItem, 0, len(ranked))
	for _, item := range ranked {
		if item.Index < 0 || item.Index >= len(rows) || rows[item.Index] == nil {
			continue
		}
		row := rows[item.Index]
		row.Score = item.Score
		result = append(result, row)
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			if result[i].GoodsDBID == result[j].GoodsDBID {
				return result[i].SkuID < result[j].SkuID
			}
			return result[i].GoodsDBID < result[j].GoodsDBID
		}
		return result[i].Score > result[j].Score
	})
	return result
}

// loadGoodsVectorConfig 读取商品向量相关配置，并回退到默认 collection。
func loadGoodsVectorConfig() (*goodsVectorConfig, error) {
	cfg := &goodsVectorConfig{}
	_ = viper.UnmarshalKey("shop.goods_vector", cfg)

	if cfg.Collection == "" {
		cfg.Collection = strings.TrimSpace(viper.GetString("milvus.goods_collection"))
	}
	if cfg.Collection == "" {
		cfg.Collection = defaultGoodsVectorCollection
	}

	return cfg, nil
}

// ensureGoodsVectorCollection 确保 Milvus collection 存在且字段与维度符合预期。
func ensureGoodsVectorCollection(ctx context.Context, client *milvusclient.Client, collectionName string, dim int) error {
	has, err := client.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if err != nil {
		return fmt.Errorf("检查 Milvus collection 失败: %w", err)
	}

	if !has {
		// 新建 collection 时同时挂上 BM25 function，让 content 自动生成稀疏向量。
		function := entity.NewFunction().
			WithName(goodsTextBm25EmbFunctionName).
			WithInputFields(goodsVectorContentField).
			WithOutputFields(goodsVectorContextSparseField).
			WithType(entity.FunctionTypeBM25)
		schema := entity.NewSchema().
			WithField(entity.NewField().WithName(goodsVectorSkuIdField).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithMaxLength(goodsVectorPKMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorGoodsIDField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorDBIDField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorCodeField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorCodeMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorCategoryField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorCategoryMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorDescriptionField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorDescMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorContentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorContentMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorSkuNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorSkuNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorSkuDescriptionField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorSkuNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorRetailPriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(goodsVectorWeightPriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(goodsVectorQuantityPriceField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorIsSaleField).WithDataType(entity.FieldTypeBool)).
			WithField(entity.NewField().WithName(goodsVectorEmbeddingField).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(dim))).
			WithField(entity.NewField().WithName(goodsVectorContextSparseField).WithDataType(entity.FieldTypeSparseVector)).
			WithField(entity.NewField().
				WithName(goodsVectorMetadataField).
				WithDataType(entity.FieldTypeJSON).
				WithNullable(true),
			).
			WithFunction(function)

		// 索引策略：
		// 1. dense 向量使用 AutoIndex(COSINE)
		// 2. 分类名使用倒排索引，便于普通标量过滤
		// 3. 稀疏向量使用 BM25 SparseInvertedIndex，供 hybrid search 融合召回
		// metadata 当前是动态 JSON，仅存储过滤字段，不预建索引。
		indexOptions := []milvusclient.CreateIndexOption{
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorEmbeddingField, index.NewAutoIndex(entity.COSINE)).WithIndexName(goodsIdxVectorField),
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorCategoryField, index.NewInvertedIndex()).WithIndexName(goodsIdxCategoryField),
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorIsSaleField, index.NewInvertedIndex()).WithIndexName(goodsIdxIsSaleField),
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorContextSparseField, index.NewSparseInvertedIndex(entity.BM25, 0.2)).WithIndexName(goodsIdxTextSparseVector),
		}

		err = client.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(collectionName, schema).WithIndexOptions(indexOptions...))
		if err != nil {
			return fmt.Errorf("创建 Milvus collection 失败: %w", err)
		}
		return nil
	}

	collection, err := client.DescribeCollection(ctx, milvusclient.NewDescribeCollectionOption(collectionName))
	if err != nil {
		return fmt.Errorf("读取 Milvus collection 结构失败: %w", err)
	}
	if collection == nil || collection.Schema == nil {
		return fmt.Errorf("Milvus collection %s 缺少 schema", collectionName)
	}

	requiredFields := map[string]bool{
		// 这里只校验写入与检索强依赖字段，避免历史 schema 缺列却继续运行。
		goodsVectorGoodsIDField:        false,
		goodsVectorDBIDField:           false,
		goodsVectorNameField:           false,
		goodsVectorCodeField:           false,
		goodsVectorCategoryField:       false,
		goodsVectorDescriptionField:    false,
		goodsVectorContentField:        false,
		goodsVectorSkuIdField:          false,
		goodsVectorSkuNameField:        false,
		goodsVectorSkuDescriptionField: false,
		goodsVectorRetailPriceField:    false,
		goodsVectorWeightPriceField:    false,
		goodsVectorQuantityPriceField:  false,
		goodsVectorIsSaleField:         false,
		goodsVectorEmbeddingField:      false,
	}
	for _, field := range collection.Schema.Fields {
		if field == nil {
			continue
		}
		if _, ok := requiredFields[field.Name]; ok {
			requiredFields[field.Name] = true
		}
		if field.Name != goodsVectorEmbeddingField {
			continue
		}
		// 已存在 collection 时重点校验向量维度，避免不同 embedding 模型写到同一张表里。
		existingDim, dimErr := field.GetDim()
		if dimErr != nil {
			return fmt.Errorf("读取 Milvus 向量维度失败: %w", dimErr)
		}
		if int(existingDim) != dim {
			return fmt.Errorf("Milvus collection %s 向量维度不匹配，现有=%d，请求=%d", collectionName, existingDim, dim)
		}
	}
	for fieldName, exists := range requiredFields {
		if !exists {
			return fmt.Errorf("Milvus collection %s 缺少 %s 字段", collectionName, fieldName)
		}
	}
	// 历史 collection 不强改 schema，但会补齐必要的标量索引。
	if err = ensureGoodsVectorScalarIndex(ctx, client, collectionName, goodsVectorCategoryField, goodsIdxCategoryField, index.NewInvertedIndex()); err != nil {
		return err
	}
	if err = ensureGoodsVectorScalarIndex(ctx, client, collectionName, goodsVectorIsSaleField, goodsIdxIsSaleField, index.NewInvertedIndex()); err != nil {
		return err
	}
	return nil
}

// ensureGoodsVectorScalarIndex 确保指定字段存在至少一个索引；若没有则补建。
func ensureGoodsVectorScalarIndex(ctx context.Context, client *milvusclient.Client, collectionName, fieldName, indexName string, idx index.Index) error {
	indexNames, err := client.ListIndexes(ctx, milvusclient.NewListIndexOption(collectionName).WithFieldName(fieldName))
	if err != nil {
		return fmt.Errorf("读取 Milvus 字段[%s]索引失败: %w", fieldName, err)
	}
	if len(indexNames) > 0 {
		return nil
	}

	task, err := client.CreateIndex(ctx, milvusclient.NewCreateIndexOption(collectionName, fieldName, idx).WithIndexName(indexName))
	if err != nil {
		return fmt.Errorf("创建 Milvus 字段[%s]索引失败: %w", fieldName, err)
	}
	if task != nil {
		if err = task.Await(ctx); err != nil {
			return fmt.Errorf("等待 Milvus 字段[%s]索引完成失败: %w", fieldName, err)
		}
	}
	return nil
}

// buildRequestContext 优先复用请求上下文，以便透传超时与取消信号。
func buildRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

// trimRunes 将字符串截断到指定最大字符数（按 Unicode 码点计算），超过则截断并返回
func trimRunes(value string, max int) string {
	value = strings.TrimSpace(value)
	if max <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}

// getInt64ColumnValue 从 Milvus 结果集中读取 int64 字段值。
func getInt64ColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (int64, error) {
	if resultSet == nil {
		return 0, fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return 0, fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsInt64(idx)
}

// getStringColumnValue 从 Milvus 结果集中读取字符串字段值。
func getStringColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (string, error) {
	if resultSet == nil {
		return "", fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return "", fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsString(idx)
}

// getFloat64ColumnValue 从 Milvus 结果集中读取浮点字段值。
func getFloat64ColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (float64, error) {
	if resultSet == nil {
		return 0, fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return 0, fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsDouble(idx)
}

// buildGoodsVectorRows 将商品及 SKU 数据展开为 Milvus 列式写入所需结构。
func buildGoodsVectorRows(goods *shopmodels.Goods, items []*shopmodels.GoodsVectorUpsertItem, extractor *goods.MetadataExtractor) (*goodsVectorRows, error) {
	goodsPK := strings.TrimSpace(goods.GoodsID)
	if goodsPK == "" {
		goodsPK = strconv.FormatInt(goods.ID, 10)
	}

	rows := &goodsVectorRows{}
	appendRow := func(item *shopmodels.GoodsVectorUpsertItem) error {
		// 每一行除了基础字段和向量本身，还会生成一份 metadata JSON，
		// 供后续按规格、分类等结构化条件过滤。
		metadataJSON, err := buildGoodsVectorMetadataJSON(goods, item, extractor)
		if err != nil {
			return err
		}
		rows.saleFlags = append(rows.saleFlags, item.IsSale)
		rows.goodsDBIDs = append(rows.goodsDBIDs, goods.ID)
		rows.goodsNames = append(rows.goodsNames, trimRunes(goods.GoodsName, goodsVectorNameMaxLength))
		rows.goodsCodes = append(rows.goodsCodes, trimRunes(goods.GoodsCode, goodsVectorCodeMaxLength))
		rows.categorys = append(rows.categorys, trimRunes(goods.ShopCategoryName, goodsVectorCategoryMaxLength))
		rows.descriptions = append(rows.descriptions, trimRunes(goods.Description, goodsVectorDescMaxLength))
		rows.contents = append(rows.contents, item.Content)
		rows.skuIDs = append(rows.skuIDs, item.SkuID)
		rows.skuDescriptions = append(rows.skuDescriptions, item.SkuDescription)
		rows.skuNames = append(rows.skuNames, trimRunes(item.SkuName, goodsVectorSkuNameMaxLength))
		rows.retailPrices = append(rows.retailPrices, item.RetailPrice)
		rows.weights = append(rows.weights, item.Weight)
		rows.quantities = append(rows.quantities, item.Quantity)
		rows.metadatas = append(rows.metadatas, metadataJSON)
		rows.vectors = append(rows.vectors, item.Vector)
		return nil
	}

	for _, item := range items {
		if item == nil {
			continue
		}
		if err := appendRow(item); err != nil {
			return nil, err
		}
	}

	return rows, nil
}

// buildGoodsVectorMetadataJSON 从商品与文本提取结果中组装 metadata JSON。
// 当前策略是：
// 1. 优先带上商品已有分类
// 2. 合并调用方显式传入的 metadata，便于更新向量时同步覆盖 metadata
// 3. 再从内容中提取规格、分类等 metadata，补齐缺失字段
// 4. 对同名 metadata 只保留首个归一化值，保证过滤表达式简单稳定
func buildGoodsVectorMetadataJSON(goodsModel *shopmodels.Goods, item *shopmodels.GoodsVectorUpsertItem, extractor *goods.MetadataExtractor) ([]byte, error) {
	metadata := make(map[string]any)
	if goodsModel != nil {
		if category := strings.TrimSpace(goodsModel.ShopCategoryName); category != "" {
			metadata["category"] = category
		}
	}
	if item != nil {
		mergeGoodsVectorExplicitMetadataMap(metadata, item.Metadata)
	}
	if extractor != nil && item != nil {
		result, err := extractor.Extract(item.Content)
		if err != nil {
			zap.L().Warn("extract goods vector metadata fail",
				zap.String("content", item.Content),
				zap.Error(err))
		} else {
			mergeGoodsVectorMetadataMap(metadata, result.Metadata)
			if _, ok := metadata["category"]; !ok && len(result.Categories) > 0 && strings.TrimSpace(result.Categories[0].Name) != "" {
				metadata["category"] = strings.TrimSpace(result.Categories[0].Name)
			}
		}
	}
	if len(metadata) == 0 {
		return []byte(`{}`), nil
	}
	body, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("序列化商品向量 metadata 失败: %w", err)
	}
	return body, nil
}

// mergeGoodsVectorExplicitMetadataMap 合并调用方显式提供的 metadata。
// 显式 metadata 代表业务侧更高优先级的结构化结果，因此会覆盖已有同名字段。
func mergeGoodsVectorExplicitMetadataMap(target map[string]any, values map[string]any) {
	if len(values) == 0 || target == nil {
		return
	}
	for key, value := range values {
		key = searchutil.NormalizeWhitespace(key)
		if key == "" {
			continue
		}
		switch typed := value.(type) {
		case string:
			typed = searchutil.NormalizeWhitespace(typed)
			if typed == "" {
				continue
			}
			target[key] = typed
		case []string:
			for _, item := range typed {
				item = searchutil.NormalizeWhitespace(item)
				if item == "" {
					continue
				}
				target[key] = item
				break
			}
		case []any:
			for _, item := range typed {
				text := searchutil.NormalizeWhitespace(fmt.Sprint(item))
				if text == "" {
					continue
				}
				target[key] = text
				break
			}
		default:
			text := searchutil.NormalizeWhitespace(fmt.Sprint(value))
			if text == "" {
				continue
			}
			target[key] = text
		}
	}
}

// mergeGoodsVectorMetadataMap 将提取器输出的 metadata 合并到写入 JSON 中。
// 目前 metadata 结构有意保持扁平，只取首个有效值，避免后续过滤表达式过于复杂。
// 若同名 key 已由显式 metadata 设置，则不再覆盖。
func mergeGoodsVectorMetadataMap(target map[string]any, values map[string][]string) {
	if len(values) == 0 || target == nil {
		return
	}
	for key, items := range values {
		key = searchutil.NormalizeWhitespace(key)
		if key == "" || len(items) == 0 {
			continue
		}
		normalized := make([]string, 0, len(items))
		for _, item := range items {
			item = searchutil.NormalizeWhitespace(item)
			if item == "" {
				continue
			}
			normalized = append(normalized, item)
		}
		if len(normalized) == 0 {
			continue
		}
		if _, exists := target[key]; exists {
			continue
		}
		target[key] = normalized[0]
	}
}

// buildGoodsMetadataFilterExpr 将提取结果转换为 Milvus filter 表达式。
// 一个 key 下如果提取出多个候选值，会被拼成 "(a or b)"；
// 多个 key 之间则使用 "and"，表示同时满足。
func buildGoodsMetadataFilterExpr(result normalization.Result) string {
	if len(result.Metadata) == 0 {
		return ""
	}
	keys := make([]string, 0, len(result.Metadata))
	for key, values := range result.Metadata {
		if strings.TrimSpace(key) == "" || len(values) == 0 {
			continue
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)
	exprs := make([]string, 0, len(keys))
	for _, key := range keys {
		valueExprs := make([]string, 0, len(result.Metadata[key]))
		for _, value := range result.Metadata[key] {
			value = searchutil.NormalizeWhitespace(value)
			if value == "" {
				continue
			}
			valueExprs = append(valueExprs, fmt.Sprintf(`%s[%s] == %s`,
				goodsVectorMetadataField,
				strconv.Quote(key),
				strconv.Quote(value),
			))
		}
		if len(valueExprs) == 0 {
			continue
		}
		if len(valueExprs) == 1 {
			exprs = append(exprs, valueExprs[0])
			continue
		}
		exprs = append(exprs, "("+strings.Join(valueExprs, " or ")+")")
	}
	return strings.Join(exprs, " and ")
}

// groupGoodsSearchRuntimeQueries 按 filterExpr 将查询分组。
// 这样同一过滤条件下的多条向量可以共用一次 Search/HybridSearch 请求。
func groupGoodsSearchRuntimeQueries(runtimeQueries []goodsSearchRuntimeQuery) []goodsSearchRuntimeGroup {
	groupIndexes := make(map[string]int, len(runtimeQueries))
	groups := make([]goodsSearchRuntimeGroup, 0, len(runtimeQueries))
	for _, runtimeQuery := range runtimeQueries {
		groupIdx, ok := groupIndexes[runtimeQuery.filterExpr]
		if !ok {
			groupIdx = len(groups)
			groupIndexes[runtimeQuery.filterExpr] = groupIdx
			groups = append(groups, goodsSearchRuntimeGroup{
				filterExpr: runtimeQuery.filterExpr,
				queries:    make([]goodsSearchRuntimeQuery, 0, 1),
			})
		}
		groups[groupIdx].queries = append(groups[groupIdx].queries, runtimeQuery)
	}
	return groups
}

// collectGoodsMetadataFallbackQueries 收集首轮 metadata 过滤未命中的 query，
// 并将其改写为无 filter 的回退查询。
func collectGoodsMetadataFallbackQueries(runtimeQueries []goodsSearchRuntimeQuery, resultSets []milvusclient.ResultSet) []goodsSearchRuntimeQuery {
	if len(runtimeQueries) == 0 || len(resultSets) == 0 {
		return nil
	}
	fallbackQueries := make([]goodsSearchRuntimeQuery, 0)
	for _, runtimeQuery := range runtimeQueries {
		if runtimeQuery.filterExpr == "" {
			continue
		}
		if runtimeQuery.index < 0 || runtimeQuery.index >= len(resultSets) {
			continue
		}
		resultSet := resultSets[runtimeQuery.index]
		if resultSet.Err != nil || resultSet.ResultCount > 0 {
			continue
		}
		runtimeQuery.filterExpr = ""
		fallbackQueries = append(fallbackQueries, runtimeQuery)
	}
	return fallbackQueries
}

// executeGoodsHybridSearchGroups 对每个分组执行一次 hybrid search。
// denseRequest 负责语义向量召回，sparseRequest 负责 BM25 文本召回，
// 最终由 RRF 重排器做融合。
func executeGoodsHybridSearchGroups(ctx context.Context, client *milvusclient.Client, collectionName string, searchLimit int,
	outputFields []string, groups []goodsSearchRuntimeGroup, resultSets []milvusclient.ResultSet) error {
	for _, group := range groups {
		groupVectors := make([]entity.Vector, 0, len(group.queries))
		groupTexts := make([]entity.Vector, 0, len(group.queries))
		for _, runtimeQuery := range group.queries {
			groupVectors = append(groupVectors, runtimeQuery.vector)
			groupTexts = append(groupTexts, runtimeQuery.text)
		}
		sparseAnnParam := index.NewSparseAnnParam()
		sparseAnnParam.WithDropRatio(0.2)

		denseRequest := milvusclient.NewAnnRequest(goodsVectorEmbeddingField, searchLimit, groupVectors...)
		sparseRequest := milvusclient.NewAnnRequest(goodsVectorContextSparseField, searchLimit, groupTexts...).
			WithAnnParam(sparseAnnParam)
		if group.filterExpr != "" {
			// 同一组内 filterExpr 相同，可以安全挂到两个召回请求上。
			denseRequest = denseRequest.WithFilter(group.filterExpr)
			sparseRequest = sparseRequest.WithFilter(group.filterExpr)
		}

		queryResultSets, err := client.HybridSearch(ctx, milvusclient.NewHybridSearchOption(
			collectionName,
			searchLimit,
			denseRequest,
			sparseRequest,
		).WithReranker(milvusclient.NewRRFReranker()).
			WithOutputFields(outputFields...))
		if err != nil {
			return fmt.Errorf("搜索 Milvus 商品向量失败: %w", err)
		}
		if err = assignGoodsSearchGroupResults(group.queries, queryResultSets, resultSets, "Milvus HybridSearch"); err != nil {
			return err
		}
	}
	return nil
}

// executeGoodsDenseSearchGroups 在不支持 hybrid 的 collection 上退化为纯 dense 搜索。
func executeGoodsDenseSearchGroups(ctx context.Context, client *milvusclient.Client, collectionName string, searchLimit int,
	outputFields []string, groups []goodsSearchRuntimeGroup, resultSets []milvusclient.ResultSet) error {
	for _, group := range groups {
		groupVectors := make([]entity.Vector, 0, len(group.queries))
		for _, runtimeQuery := range group.queries {
			groupVectors = append(groupVectors, runtimeQuery.vector)
		}
		searchOption := milvusclient.NewSearchOption(collectionName, searchLimit, groupVectors).
			WithANNSField(goodsVectorEmbeddingField).
			WithOutputFields(outputFields...)
		if group.filterExpr != "" {
			searchOption = searchOption.WithFilter(group.filterExpr)
		}
		queryResultSets, err := client.Search(ctx, searchOption)
		if err != nil {
			return fmt.Errorf("搜索 Milvus 商品向量失败: %w", err)
		}
		if err = assignGoodsSearchGroupResults(group.queries, queryResultSets, resultSets, "Milvus Search"); err != nil {
			return err
		}
	}
	return nil
}

// assignGoodsSearchGroupResults 将分组查询结果按原始 query 顺序回填。
// 分组执行后返回顺序只在组内有意义，因此必须借助 runtimeQuery.index 恢复到调用方顺序。
func assignGoodsSearchGroupResults(groupQueries []goodsSearchRuntimeQuery, queryResultSets []milvusclient.ResultSet,
	resultSets []milvusclient.ResultSet, label string) error {
	if len(queryResultSets) == 0 {
		for _, runtimeQuery := range groupQueries {
			resultSets[runtimeQuery.index] = milvusclient.ResultSet{}
		}
		return nil
	}
	if len(queryResultSets) != len(groupQueries) {
		return fmt.Errorf("%s 结果数量不匹配，expected=%d actual=%d", label, len(groupQueries), len(queryResultSets))
	}
	for idx, runtimeQuery := range groupQueries {
		resultSets[runtimeQuery.index] = queryResultSets[idx]
	}
	return nil
}
