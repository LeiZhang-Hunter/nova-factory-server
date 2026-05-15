package shopserviceimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"nova-factory-server/app/business/shop/product/shopmodels"
	embeddingutil "nova-factory-server/app/utils/llm/embedding"
)

const goodsVectorContentMaxLength = 16384
const (
	defaultGoodsVectorSearchLimit = 10
	maxGoodsVectorSearchLimit     = 50
)

type goodsVectorConfig struct {
	Embedding embeddingutil.ProviderConfig `mapstructure:"embedding"`
}

type goodsEmbeddingPayload struct {
	sku     *shopmodels.GoodsSku
	content string
}

func (s *ShopGoodsServiceImpl) GenerateVector(c *gin.Context, req *shopmodels.GenVectorReq) (*shopmodels.GoodsVectorResult, error) {
	goods, err := s.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	if goods.IsOnSale <= 0 {
		return nil, errors.New("商品未上架，不能生成向量")
	}

	if err = s.attachSkus(c, []*shopmodels.Goods{goods}); err != nil {
		return nil, err
	}

	if err = s.attachCategoryNames(c, []*shopmodels.Goods{goods}); err != nil {
		return nil, err
	}

	cfg, err := s.loadGoodsVectorConfig(c, req)
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, &cfg.Embedding)
	if err != nil {
		return nil, fmt.Errorf("初始化向量模型失败: %w", err)
	}

	payloads := buildGoodsEmbeddingPayloads(goods)
	if len(payloads) == 0 {
		return nil, errors.New("商品内容为空，无法生成向量")
	}

	texts := make([]string, 0, len(payloads))
	for _, payload := range payloads {
		texts = append(texts, payload.content)
	}

	vectors, err := embedder.EmbedStrings(requestCtx, texts)
	if err != nil {
		return nil, fmt.Errorf("生成商品向量失败: %w", err)
	}
	if len(vectors) != len(payloads) {
		return nil, fmt.Errorf("向量模型返回结果数量不匹配，expected=%d actual=%d", len(payloads), len(vectors))
	}
	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, errors.New("向量模型未返回有效结果")
	}

	items := make([]*shopmodels.GoodsVectorUpsertItem, 0, len(payloads))
	for idx, payload := range payloads {
		if len(vectors[idx]) == 0 {
			return nil, fmt.Errorf("第%d条SKU向量为空", idx+1)
		}
		item := &shopmodels.GoodsVectorUpsertItem{
			Content: payload.content,
			Vector:  float64ToFloat32(vectors[idx]),
		}
		if payload.sku != nil {
			item.SkuID = int64(payload.sku.ID)
			item.SkuName = payload.sku.SkuName
			item.SkuDescription = payload.sku.Description
			item.RetailPrice = payload.sku.RetailPrice
			item.Weight = payload.sku.Weight
			item.Quantity = payload.sku.Quantity
		}
		items = append(items, item)
	}
	result, err := s.vectorDao.Upsert(c, goods, items)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("未生成有效的商品向量数据")
	}
	return result, nil
}

func (s *ShopGoodsServiceImpl) SearchVector(c *gin.Context, req *shopmodels.GoodsVectorSearchReq) (*shopmodels.GoodsVectorSearchData, error) {
	if req == nil {
		return nil, errors.New("搜索参数不能为空")
	}
	queries, err := normalizeGoodsVectorSearchQueries([]string{req.Query})
	if err != nil {
		return nil, err
	}
	items, err := s.batchSearchVector(c, queries, req.Limit, req.Embedding)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return &shopmodels.GoodsVectorSearchData{
			Rows:  make([]*shopmodels.GoodsVectorSearchItem, 0),
			Total: 0,
		}, nil
	}
	return &shopmodels.GoodsVectorSearchData{
		Rows:  items[0].Rows,
		Total: items[0].Total,
	}, nil
}

func (s *ShopGoodsServiceImpl) BatchSearchVector(c *gin.Context,
	req *shopmodels.GoodsVectorBatchSearchReq) (*shopmodels.GoodsVectorBatchSearchData, error) {
	if req == nil {
		return nil, errors.New("批量搜索参数不能为空")
	}
	queries, err := normalizeGoodsVectorSearchQueries(req.Queries)
	if err != nil {
		return nil, err
	}
	rows, err := s.batchSearchVector(c, queries, req.Limit, req.Embedding)
	if err != nil {
		return nil, err
	}
	return &shopmodels.GoodsVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func (s *ShopGoodsServiceImpl) batchSearchVector(c *gin.Context, queries []string, limit int,
	embedding *shopmodels.EmbeddingConfig) ([]*shopmodels.GoodsVectorBatchSearchItem, error) {
	cfg, err := loadEmbeddingProviderConfig(embedding)
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化向量模型失败: %w", err)
	}

	vectors, err := embedder.EmbedStrings(requestCtx, queries)
	if err != nil {
		return nil, fmt.Errorf("生成搜索向量失败: %w", err)
	}
	if len(vectors) != len(queries) {
		return nil, fmt.Errorf("向量模型返回结果数量不匹配，expected=%d actual=%d", len(queries), len(vectors))
	}

	normalizedVectors := make([][]float32, 0, len(vectors))
	for idx := range queries {
		if len(vectors[idx]) == 0 {
			return nil, fmt.Errorf("第%d条搜索向量为空", idx+1)
		}
		normalizedVectors = append(normalizedVectors, float64ToFloat32(vectors[idx]))
	}
	data, err := s.vectorDao.BatchSearch(c, &shopmodels.GoodsVectorBatchSearchReq{
		Queries: queries,
		Limit:   normalizeGoodsVectorSearchLimit(limit),
	}, normalizedVectors)
	if err != nil {
		return nil, err
	}
	if data == nil || data.Rows == nil {
		return make([]*shopmodels.GoodsVectorBatchSearchItem, 0), nil
	}
	return data.Rows, nil
}

func (s *ShopGoodsServiceImpl) loadGoodsVectorConfig(c *gin.Context,
	req *shopmodels.GenVectorReq) (*goodsVectorConfig, error) {
	cfg := &goodsVectorConfig{}
	embeddingCfg, err := loadEmbeddingProviderConfig(req.Embedding)
	if err != nil {
		return nil, err
	}
	cfg.Embedding = *embeddingCfg
	return cfg, nil
}

func loadEmbeddingProviderConfig(req *shopmodels.EmbeddingConfig) (*embeddingutil.ProviderConfig, error) {
	if req == nil {
		return nil, errors.New("未配置 embedding.model_id 或 shop.goods_vector.embedding.model_id")
	}
	cfg := &embeddingutil.ProviderConfig{
		ProviderType: req.ProviderType,
		ProviderID:   req.ProviderID,
		APIEndpoint:  req.APIEndpoint,
		ModelID:      req.ModelID,
		APIKey:       req.ApiKey,
	}
	if cfg.ProviderType == "" {
		cfg.ProviderType = "openai"
	}
	if cfg.ProviderID == "" {
		cfg.ProviderID = cfg.ProviderType
	}
	if cfg.ModelID == "" {
		return nil, errors.New("未配置 embedding.model_id 或 shop.goods_vector.embedding.model_id")
	}
	return cfg, nil
}

func buildRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

func float64ToFloat32(vector []float64) []float32 {
	result := make([]float32, 0, len(vector))
	for _, value := range vector {
		result = append(result, float32(value))
	}
	return result
}

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

func normalizeGoodsVectorSearchQueries(queries []string) ([]string, error) {
	result := make([]string, 0, len(queries))
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		result = append(result, query)
	}
	if len(result) == 0 {
		return nil, errors.New("搜索内容不能为空")
	}
	return result, nil
}

func buildGoodsEmbeddingPayloads(goods *shopmodels.Goods) []goodsEmbeddingPayload {
	if goods == nil {
		return nil
	}

	baseLines := make([]string, 0, 9)
	appendLine := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		baseLines = append(baseLines, label+": "+value)
	}

	appendLine("商品名称", goods.GoodsName)
	appendLine("商品分类", goods.ShopCategoryName)
	appendLine("商品描述", goods.Description)
	//if goods.RetailPrice > 0 {
	//	baseLines = append(baseLines, fmt.Sprintf("零售价: %.2f", goods.RetailPrice))
	//}
	//if goods.Weight > 0 {
	//	baseLines = append(baseLines, fmt.Sprintf("重量: %.2f%s", goods.Weight, strings.TrimSpace(goods.WeightUnit)))
	//}
	if goods.Unit != "" {
		baseLines = append(baseLines, "销售单位: "+strings.TrimSpace(goods.Unit))
	}
	//baseLines = append(baseLines, fmt.Sprintf("库存: %d", goods.Quantity))
	//if goods.IsOnSale > 0 {
	//	baseLines = append(baseLines, "上架状态: 上架")
	//} else {
	//	baseLines = append(baseLines, "上架状态: 下架")
	//}

	payloads := make([]goodsEmbeddingPayload, 0, maxInt(len(goods.Skus), 1))
	for _, sku := range goods.Skus {
		if sku == nil {
			continue
		}
		lines := append([]string{}, baseLines...)
		appendSkuLines(&lines, sku)
		content := trimRunes(strings.Join(lines, "\n"), goodsVectorContentMaxLength)
		if strings.TrimSpace(content) == "" {
			continue
		}
		payloads = append(payloads, goodsEmbeddingPayload{
			sku:     sku,
			content: content,
		})
	}
	if len(payloads) > 0 {
		return payloads
	}
	content := trimRunes(strings.Join(baseLines, "\n"), goodsVectorContentMaxLength)
	if strings.TrimSpace(content) == "" {
		return nil
	}
	return []goodsEmbeddingPayload{{content: content}}
}

func appendSkuLines(lines *[]string, sku *shopmodels.GoodsSku) {
	if sku == nil {
		return
	}
	appendLine := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		*lines = append(*lines, label+": "+value)
	}
	appendLine("SKU名称", sku.SkuName)
	appendLine("SKU描述", sku.Description)
	if sku.Unit != "" {
		*lines = append(*lines, "SKU单位: "+strings.TrimSpace(sku.Unit))
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalizeGoodsVectorSearchLimit(limit int) int {
	if limit <= 0 {
		return defaultGoodsVectorSearchLimit
	}
	if limit > maxGoodsVectorSearchLimit {
		return maxGoodsVectorSearchLimit
	}
	return limit
}
