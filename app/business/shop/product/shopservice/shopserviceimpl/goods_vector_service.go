package shopserviceimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"nova-factory-server/app/business/shop/product/shopmodels"
	embeddingutil "nova-factory-server/app/utils/llm/embedding"
)

const goodsVectorContentMaxLength = 16384

type goodsVectorConfig struct {
	Embedding embeddingutil.ProviderConfig `mapstructure:"embedding"`
}

func (s *ShopGoodsServiceImpl) GenerateVector(c *gin.Context, id int64) (*shopmodels.GoodsVectorResult, error) {
	goods, err := s.GetByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}
	if goods == nil {
		return nil, errors.New("商品不存在")
	}

	cfg, err := loadGoodsVectorConfig()
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, &cfg.Embedding)
	if err != nil {
		return nil, fmt.Errorf("初始化向量模型失败: %w", err)
	}

	content := trimRunes(buildGoodsEmbeddingText(goods), goodsVectorContentMaxLength)
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("商品内容为空，无法生成向量")
	}

	vectors, err := embedder.EmbedStrings(requestCtx, []string{content})
	if err != nil {
		return nil, fmt.Errorf("生成商品向量失败: %w", err)
	}
	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, errors.New("向量模型未返回有效结果")
	}

	return s.vectorDao.Upsert(c, goods, content, float64ToFloat32(vectors[0]))
}

func loadGoodsVectorConfig() (*goodsVectorConfig, error) {
	cfg := &goodsVectorConfig{}
	_ = viper.UnmarshalKey("shop.goods_vector", cfg)

	if cfg.Embedding.ModelID == "" && viper.IsSet("embedding") {
		if err := viper.UnmarshalKey("embedding", &cfg.Embedding); err != nil {
			return nil, fmt.Errorf("读取 embedding 配置失败: %w", err)
		}
	}

	cfg.Embedding.ProviderType = strings.TrimSpace(cfg.Embedding.ProviderType)
	cfg.Embedding.ProviderID = strings.TrimSpace(cfg.Embedding.ProviderID)
	cfg.Embedding.APIEndpoint = strings.TrimSpace(cfg.Embedding.APIEndpoint)
	cfg.Embedding.ModelID = strings.TrimSpace(cfg.Embedding.ModelID)

	if cfg.Embedding.ProviderType == "" {
		cfg.Embedding.ProviderType = "openai"
	}
	if cfg.Embedding.ProviderID == "" {
		cfg.Embedding.ProviderID = cfg.Embedding.ProviderType
	}
	if cfg.Embedding.ModelID == "" {
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

func buildGoodsEmbeddingText(goods *shopmodels.Goods) string {
	if goods == nil {
		return ""
	}

	lines := make([]string, 0, 8+len(goods.Skus))
	appendLine := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		lines = append(lines, label+": "+value)
	}

	appendLine("商品名称", goods.GoodsName)
	appendLine("商品编码", goods.GoodsCode)
	appendLine("商品分类", goods.ShopCategoryName)
	appendLine("商品描述", goods.Description)
	if goods.RetailPrice > 0 {
		lines = append(lines, fmt.Sprintf("零售价: %.2f", goods.RetailPrice))
	}
	if goods.Weight > 0 {
		lines = append(lines, fmt.Sprintf("重量: %.2f%s", goods.Weight, strings.TrimSpace(goods.WeightUnit)))
	}
	if goods.Unit != "" {
		lines = append(lines, "销售单位: "+strings.TrimSpace(goods.Unit))
	}
	lines = append(lines, fmt.Sprintf("库存: %d", goods.Quantity))
	if goods.IsOnSale > 0 {
		lines = append(lines, "上架状态: 上架")
	} else {
		lines = append(lines, "上架状态: 下架")
	}

	for idx, sku := range goods.Skus {
		if sku == nil {
			continue
		}
		parts := make([]string, 0, 8)
		if sku.SkuName != "" {
			parts = append(parts, "名称="+strings.TrimSpace(sku.SkuName))
		}
		if sku.SkuCode != "" {
			parts = append(parts, "编码="+strings.TrimSpace(sku.SkuCode))
		}
		if sku.Barcode != "" {
			parts = append(parts, "条码="+strings.TrimSpace(sku.Barcode))
		}
		if sku.Description != "" {
			parts = append(parts, "描述="+strings.TrimSpace(sku.Description))
		}
		if sku.RetailPrice > 0 {
			parts = append(parts, fmt.Sprintf("价格=%.2f", sku.RetailPrice))
		}
		if sku.Weight > 0 {
			parts = append(parts, fmt.Sprintf("重量=%.2f%s", sku.Weight, strings.TrimSpace(sku.WeightUnit)))
		}
		if sku.Unit != "" {
			parts = append(parts, "单位="+strings.TrimSpace(sku.Unit))
		}
		parts = append(parts, fmt.Sprintf("库存=%d", sku.Quantity))
		if len(parts) == 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("规格%d: %s", idx+1, strings.Join(parts, "，")))
	}

	return strings.Join(lines, "\n")
}
