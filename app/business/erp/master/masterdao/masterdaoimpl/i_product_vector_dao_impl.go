package masterdaoimpl

import (
	"errors"
	"fmt"
	"strings"

	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/datasource/milvus"
	searchutil "nova-factory-server/app/utils/vectorsearch"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// Upsert 写入或更新产品向量及其检索字段。
func (d *ProductVectorDaoImpl) Upsert(c *gin.Context, product *mastermodels.Product, item *mastermodels.ProductVectorUpsertItem) (*mastermodels.ProductVectorResult, error) {
	if product == nil || product.ID <= 0 {
		return nil, errors.New("产品不存在")
	}
	if item == nil {
		return nil, errors.New("产品向量写入项为空")
	}
	item.Content = searchutil.TrimRunes(item.Content, productVectorContentMaxLength)
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

	// 同步写入结构化字段、全文检索字段和稠密向量，便于后续做混合检索与二次重排。
	_, err = client.Upsert(requestCtx, milvusclient.NewColumnBasedInsertOption(cfg.Collection).
		WithInt64Column(productVectorPKField, []int64{product.ID}).
		WithVarcharColumn(productVectorNameField, []string{searchutil.TrimRunes(product.Name, productVectorNameMaxLength)}).
		WithVarcharColumn(productVectorBarCodeField, []string{searchutil.TrimRunes(product.BarCode, productVectorBarCodeMaxLength)}).
		WithInt64Column(productVectorCategoryIDField, []int64{product.CategoryId}).
		WithVarcharColumn(productVectorCategoryNameField, []string{searchutil.TrimRunes(product.CategoryName, productVectorCategoryMaxLength)}).
		WithInt64Column(productVectorUnitIDField, []int64{product.UnitId}).
		WithVarcharColumn(productVectorUnitNameField, []string{searchutil.TrimRunes(product.UnitName, productVectorUnitMaxLength)}).
		WithVarcharColumn(productVectorStandardField, []string{searchutil.TrimRunes(product.Standard, productVectorStandardMaxLength)}).
		WithVarcharColumn(productVectorRemarkField, []string{searchutil.TrimRunes(product.Remark, productVectorRemarkMaxLength)}).
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

// Search 单条搜索复用批量搜索逻辑，统一召回与重排行为。
func (d *ProductVectorDaoImpl) Search(c *gin.Context, req *mastermodels.ProductVectorSearchReq, vector []float32) (*mastermodels.ProductVectorSearchData, error) {
	if req == nil {
		return nil, errors.New("产品向量搜索参数为空")
	}
	if len(vector) == 0 {
		return nil, errors.New("产品搜索向量为空")
	}
	data, err := d.BatchSearch(c, &mastermodels.ProductVectorBatchSearchReq{
		Queries:     []string{req.Query},
		SearchTexts: []string{req.SearchText},
		Limit:       req.Limit,
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

// BatchSearch 先从 Milvus 扩大候选召回，再结合字段特征做应用层重排。
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
	if err = milvus.EnsureCollectionLoaded(requestCtx, client, cfg.Collection); err != nil {
		return nil, err
	}

	searchLimit := resolveProductVectorSearchCandidateLimit(req.Limit)
	plans, err := buildProductVectorSearchPlans(req, vectors)
	if err != nil {
		return nil, err
	}
	resultSets, err := executeProductVectorSearch(requestCtx, client, cfg.Collection, searchLimit, buildProductVectorOutputFields(), plans)
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
		// 应用层重排会综合标题、条码、分类、规格、正文等字段信号重新排序，并动态裁剪弱相关结果。
		data.Rows = rerankProductVectorSearchRows(req.Queries[idx], data.Rows, req.Limit)
		rows = append(rows, &mastermodels.ProductVectorBatchSearchItem{
			Query: req.Queries[idx],
			Rows:  data.Rows,
			Total: int64(len(data.Rows)),
		})
	}

	return &mastermodels.ProductVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}
