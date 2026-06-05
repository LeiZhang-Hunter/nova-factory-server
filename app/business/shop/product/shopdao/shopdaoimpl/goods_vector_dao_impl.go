package shopdaoimpl

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/datasource/milvus"
	searchutil "nova-factory-server/app/utils/vectorsearch"
)

const (
	defaultGoodsVectorCollection = "shop_goods_vectors"
	goodsTextBm25EmbFunctionName = "text_bm25_emb"

	goodsVectorPKField          = "goods_id"
	goodsVectorDBIDField        = "goods_db_id"
	goodsVectorNameField        = "goods_name"
	goodsVectorCodeField        = "goods_code"
	goodsVectorCategoryField    = "category_name"
	goodsVectorDescriptionField = "description"

	goodsVectorSkuIdField          = "sku_id"
	goodsVectorSkuNameField        = "sku_name"
	goodsVectorSkuDescriptionField = "sku_description"

	goodsVectorRetailPriceField   = "retail_price"
	goodsVectorWeightPriceField   = "weight"
	goodsVectorQuantityPriceField = "quantity"

	goodsVectorMetadataField      = "metadata"
	goodsVectorContentField       = "content"
	goodsVectorEmbeddingField     = "vector"
	goodsVectorContextSparseField = "text_sparse_vector"
	goodsIdxVectorField           = "idx_goods_vector"
	goodsIdxTextSparseVector      = "idx_goods_text_sparse_vector"

	goodsVectorPKMaxLength       = 128
	goodsVectorNameMaxLength     = 512
	goodsVectorCodeMaxLength     = 128
	goodsVectorCategoryMaxLength = 256
	goodsVectorSkuNameMaxLength  = 512
	goodsVectorDescMaxLength     = 4096
	goodsVectorContentMaxLength  = 16384

	goodsVectorSearchCandidateMultiplier = 3
	goodsVectorSearchMinCandidates       = 20
	goodsVectorSearchMaxCandidates       = 100
	goodsVectorSearchDefaultLimit        = 10
	goodsVectorSearchMaxLimit            = 50
)

type ShopGoodsVectorDaoImpl struct{}

type goodsVectorConfig struct {
	Collection string `mapstructure:"collection"`
}

type goodsVectorRows struct {
	pks             []int64
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
	vectors         [][]float32
}

func NewShopGoodsVectorDao() shopdao.IShopGoodsVectorDao {
	return &ShopGoodsVectorDaoImpl{}
}

func (d *ShopGoodsVectorDaoImpl) Upsert(c *gin.Context, goods *shopmodels.Goods, items []*shopmodels.GoodsVectorUpsertItem) (*shopmodels.GoodsVectorResult, error) {
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	if len(items) == 0 {
		return nil, errors.New("商品向量写入项为空")
	}

	dim := 0
	for idx, item := range items {
		if item == nil {
			return nil, fmt.Errorf("第%d条商品向量写入项为空", idx+1)
		}
		item.Content = trimRunes(item.Content, goodsVectorContentMaxLength)
		if strings.TrimSpace(item.Content) == "" {
			return nil, fmt.Errorf("第%d条商品内容为空", idx+1)
		}
		if len(item.Vector) == 0 {
			return nil, fmt.Errorf("第%d条商品向量为空", idx+1)
		}
		if dim == 0 {
			dim = len(item.Vector)
			continue
		}
		if len(item.Vector) != dim {
			return nil, fmt.Errorf("第%d条商品向量维度不一致，expected=%d actual=%d", idx+1, dim, len(item.Vector))
		}
	}

	cfg, err := loadGoodsVectorConfig()
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	client, err := milvus.GetClient(requestCtx)
	if err != nil {
		return nil, fmt.Errorf("初始化 Milvus 客户端失败: %w", err)
	}

	if err = ensureGoodsVectorCollection(requestCtx, client, cfg.Collection, dim); err != nil {
		return nil, err
	}

	rows := buildGoodsVectorRows(goods, items)
	_, err = client.Upsert(requestCtx, milvusclient.NewColumnBasedInsertOption(cfg.Collection).
		WithInt64Column(goodsVectorPKField, rows.pks).
		WithInt64Column(goodsVectorDBIDField, rows.goodsDBIDs).
		WithVarcharColumn(goodsVectorNameField, rows.goodsNames).
		WithVarcharColumn(goodsVectorCodeField, rows.goodsCodes).
		WithVarcharColumn(goodsVectorCategoryField, rows.categorys).
		WithVarcharColumn(goodsVectorDescriptionField, rows.descriptions).
		WithVarcharColumn(goodsVectorContentField, rows.contents).
		WithInt64Column(goodsVectorSkuIdField, rows.skuIDs).
		WithVarcharColumn(goodsVectorSkuNameField, rows.skuNames).
		WithVarcharColumn(goodsVectorSkuDescriptionField, rows.skuDescriptions).
		WithColumns(
			column.NewColumnDouble(goodsVectorRetailPriceField, rows.retailPrices),
			column.NewColumnDouble(goodsVectorWeightPriceField, rows.weights),
		).
		WithInt64Column(goodsVectorQuantityPriceField, rows.quantities).
		WithFloatVectorColumn(goodsVectorEmbeddingField, dim, rows.vectors),
	)
	if err != nil {
		return nil, fmt.Errorf("写入 Milvus 失败: %w", err)
	}

	goodsPK := strings.TrimSpace(goods.GoodsID)
	if goodsPK == "" {
		goodsPK = strconv.FormatInt(goods.ID, 10)
	}

	return &shopmodels.GoodsVectorResult{
		GoodsDBID:  goods.ID,
		GoodsID:    goodsPK,
		Collection: cfg.Collection,
		Dimension:  dim,
		SkuCount:   len(rows.pks),
	}, nil
}

func (d *ShopGoodsVectorDaoImpl) Search(c *gin.Context, req *shopmodels.GoodsVectorSearchReq, vector []float32) (*shopmodels.GoodsVectorSearchData, error) {
	if req == nil {
		return nil, errors.New("商品向量搜索参数为空")
	}
	if len(vector) == 0 {
		return nil, errors.New("商品搜索向量为空")
	}
	data, err := d.BatchSearch(c, &shopmodels.GoodsVectorBatchSearchReq{
		Queries:     []string{req.Query},
		SearchTexts: []string{req.SearchText},
		Limit:       req.Limit,
	}, [][]float32{vector})
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Rows) == 0 || data.Rows[0] == nil {
		return &shopmodels.GoodsVectorSearchData{
			Rows:  make([]*shopmodels.GoodsVectorSearchItem, 0),
			Total: 0,
		}, nil
	}
	return &shopmodels.GoodsVectorSearchData{
		Rows:  data.Rows[0].Rows,
		Total: data.Rows[0].Total,
	}, nil
}

func (d *ShopGoodsVectorDaoImpl) BatchSearch(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq,
	vectors [][]float32) (*shopmodels.GoodsVectorBatchSearchData, error) {
	if req == nil {
		return nil, errors.New("商品批量向量搜索参数为空")
	}
	if len(vectors) == 0 {
		return nil, errors.New("商品批量搜索向量为空")
	}
	if len(req.Queries) != len(vectors) {
		return nil, fmt.Errorf("商品批量搜索参数数量不匹配，expected=%d actual=%d", len(req.Queries), len(vectors))
	}

	cfg, err := loadGoodsVectorConfig()
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	client, err := milvus.GetClient(requestCtx)
	if err != nil {
		return nil, fmt.Errorf("初始化 Milvus 客户端失败: %w", err)
	}

	has, err := client.HasCollection(requestCtx, milvusclient.NewHasCollectionOption(cfg.Collection))
	if err != nil {
		return nil, fmt.Errorf("检查 Milvus collection 失败: %w", err)
	}
	if !has {
		return buildEmptyGoodsVectorBatchSearchData(req.Queries), nil
	}

	searchLimit := resolveGoodsVectorSearchCandidateLimit(req.Limit)
	searchVectors := make([]entity.Vector, 0, len(vectors))
	textQueries := make([]entity.Vector, 0, len(req.Queries))
	for idx, vector := range vectors {
		if len(vector) == 0 {
			return nil, fmt.Errorf("第%d条商品搜索向量为空", idx+1)
		}
		query := strings.TrimSpace(req.Queries[idx])
		if query == "" {
			return nil, fmt.Errorf("第%d条商品搜索文本为空", idx+1)
		}
		searchText := query
		if len(req.SearchTexts) == len(req.Queries) {
			if candidate := strings.TrimSpace(req.SearchTexts[idx]); candidate != "" {
				searchText = candidate
			}
		}
		searchVectors = append(searchVectors, entity.FloatVector(vector))
		textQueries = append(textQueries, entity.Text(searchText))
	}

	outputFields := []string{
		goodsVectorDBIDField,
		goodsVectorNameField,
		goodsVectorCodeField,
		goodsVectorCategoryField,
		goodsVectorDescriptionField,
		goodsVectorSkuIdField,
		goodsVectorSkuNameField,
		goodsVectorSkuDescriptionField,
		goodsVectorRetailPriceField,
		goodsVectorWeightPriceField,
		goodsVectorQuantityPriceField,
		goodsVectorContentField,
	}

	supportHybrid, err := supportsGoodsVectorHybridSearch(requestCtx, client, cfg.Collection)
	if err != nil {
		return nil, err
	}

	var resultSets []milvusclient.ResultSet
	if supportHybrid {
		sparseAnnParam := index.NewSparseAnnParam()
		sparseAnnParam.WithDropRatio(0.2)

		denseRequest := milvusclient.NewAnnRequest(goodsVectorEmbeddingField, searchLimit, searchVectors...)
		sparseRequest := milvusclient.NewAnnRequest(goodsVectorContextSparseField, searchLimit, textQueries...).
			WithAnnParam(sparseAnnParam)

		resultSets, err = client.HybridSearch(requestCtx, milvusclient.NewHybridSearchOption(
			cfg.Collection,
			searchLimit,
			denseRequest,
			sparseRequest,
		).WithReranker(milvusclient.NewRRFReranker()).
			WithOutputFields(outputFields...))
	} else {
		resultSets, err = client.Search(requestCtx, milvusclient.NewSearchOption(cfg.Collection, searchLimit,
			searchVectors).
			WithANNSField(goodsVectorEmbeddingField).
			WithOutputFields(outputFields...))
	}
	if err != nil {
		return nil, fmt.Errorf("搜索 Milvus 商品向量失败: %w", err)
	}
	if len(resultSets) == 0 {
		return buildEmptyGoodsVectorBatchSearchData(req.Queries), nil
	}
	if len(resultSets) != len(req.Queries) {
		return nil, fmt.Errorf("Milvus 批量搜索结果数量不匹配，expected=%d actual=%d", len(req.Queries), len(resultSets))
	}

	rows := make([]*shopmodels.GoodsVectorBatchSearchItem, 0, len(resultSets))
	for idx, resultSet := range resultSets {
		data, parseErr := parseGoodsVectorSearchResultSet(resultSet)
		if parseErr != nil {
			return nil, fmt.Errorf("解析第%d条商品向量搜索结果失败: %w", idx+1, parseErr)
		}
		data.Rows = rerankGoodsVectorSearchRows(req.Queries[idx], data.Rows, req.Limit)
		rows = append(rows, &shopmodels.GoodsVectorBatchSearchItem{
			Query: req.Queries[idx],
			Rows:  data.Rows,
			Total: int64(len(data.Rows)),
		})
	}

	return &shopmodels.GoodsVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

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

func buildEmptyGoodsVectorBatchSearchData(queries []string) *shopmodels.GoodsVectorBatchSearchData {
	rows := make([]*shopmodels.GoodsVectorBatchSearchItem, 0, len(queries))
	for _, query := range queries {
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

func parseGoodsVectorSearchResultSet(resultSet milvusclient.ResultSet) (*shopmodels.GoodsVectorSearchData, error) {
	if resultSet.Err != nil {
		return nil, resultSet.Err
	}

	rows := make([]*shopmodels.GoodsVectorSearchItem, 0, resultSet.ResultCount)
	for i := 0; i < resultSet.ResultCount; i++ {
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

func normalizeGoodsVectorSearchLimit(limit int) int {
	if limit <= 0 {
		return goodsVectorSearchDefaultLimit
	}
	if limit > goodsVectorSearchMaxLimit {
		return goodsVectorSearchMaxLimit
	}
	return limit
}

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

func ensureGoodsVectorCollection(ctx context.Context, client *milvusclient.Client, collectionName string, dim int) error {
	has, err := client.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if err != nil {
		return fmt.Errorf("检查 Milvus collection 失败: %w", err)
	}

	if !has {
		function := entity.NewFunction().
			WithName(goodsTextBm25EmbFunctionName).
			WithInputFields(goodsVectorContentField).
			WithOutputFields(goodsVectorContextSparseField).
			WithType(entity.FunctionTypeBM25)
		schema := entity.NewSchema().
			WithField(entity.NewField().WithName(goodsVectorPKField).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithMaxLength(goodsVectorPKMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorDBIDField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorCodeField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorCodeMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorCategoryField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorCategoryMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorDescriptionField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorDescMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorContentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorContentMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorSkuIdField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorSkuNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorSkuNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorSkuDescriptionField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorSkuNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorRetailPriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(goodsVectorWeightPriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(goodsVectorQuantityPriceField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorEmbeddingField).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(dim))).
			WithField(entity.NewField().WithName(goodsVectorContextSparseField).WithDataType(entity.FieldTypeSparseVector)).
			WithField(entity.NewField().
				WithName(goodsVectorMetadataField).
				WithDataType(entity.FieldTypeJSON).
				WithNullable(true),
			).
			WithFunction(function)

		err = client.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(collectionName, schema).WithIndexOptions(
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorEmbeddingField, index.NewAutoIndex(entity.COSINE)).WithIndexName(goodsIdxVectorField),
		).WithIndexOptions(
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorContextSparseField, index.NewSparseInvertedIndex(entity.BM25, 0.2)).WithIndexName(goodsIdxTextSparseVector),
		))
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
		goodsVectorPKField:             false,
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
	return nil
}

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

func buildGoodsVectorRows(goods *shopmodels.Goods, items []*shopmodels.GoodsVectorUpsertItem) *goodsVectorRows {
	goodsPK := strings.TrimSpace(goods.GoodsID)
	if goodsPK == "" {
		goodsPK = strconv.FormatInt(goods.ID, 10)
	}

	rows := &goodsVectorRows{}
	appendRow := func(item *shopmodels.GoodsVectorUpsertItem) {
		rows.pks = append(rows.pks, item.SkuID)
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
		rows.vectors = append(rows.vectors, item.Vector)
	}

	for _, item := range items {
		if item == nil {
			continue
		}
		appendRow(item)
	}

	return rows
}
