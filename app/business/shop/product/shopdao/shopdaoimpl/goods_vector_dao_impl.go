package shopdaoimpl

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/datasource/milvus"
)

const (
	defaultGoodsVectorCollection = "shop_goods_vectors"

	goodsVectorPKField           = "goods_id"
	goodsVectorDBIDField         = "goods_db_id"
	goodsVectorNameField         = "goods_name"
	goodsVectorCodeField         = "goods_code"
	goodsVectorCategoryField     = "category_name"
	goodsVectorDescriptionField  = "description"
	goodsVectorContentField      = "content"
	goodsVectorEmbeddingField    = "vector"
	goodsVectorPKMaxLength       = 128
	goodsVectorNameMaxLength     = 512
	goodsVectorCodeMaxLength     = 128
	goodsVectorCategoryMaxLength = 256
	goodsVectorDescMaxLength     = 4096
	goodsVectorContentMaxLength  = 16384
)

type ShopGoodsVectorDaoImpl struct{}

type goodsVectorConfig struct {
	Collection string `mapstructure:"collection"`
}

func NewShopGoodsVectorDao() shopdao.IShopGoodsVectorDao {
	return &ShopGoodsVectorDaoImpl{}
}

func (d *ShopGoodsVectorDaoImpl) Upsert(c *gin.Context, goods *shopmodels.Goods, content string, vector []float32) (*shopmodels.GoodsVectorResult, error) {
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	content = trimRunes(content, goodsVectorContentMaxLength)
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("商品内容为空")
	}
	if len(vector) == 0 {
		return nil, errors.New("商品向量为空")
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

	if err = ensureGoodsVectorCollection(requestCtx, client, cfg.Collection, len(vector)); err != nil {
		return nil, err
	}

	goodsPK := strings.TrimSpace(goods.GoodsID)
	if goodsPK == "" {
		goodsPK = strconv.FormatInt(goods.ID, 10)
	}

	_, err = client.Upsert(requestCtx, milvusclient.NewColumnBasedInsertOption(cfg.Collection).
		WithVarcharColumn(goodsVectorPKField, []string{trimRunes(goodsPK, goodsVectorPKMaxLength)}).
		WithInt64Column(goodsVectorDBIDField, []int64{goods.ID}).
		WithVarcharColumn(goodsVectorNameField, []string{trimRunes(goods.GoodsName, goodsVectorNameMaxLength)}).
		WithVarcharColumn(goodsVectorCodeField, []string{trimRunes(goods.GoodsCode, goodsVectorCodeMaxLength)}).
		WithVarcharColumn(goodsVectorCategoryField, []string{trimRunes(goods.ShopCategoryName, goodsVectorCategoryMaxLength)}).
		WithVarcharColumn(goodsVectorDescriptionField, []string{trimRunes(goods.Description, goodsVectorDescMaxLength)}).
		WithVarcharColumn(goodsVectorContentField, []string{content}).
		WithFloatVectorColumn(goodsVectorEmbeddingField, len(vector), [][]float32{vector}),
	)
	if err != nil {
		return nil, fmt.Errorf("写入 Milvus 失败: %w", err)
	}

	return &shopmodels.GoodsVectorResult{
		GoodsDBID:  goods.ID,
		GoodsID:    goodsPK,
		Collection: cfg.Collection,
		Dimension:  len(vector),
	}, nil
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
		schema := entity.NewSchema().
			WithField(entity.NewField().WithName(goodsVectorPKField).WithDataType(entity.FieldTypeVarChar).WithIsPrimaryKey(true).WithMaxLength(goodsVectorPKMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorDBIDField).WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().WithName(goodsVectorNameField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorNameMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorCodeField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorCodeMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorCategoryField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorCategoryMaxLength)).
			WithField(entity.NewField().WithName(goodsVectorDescriptionField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorDescMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorContentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(goodsVectorContentMaxLength).WithEnableAnalyzer(true).WithEnableMatch(true)).
			WithField(entity.NewField().WithName(goodsVectorEmbeddingField).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(dim)))

		err = client.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(collectionName, schema).WithIndexOptions(
			milvusclient.NewCreateIndexOption(collectionName, goodsVectorEmbeddingField, index.NewAutoIndex(entity.COSINE)).WithIndexName("idx_goods_vector"),
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

	for _, field := range collection.Schema.Fields {
		if field == nil || field.Name != goodsVectorEmbeddingField {
			continue
		}
		existingDim, dimErr := field.GetDim()
		if dimErr != nil {
			return fmt.Errorf("读取 Milvus 向量维度失败: %w", dimErr)
		}
		if int(existingDim) != dim {
			return fmt.Errorf("Milvus collection %s 向量维度不匹配，现有=%d，请求=%d", collectionName, existingDim, dim)
		}
		return nil
	}

	return fmt.Errorf("Milvus collection %s 缺少 %s 字段", collectionName, goodsVectorEmbeddingField)
}

func buildRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
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
