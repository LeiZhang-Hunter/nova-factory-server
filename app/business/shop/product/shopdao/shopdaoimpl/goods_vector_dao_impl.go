package shopdaoimpl

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"nova-factory-server/app/utils/vectorsearch/goods"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/datasource/milvus"
)

// goodsSearchRuntimeGroup 用于把 filterExpr 相同的 query 合并成一个 Milvus 请求。
// Milvus 一次 Search/HybridSearch 只能挂一份 filter，所以不同过滤条件不能混在同一个请求里。
type goodsSearchRuntimeGroup struct {
	filterExpr string
	queries    []goodsSearchRuntimeQuery
}

// NewShopGoodsVectorDao 创建商品向量 DAO 实现。
func NewShopGoodsVectorDao() shopdao.IShopGoodsVectorDao {
	metadataExtractor := goods.NewMetadataExtractor()
	err := metadataExtractor.Init()
	if err != nil {
		zap.L().Error("init metadata extractor failed", zap.Error(err))
	}
	return &ShopGoodsVectorDaoImpl{
		metadataExtractor: metadataExtractor,
	}
}

// Upsert 将商品及其 SKU 向量写入 Milvus。
// 每条 SKU 会落为一行，便于检索时保留更细的规格粒度。
func (d *ShopGoodsVectorDaoImpl) Upsert(c *gin.Context, goods *shopmodels.Goods, items []*shopmodels.GoodsVectorUpsertItem) (*shopmodels.GoodsVectorResult, error) {
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	if len(items) == 0 {
		return nil, errors.New("商品向量写入项为空")
	}

	dim := 0
	for idx, item := range items {
		// 写入前先做最基本的内容与向量校验，避免半路触发 Milvus 参数错误。
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

	// 先将业务对象展开成列式写入所需的数据结构，再统一 upsert。
	rows, err := buildGoodsVectorRows(goods, items, d.metadataExtractor)
	if err != nil {
		return nil, err
	}
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
			column.NewColumnJSONBytes(goodsVectorMetadataField, rows.metadatas),
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

// Search 复用批量检索入口处理单条查询。
func (d *ShopGoodsVectorDaoImpl) Search(c *gin.Context, req *shopmodels.GoodsVectorSearchReq, vector []float32, FallbackWithoutMetadata bool) (*shopmodels.GoodsVectorSearchData, error) {
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
	}, [][]float32{vector}, FallbackWithoutMetadata)
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

// BatchSearch 执行商品批量向量检索。
// 支持 dense-only 检索，以及 dense + BM25 稀疏向量的 hybrid 检索。
func (d *ShopGoodsVectorDaoImpl) BatchSearch(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq,
	vectors [][]float32, FallbackWithoutMetadata bool) (*shopmodels.GoodsVectorBatchSearchData, error) {
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

	// 先扩大召回候选集，再在应用层做一次业务重排。
	searchLimit := resolveGoodsVectorSearchCandidateLimit(req.Limit)
	runtimeQueries := make([]goodsSearchRuntimeQuery, 0, len(vectors))
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
			// SearchTexts 允许调用方覆盖 query，用于传入更适合召回的文本。
			if candidate := strings.TrimSpace(req.SearchTexts[idx]); candidate != "" {
				searchText = candidate
			}
		}

		filterExpr := ""
		if d.metadataExtractor != nil {
			// 对搜索文本做规格/分类提取：
			// 1. 提取后的 Value 可替换为更干净的文本参与召回；
			// 2. 提取后的 Metadata 会被拼成 Milvus filter，用于精确过滤。
			extract, extractErr := d.metadataExtractor.Extract(searchText)
			if extractErr != nil {
				zap.L().Error("extract goods search metadata fail",
					zap.String("query", query),
					zap.String("searchText", searchText),
					zap.Error(extractErr))
			} else {
				if strings.TrimSpace(extract.Value) != "" {
					searchText = extract.Value
				}
				filterExpr = buildGoodsMetadataFilterExpr(extract)
			}
		}
		runtimeQueries = append(runtimeQueries, goodsSearchRuntimeQuery{
			index:      len(runtimeQueries),
			query:      query,
			vector:     entity.FloatVector(vector),
			text:       entity.Text(searchText),
			filterExpr: filterExpr,
		})
	}
	if len(runtimeQueries) == 0 {
		return buildEmptyGoodsVectorBatchSearchData(req.Queries), nil
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

	// 由于 filter 是请求级配置，不同 filterExpr 的 query 不能合并到同一次请求。
	// 这里按 filterExpr 分组后批量查询，兼顾正确性和请求数控制。
	groups := groupGoodsSearchRuntimeQueries(runtimeQueries)

	resultSets := make([]milvusclient.ResultSet, len(runtimeQueries))
	if supportHybrid {
		if err = executeGoodsHybridSearchGroups(requestCtx, client, cfg.Collection, searchLimit, outputFields, groups, resultSets); err != nil {
			return nil, err
		}
	} else {
		if err = executeGoodsDenseSearchGroups(requestCtx, client, cfg.Collection, searchLimit, outputFields, groups, resultSets); err != nil {
			return nil, err
		}
	}
	if FallbackWithoutMetadata {
		fallbackQueries := collectGoodsMetadataFallbackQueries(runtimeQueries, resultSets)
		if len(fallbackQueries) > 0 {
			fallbackGroups := groupGoodsSearchRuntimeQueries(fallbackQueries)
			if supportHybrid {
				if err = executeGoodsHybridSearchGroups(requestCtx, client, cfg.Collection, searchLimit, outputFields, fallbackGroups, resultSets); err != nil {
					return nil, err
				}
			} else {
				if err = executeGoodsDenseSearchGroups(requestCtx, client, cfg.Collection, searchLimit, outputFields, fallbackGroups, resultSets); err != nil {
					return nil, err
				}
			}
		}
	}
	if len(resultSets) == 0 {
		queries := make([]string, 0, len(runtimeQueries))
		for _, runtimeQuery := range runtimeQueries {
			queries = append(queries, runtimeQuery.query)
		}
		return buildEmptyGoodsVectorBatchSearchData(queries), nil
	}
	if len(resultSets) != len(runtimeQueries) {
		return nil, fmt.Errorf("Milvus 批量搜索结果数量不匹配，expected=%d actual=%d", len(runtimeQueries), len(resultSets))
	}

	rows := make([]*shopmodels.GoodsVectorBatchSearchItem, 0, len(resultSets))
	for idx, resultSet := range resultSets {
		data, parseErr := parseGoodsVectorSearchResultSet(resultSet)
		if parseErr != nil {
			return nil, fmt.Errorf("解析第%d条商品向量搜索结果失败: %w", idx+1, parseErr)
		}
		data.Rows = rerankGoodsVectorSearchRows(runtimeQueries[idx].query, data.Rows, req.Limit)
		rows = append(rows, &shopmodels.GoodsVectorBatchSearchItem{
			Query: runtimeQueries[idx].query,
			Rows:  data.Rows,
			Total: int64(len(data.Rows)),
		})
	}

	return &shopmodels.GoodsVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}
