package masterdaoimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/spf13/viper"
)

// supportsProductVectorHybridSearch 检查当前 collection 是否具备稀疏向量字段。
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

// loadProductVectorConfig 读取产品向量相关配置。
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

// ensureProductVectorCollection 确保 Milvus collection 存在且字段结构与向量维度符合预期。
func ensureProductVectorCollection(ctx context.Context, client *milvusclient.Client, collectionName string, dim int) error {
	has, err := client.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if err != nil {
		return fmt.Errorf("检查 Milvus collection 失败: %w", err)
	}
	if !has {
		// 新建 collection 时同时挂上 BM25 function，让 content 字段自动生成稀疏向量。
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

	// 已存在 collection 时只做校验，不在运行期强行修改结构，避免影响线上数据。
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
