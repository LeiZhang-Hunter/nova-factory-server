package masterdaoimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/datasource/milvus"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"
)

const (
	defaultProductVectorCollection = "erp_master_product_vectors"

	productTextBm25EmbFunctionName = "text_bm25_emb"

	productVectorPKField            = "product_id"
	productVectorNameField          = "name"
	productVectorBarCodeField       = "bar_code"
	productVectorCategoryIDField    = "category_id"
	productVectorCategoryNameField  = "category_name"
	productVectorUnitIDField        = "unit_id"
	productVectorUnitNameField      = "unit_name"
	productVectorStandardField      = "standard"
	productVectorRemarkField        = "remark"
	productVectorExpiryDayField     = "expiry_day"
	productVectorWeightField        = "weight"
	productVectorPurchasePriceField = "purchase_price"
	productVectorSalePriceField     = "sale_price"
	productVectorMinPriceField      = "min_price"
	productVectorContentField       = "content"
	productVectorEmbeddingField     = "vector"
	productVectorContextSparseField = "text_sparse_vector"

	productIdxProductVector           = "idx_product_vector"
	productIdxProductTextSparseVector = "idx_product_text_sparse_vector"

	productVectorNameMaxLength     = 512
	productVectorBarCodeMaxLength  = 128
	productVectorCategoryMaxLength = 256
	productVectorUnitMaxLength     = 128
	productVectorStandardMaxLength = 512
	productVectorRemarkMaxLength   = 2048
	productVectorContentMaxLength  = 16384
)

type ProductVectorDaoImpl struct{}

type productVectorConfig struct {
	Collection string `mapstructure:"collection"`
}

func NewProductVectorDao() masterdao.IProductVectorDao {
	return &ProductVectorDaoImpl{}
}

func (d *ProductVectorDaoImpl) Upsert(c *gin.Context, product *mastermodels.Product, item *mastermodels.ProductVectorUpsertItem) (*mastermodels.ProductVectorResult, error) {
	if product == nil || product.ID <= 0 {
		return nil, errors.New("产品不存在")
	}
	if item == nil {
		return nil, errors.New("产品向量写入项为空")
	}
	item.Content = trimProductVectorRunes(item.Content, productVectorContentMaxLength)
	if strings.TrimSpace(item.Content) == "" {
		return nil, errors.New("产品向量内容不能为空")
	}
	if len(item.Vector) == 0 {
		return nil, errors.New("产品向量不能为空")
	}

	cfg, err := loadProductVectorConfig()
	if err != nil {
		return nil, err
	}

	requestCtx := buildProductVectorRequestContext(c)
	client, err := milvus.GetClient(requestCtx)
	if err != nil {
		return nil, fmt.Errorf("初始化 Milvus 客户端失败: %w", err)
	}

	dim := len(item.Vector)
	if err = ensureProductVectorCollection(requestCtx, client, cfg.Collection, dim); err != nil {
		return nil, err
	}

	_, err = client.Upsert(requestCtx, milvusclient.NewColumnBasedInsertOption(cfg.Collection).
		WithInt64Column(productVectorPKField, []int64{product.ID}).
		WithVarcharColumn(productVectorNameField, []string{trimProductVectorRunes(product.Name, productVectorNameMaxLength)}).
		WithVarcharColumn(productVectorBarCodeField, []string{trimProductVectorRunes(product.BarCode, productVectorBarCodeMaxLength)}).
		WithInt64Column(productVectorCategoryIDField, []int64{product.CategoryId}).
		WithVarcharColumn(productVectorCategoryNameField, []string{trimProductVectorRunes(product.CategoryName, productVectorCategoryMaxLength)}).
		WithInt64Column(productVectorUnitIDField, []int64{product.UnitId}).
		WithVarcharColumn(productVectorUnitNameField, []string{trimProductVectorRunes(product.UnitName, productVectorUnitMaxLength)}).
		WithVarcharColumn(productVectorStandardField, []string{trimProductVectorRunes(product.Standard, productVectorStandardMaxLength)}).
		WithVarcharColumn(productVectorRemarkField, []string{trimProductVectorRunes(product.Remark, productVectorRemarkMaxLength)}).
		WithInt64Column(productVectorExpiryDayField, []int64{int64(product.ExpiryDay)}).
		WithColumns(
			column.NewColumnDouble(productVectorWeightField, []float64{product.Weight}),
			column.NewColumnDouble(productVectorPurchasePriceField, []float64{product.PurchasePrice}),
			column.NewColumnDouble(productVectorSalePriceField, []float64{product.SalePrice}),
			column.NewColumnDouble(productVectorMinPriceField, []float64{product.MinPrice}),
		).
		WithVarcharColumn(productVectorContentField, []string{item.Content}).
		WithFloatVectorColumn(productVectorEmbeddingField, dim, [][]float32{item.Vector}),
	)
	if err != nil {
		return nil, fmt.Errorf("写入 Milvus 失败: %w", err)
	}

	return &mastermodels.ProductVectorResult{
		ProductID:  product.ID,
		Collection: cfg.Collection,
		Dimension:  dim,
	}, nil
}

func (d *ProductVectorDaoImpl) Search(c *gin.Context, req *mastermodels.ProductVectorSearchReq, vector []float32) (*mastermodels.ProductVectorSearchData, error) {
	if req == nil {
		return nil, errors.New("产品向量搜索参数为空")
	}
	if len(vector) == 0 {
		return nil, errors.New("产品搜索向量为空")
	}
	data, err := d.BatchSearch(c, &mastermodels.ProductVectorBatchSearchReq{
		Queries: []string{req.Query},
		Limit:   req.Limit,
	}, [][]float32{vector})
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Rows) == 0 || data.Rows[0] == nil {
		return &mastermodels.ProductVectorSearchData{
			Rows:  make([]*mastermodels.ProductVectorSearchItem, 0),
			Total: 0,
		}, nil
	}
	return &mastermodels.ProductVectorSearchData{
		Rows:  data.Rows[0].Rows,
		Total: data.Rows[0].Total,
	}, nil
}

func (d *ProductVectorDaoImpl) BatchSearch(c *gin.Context, req *mastermodels.ProductVectorBatchSearchReq, vectors [][]float32) (*mastermodels.ProductVectorBatchSearchData, error) {
	if req == nil {
		return nil, errors.New("产品批量向量搜索参数为空")
	}
	if len(vectors) == 0 {
		return nil, errors.New("产品批量搜索向量为空")
	}
	if len(req.Queries) != len(vectors) {
		return nil, fmt.Errorf("产品批量搜索参数数量不匹配，expected=%d actual=%d", len(req.Queries), len(vectors))
	}

	cfg, err := loadProductVectorConfig()
	if err != nil {
		return nil, err
	}

	requestCtx := buildProductVectorRequestContext(c)
	client, err := milvus.GetClient(requestCtx)
	if err != nil {
		return nil, fmt.Errorf("初始化 Milvus 客户端失败: %w", err)
	}

	has, err := client.HasCollection(requestCtx, milvusclient.NewHasCollectionOption(cfg.Collection))
	if err != nil {
		return nil, fmt.Errorf("检查 Milvus collection 失败: %w", err)
	}
	if !has {
		return buildEmptyProductVectorBatchSearchData(req.Queries), nil
	}

	searchVectors := make([]entity.Vector, 0, len(vectors))
	textQueries := make([]entity.Vector, 0, len(req.Queries))
	for idx, vector := range vectors {
		if len(vector) == 0 {
			return nil, fmt.Errorf("第%d条产品搜索向量为空", idx+1)
		}
		query := strings.TrimSpace(req.Queries[idx])
		if query == "" {
			return nil, fmt.Errorf("第%d条产品搜索文本为空", idx+1)
		}
		searchVectors = append(searchVectors, entity.FloatVector(vector))
		textQueries = append(textQueries, entity.Text(query))
	}

	outputFields := []string{
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

	supportHybrid, err := supportsProductVectorHybridSearch(requestCtx, client, cfg.Collection)
	if err != nil {
		return nil, err
	}

	var resultSets []milvusclient.ResultSet
	if supportHybrid {
		sparseAnnParam := index.NewSparseAnnParam()
		sparseAnnParam.WithDropRatio(0.2)

		denseRequest := milvusclient.NewAnnRequest(productVectorEmbeddingField, req.Limit, searchVectors...)
		sparseRequest := milvusclient.NewAnnRequest(productVectorContextSparseField, req.Limit, textQueries...).
			WithAnnParam(sparseAnnParam)

		ranker := entity.NewFunction().
			WithName("rrf").
			WithType(entity.FunctionTypeRerank).
			WithParam("reranker", "rrf").
			WithParam("k", "100")

		resultSets, err = client.HybridSearch(requestCtx, milvusclient.NewHybridSearchOption(
			cfg.Collection,
			req.Limit,
			denseRequest,
			sparseRequest,
		).WithReranker(milvusclient.NewRRFReranker()).WithFunctionRerankers(ranker).
			WithOutputFields(outputFields...))
	} else {
		resultSets, err = client.Search(requestCtx, milvusclient.NewSearchOption(cfg.Collection, req.Limit,
			searchVectors).
			WithANNSField(productVectorEmbeddingField).
			WithOutputFields(outputFields...))
	}
	if err != nil {
		return nil, fmt.Errorf("搜索 Milvus 产品向量失败: %w", err)
	}
	if len(resultSets) == 0 {
		return buildEmptyProductVectorBatchSearchData(req.Queries), nil
	}
	if len(resultSets) != len(req.Queries) {
		return nil, fmt.Errorf("Milvus 批量搜索结果数量不匹配，expected=%d actual=%d", len(req.Queries), len(resultSets))
	}

	rows := make([]*mastermodels.ProductVectorBatchSearchItem, 0, len(resultSets))
	for idx, resultSet := range resultSets {
		data, parseErr := parseProductVectorSearchResultSet(resultSet)
		if parseErr != nil {
			return nil, fmt.Errorf("解析第%d条产品向量搜索结果失败: %w", idx+1, parseErr)
		}
		rows = append(rows, &mastermodels.ProductVectorBatchSearchItem{
			Query: req.Queries[idx],
			Rows:  data.Rows,
			Total: data.Total,
		})
	}

	return &mastermodels.ProductVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func supportsProductVectorHybridSearch(ctx context.Context, client *milvusclient.Client, collectionName string) (bool, error) {
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
		if field.Name == productVectorContextSparseField {
			return true, nil
		}
	}
	return false, nil
}

func loadProductVectorConfig() (*productVectorConfig, error) {
	cfg := &productVectorConfig{}
	_ = viper.UnmarshalKey("erp.master.product_vector", cfg)
	if cfg.Collection == "" {
		cfg.Collection = strings.TrimSpace(viper.GetString("milvus.product_collection"))
	}
	if cfg.Collection == "" {
		cfg.Collection = defaultProductVectorCollection
	}
	return cfg, nil
}

func ensureProductVectorCollection(ctx context.Context, client *milvusclient.Client, collectionName string, dim int) error {
	has, err := client.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if err != nil {
		return fmt.Errorf("检查 Milvus collection 失败: %w", err)
	}
	if !has {
		function := entity.NewFunction().
			WithName(productTextBm25EmbFunctionName).
			WithInputFields(productVectorContentField).
			WithOutputFields(productVectorContextSparseField).
			WithType(entity.FunctionTypeBM25)
		schema := entity.NewSchema().
			WithField(entity.NewField().WithName(productVectorPKField).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
			WithField(entity.NewField().WithName(productVectorNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(productVectorBarCodeField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorBarCodeMaxLength)).
			WithField(entity.NewField().WithName(productVectorCategoryIDField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(productVectorCategoryNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorCategoryMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(productVectorUnitIDField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(productVectorUnitNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorUnitMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(productVectorStandardField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorStandardMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(productVectorRemarkField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorRemarkMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(productVectorExpiryDayField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(productVectorWeightField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(productVectorPurchasePriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(productVectorSalePriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(productVectorMinPriceField).WithDataType(entity.FieldTypeDouble)).
			WithField(entity.NewField().WithName(productVectorContentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(productVectorContentMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(productVectorEmbeddingField).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(dim))).
			WithField(entity.NewField().
				WithName(productVectorContextSparseField).
				WithDataType(entity.FieldTypeSparseVector),
			).WithFunction(function)

		err = client.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(collectionName, schema).WithIndexOptions(
			milvusclient.NewCreateIndexOption(collectionName, productVectorEmbeddingField, index.NewAutoIndex(entity.COSINE)).
				WithIndexName(productIdxProductVector),
		).WithIndexOptions(milvusclient.NewCreateIndexOption(collectionName, productVectorContextSparseField,
			index.NewSparseInvertedIndex(entity.BM25, 0.2)).WithIndexName(productIdxProductTextSparseVector)))
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
		productVectorPKField:            false,
		productVectorNameField:          false,
		productVectorBarCodeField:       false,
		productVectorCategoryIDField:    false,
		productVectorCategoryNameField:  false,
		productVectorUnitIDField:        false,
		productVectorUnitNameField:      false,
		productVectorStandardField:      false,
		productVectorRemarkField:        false,
		productVectorExpiryDayField:     false,
		productVectorWeightField:        false,
		productVectorPurchasePriceField: false,
		productVectorSalePriceField:     false,
		productVectorMinPriceField:      false,
		productVectorContentField:       false,
		productVectorEmbeddingField:     false,
	}
	for _, field := range collection.Schema.Fields {
		if field == nil {
			continue
		}
		if _, ok := requiredFields[field.Name]; ok {
			requiredFields[field.Name] = true
		}
		if field.Name != productVectorEmbeddingField {
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

func buildEmptyProductVectorBatchSearchData(queries []string) *mastermodels.ProductVectorBatchSearchData {
	rows := make([]*mastermodels.ProductVectorBatchSearchItem, 0, len(queries))
	for _, query := range queries {
		rows = append(rows, &mastermodels.ProductVectorBatchSearchItem{
			Query: query,
			Rows:  make([]*mastermodels.ProductVectorSearchItem, 0),
			Total: 0,
		})
	}
	return &mastermodels.ProductVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}
}

func parseProductVectorSearchResultSet(resultSet milvusclient.ResultSet) (*mastermodels.ProductVectorSearchData, error) {
	if resultSet.Err != nil {
		return nil, resultSet.Err
	}

	rows := make([]*mastermodels.ProductVectorSearchItem, 0, resultSet.ResultCount)
	for i := 0; i < resultSet.ResultCount; i++ {
		item := &mastermodels.ProductVectorSearchItem{
			Score: resultSet.Scores[i],
		}
		var err error
		if item.ProductID, err = getProductVectorInt64ColumnValue(&resultSet, productVectorPKField, i); err != nil {
			return nil, err
		}
		if item.Name, err = getProductVectorStringColumnValue(&resultSet, productVectorNameField, i); err != nil {
			return nil, err
		}
		if item.BarCode, err = getProductVectorStringColumnValue(&resultSet, productVectorBarCodeField, i); err != nil {
			return nil, err
		}
		if item.CategoryId, err = getProductVectorInt64ColumnValue(&resultSet, productVectorCategoryIDField, i); err != nil {
			return nil, err
		}
		if item.CategoryName, err = getProductVectorStringColumnValue(&resultSet, productVectorCategoryNameField, i); err != nil {
			return nil, err
		}
		if item.UnitId, err = getProductVectorInt64ColumnValue(&resultSet, productVectorUnitIDField, i); err != nil {
			return nil, err
		}
		if item.UnitName, err = getProductVectorStringColumnValue(&resultSet, productVectorUnitNameField, i); err != nil {
			return nil, err
		}
		if item.Standard, err = getProductVectorStringColumnValue(&resultSet, productVectorStandardField, i); err != nil {
			return nil, err
		}
		if item.Remark, err = getProductVectorStringColumnValue(&resultSet, productVectorRemarkField, i); err != nil {
			return nil, err
		}
		expiryDay, err := getProductVectorInt64ColumnValue(&resultSet, productVectorExpiryDayField, i)
		if err != nil {
			return nil, err
		}
		item.ExpiryDay = int32(expiryDay)
		if item.Weight, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorWeightField, i); err != nil {
			return nil, err
		}
		if item.PurchasePrice, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorPurchasePriceField, i); err != nil {
			return nil, err
		}
		if item.SalePrice, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorSalePriceField, i); err != nil {
			return nil, err
		}
		if item.MinPrice, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorMinPriceField, i); err != nil {
			return nil, err
		}
		if item.Content, err = getProductVectorStringColumnValue(&resultSet, productVectorContentField, i); err != nil {
			return nil, err
		}
		rows = append(rows, item)
	}

	return &mastermodels.ProductVectorSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func buildProductVectorRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

func trimProductVectorRunes(value string, max int) string {
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

func getProductVectorInt64ColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (int64, error) {
	if resultSet == nil {
		return 0, fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return 0, fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsInt64(idx)
}

func getProductVectorStringColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (string, error) {
	if resultSet == nil {
		return "", fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return "", fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsString(idx)
}

func getProductVectorFloat64ColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (float64, error) {
	if resultSet == nil {
		return 0, fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return 0, fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsDouble(idx)
}
